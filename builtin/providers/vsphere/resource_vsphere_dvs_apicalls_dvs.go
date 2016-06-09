package vsphere

import (
	"fmt"
	"strings"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
	"golang.org/x/net/context"
)

// methods for dvs objects

func (d *dvs) makeDVSCreateSpec() types.DVSCreateSpec {
	return types.DVSCreateSpec{
		ConfigSpec: &types.DVSConfigSpec{
			Contact: &types.DVSContactInfo{
				Contact: d.contact.infos,
				Name:    d.contact.name,
			},
			ExtensionKey:       d.extensionKey,
			Description:        d.description,
			Name:               d.name,
			NumStandalonePorts: int32(d.numStandalonePorts),
			Policy: &types.DVSPolicy{
				AutoPreInstallAllowed: &d.switchUsagePolicy.autoPreinstallAllowed,
				AutoUpgradeAllowed:    &d.switchUsagePolicy.autoUpgradeAllowed,
				PartialUpgradeAllowed: &d.switchUsagePolicy.partialUpgradeAllowed,
			},
			SwitchIpAddress: d.switchIPAddress,
		},
	}

}

func (d *dvs) getDCAndFolders(c *govmomi.Client) (*object.Datacenter, *object.DatacenterFolders, error) {
	datacenter, err := getDatacenter(c, d.datacenter)
	if err != nil {
		return nil, nil, fmt.Errorf("Cannot get datacenter from %+v [%+v]", d, err)
	}

	// get Network Folder from datacenter
	dcFolders, err := datacenter.Folders(context.TODO())
	if err != nil {
		return nil, nil, fmt.Errorf("Cannot get folders for datacenter %+v [%+v]", datacenter, err)
	}
	return datacenter, dcFolders, nil
}

func (d *dvs) addHost(c *govmomi.Client, host string, nicNames []string) error {
	dvsPath := fmt.Sprintf("%s/%s", d.folder, d.name)
	dvsItem, err := d.getDVS(c, dvsPath)
	dvsStruct := mo.DistributedVirtualSwitch{}
	if err = dvsItem.Properties(
		context.TODO(),
		dvsItem.Reference(),
		[]string{"capability", "config", "networkResourcePool", "portgroup", "summary", "uuid"},
		&dvsStruct); err != nil {
		return fmt.Errorf("Could not get properties for %s", dvsItem)
	}
	config := dvsStruct.Config.GetDVSConfigInfo()
	hostObj, err := getHost(c, d.datacenter, host)
	if err != nil {
		return fmt.Errorf("Could not get host %s: %+v", host, err)
	}
	hostref := hostObj.Reference()

	var pnicSpecs []types.DistributedVirtualSwitchHostMemberPnicSpec
	for _, nic := range nicNames {
		pnicSpecs = append(
			pnicSpecs,
			types.DistributedVirtualSwitchHostMemberPnicSpec{
				PnicDevice: nic,
			})
	}
	newHost := types.DistributedVirtualSwitchHostMemberConfigSpec{
		Host:      hostref,
		Operation: "add",
		Backing: &types.DistributedVirtualSwitchHostMemberPnicBacking{
			PnicSpec: pnicSpecs,
		},
	}
	configSpec := types.DVSConfigSpec{
		ConfigVersion: config.ConfigVersion,
		Host:          []types.DistributedVirtualSwitchHostMemberConfigSpec{newHost},
	}
	task, err := dvsItem.Reconfigure(context.TODO(), &configSpec)
	if err != nil {
		return fmt.Errorf("Could not reconfigure the DVS: %+v", err)
	}
	return waitForTaskEnd(task, "Could not reconfigure the DVS: %+v")
}

func (d *dvs) createSwitch(c *govmomi.Client) error {
	_, folders, err := d.getDCAndFolders(c)
	if err != nil {
		return fmt.Errorf("Could not get datacenter and  folders: %+v", err)
	}
	folder := folders.NetworkFolder

	// using Network Folder, create the DVSCreateSpec (pretty much a mapping of the config)
	spec := d.makeDVSCreateSpec()
	task, err := folder.CreateDVS(context.TODO(), spec)
	if err != nil {
		return fmt.Errorf("Could not create the DVS: %+v", err)
	}
	_, err = task.WaitForResult(context.TODO(), nil)
	if err != nil {
		return fmt.Errorf("Could not create the DVS: %+v", err)
	}
	return nil
}

// get a DVS from its name and populate the DVS with its infos
func (d *dvs) getDVS(c *govmomi.Client, dvsPath string) (*object.DistributedVirtualSwitch, error) {
	datacenter, _, err := d.getDCAndFolders(c)
	if err != nil {
		return nil, fmt.Errorf("Could not get datacenter and  folders: %+v", err)
	}
	finder := find.NewFinder(c.Client, true)
	finder.SetDatacenter(datacenter)
	res, err := finder.Network(context.TODO(), dvsPath)
	if err != nil {
		return nil, fmt.Errorf("Cannot get DVS %s:( %+v", dvsPath, err)
	}
	castedobj, casted := res.(*object.DistributedVirtualSwitch)
	if !casted {
		return nil, fmt.Errorf("Oops! Object %s is not a DVS but a %T", res, res)
	}
	return castedobj, nil
}

// load a DVS and populate the struct with it
func (d *dvs) loadDVS(c *govmomi.Client, datacenter, dvsName string) error {
	var dvsMo mo.DistributedVirtualSwitch
	dvsobj, err := d.getDVS(c, dvsName)
	if err != nil {
		return err
	}
	tokens := strings.Split(dvsName, "/")
	folder := strings.Join(tokens[:len(tokens)-1], "/")
	// retrieve the DVS properties

	err = dvsobj.Properties(
		context.TODO(),
		dvsobj.Reference(),
		[]string{"capability", "config", "networkResourcePool", "portgroup", "summary", "uuid"},
		&dvsMo)
	if err != nil {
		return fmt.Errorf("Could not retrieve properties: %+v", err)
	}
	// populate the struct from the data
	dvsci := dvsMo.Config.GetDVSConfigInfo()
	d.folder = folder
	d.contact.infos = dvsci.Contact.Contact
	d.contact.name = dvsci.Contact.Name
	d.description = dvsci.Description
	d.extensionKey = dvsci.ExtensionKey
	d.name = dvsci.Name
	d.numStandalonePorts = int(dvsci.NumStandalonePorts)
	d.switchIPAddress = dvsci.SwitchIpAddress
	d.switchUsagePolicy.autoPreinstallAllowed = *dvsci.Policy.AutoPreInstallAllowed
	d.switchUsagePolicy.autoUpgradeAllowed = *dvsci.Policy.AutoUpgradeAllowed
	d.switchUsagePolicy.partialUpgradeAllowed = *dvsci.Policy.PartialUpgradeAllowed
	// return nil: no error
	return nil
}

func (d *dvs) getDVSHostMembers(c *govmomi.Client) (out map[string]*dvs_map_host_dvs, err error) {
	properties, err := d.getProperties(c)
	if err != nil {
		return
	}
	hostInfos := properties.Config.GetDVSConfigInfo().Host
	// now we need to populate out
	for _, hostmember := range hostInfos {
		mapObj, err := loadMapHostDVS(d.getFullName(), hostmember)
		if err != nil {
			return nil, err
		}
		out[mapObj.hostName] = mapObj
	}
	return
}

func (d *dvs) getProperties(c *govmomi.Client) (out *mo.DistributedVirtualSwitch, err error) {

	dvsMo := mo.DistributedVirtualSwitch{}
	dvsobj, err := d.getDVS(c, d.name)
	if err != nil {
		return nil, err
	}
	return &dvsMo, dvsobj.Properties(
		context.TODO(),
		dvsobj.Reference(),
		[]string{"capability", "config", "host", "networkResourcePool", "portgroup", "summary", "uuid"},
		&dvsMo)
}
