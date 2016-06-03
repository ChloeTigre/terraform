package vsphere

import (
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/govmomi/vim25/mo"
	"golang.org/x/net/context"
	"fmt"
	"log"
	"strings"
)

func (d *dvs) makeDVSConfigSpec() (types.DVSCreateSpec, error) {
	return types.DVSCreateSpec{
		ConfigSpec: &types.DVSConfigSpec{
			Contact: &types.DVSContactInfo{
				Contact: d.contact.infos,
				Name: d.contact.name,
			},
			ExtensionKey: d.extensionKey,
			Description: d.description,
			Name: d.name,
			NumStandalonePorts: int32(d.numStandalonePorts),
			Policy: &types.DVSPolicy{
				AutoPreInstallAllowed: &d.switchUsagePolicy.autoPreinstallAllowed,
				AutoUpgradeAllowed: &d.switchUsagePolicy.autoUpgradeAllowed,
				PartialUpgradeAllowed: &d.switchUsagePolicy.partialUpgradeAllowed,
			},
			SwitchIpAddress: d.switchIPAddress,

		},
	}, nil

}

func (d *dvs) switchBoilerplate(c *govmomi.Client) (*object.Datacenter, *object.DatacenterFolders, error) {
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


func (d *dvs) createSwitch(c *govmomi.Client) error {
	_, folders, err := d.switchBoilerplate(c)
	if err != nil {
		return fmt.Errorf("Could not get datacenter and  folders: %+v", err)
	}
	folder := folders.NetworkFolder

	// using Network Folder, create the DVSCreateSpec (pretty much a mapping of the config)
	spec, err := d.makeDVSConfigSpec()
	if err != nil {
		return fmt.Errorf("Could not make DVSConfig properly: %+v", err)
	}
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
func (d *dvs) getDVS(c *govmomi.Client, dvsName string) (object.NetworkReference, error) {
	datacenter, _, err := d.switchBoilerplate(c)
	if err != nil {
		return nil, fmt.Errorf("Could not get datacenter and  folders: %+v", err)
	}
	finder := find.NewFinder(c.Client, true)
	finder.SetDatacenter(datacenter)
	res, err := finder.Network(context.TODO(), dvsName)
	if err != nil {
		return nil, fmt.Errorf("Cannot get DVS %s:( %+v", dvsName, err)
	}
	log.Printf("Result of getDVS: %+v", res)
	return res, nil
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

	castedobj, casted := dvsobj.(*object.DistributedVirtualSwitch)
	if !casted {
		return fmt.Errorf("Oops! Object %s is not a DVS", dvsName)
	}
	err = castedobj.Properties(
		context.TODO(),
		castedobj.Reference(),
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

func loadDVS(c *govmomi.Client, datacenter, dvsName string, output dvs) error {
	output.datacenter = datacenter
	err := output.loadDVS(c, datacenter, dvsName)
	return err
}
