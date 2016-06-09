package vsphere

import (
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/methods"
	"golang.org/x/net/context"
	"fmt"
	"log"
	"strings"
)

/*
this part should be merged in govmomi (module object) when we get a chance to do so
*/

func createDVPortgroup(c *govmomi.Client, dvsRef object.NetworkReference, spec types.DVPortgroupConfigSpec) (*object.Task, error) {
	req := types.CreateDVPortgroup_Task{
		Spec: spec,
		This: dvsRef.Reference(),
	}
	res, err := methods.CreateDVPortgroup_Task(context.TODO(), c.Client, &req)
	if err != nil {
		return nil, err
	}
	return object.NewTask(c.Client, res.Returnval), nil
}

/* end of what should be contributed to govmomi */

// methods for dvs objects

func (d *dvs) makeDVSCreateSpec() types.DVSCreateSpec {
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
		Host: hostref,
		Operation: "add",
		Backing: &types.DistributedVirtualSwitchHostMemberPnicBacking{
			PnicSpec: pnicSpecs,
		},
	}
	configSpec := types.DVSConfigSpec{
		ConfigVersion: config.ConfigVersion,
		Host: []types.DistributedVirtualSwitchHostMemberConfigSpec{newHost},

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

// portgroup methods

func (p *dvs_port_group) getVmomiDVPG(c *govmomi.Client, datacenter, switchName, name string) (*object.DistributedVirtualPortgroup, error) {
	datacenterO, _, err := getDCAndFolders(c, datacenter)
	if err != nil {
		return nil, fmt.Errorf("Could not get datacenter and  folders: %+v", err)
	}
	finder := find.NewFinder(c.Client, true)
	finder.SetDatacenter(datacenterO)
	pgPath := fmt.Sprintf("%s/%s", switchName, name)
	res, err := finder.Network(context.TODO(), pgPath)
	if err != nil {
		return nil, fmt.Errorf("Cannot get DVPG %s: %+v", pgPath, err)
	}
	castedobj, casted := res.(*object.DistributedVirtualPortgroup)
	if !casted {
		return nil, fmt.Errorf("Cannot cast %s to DVPG", pgPath)
	}
	return castedobj, nil
}

// load a DVPG and populate the passed struct with
func (p *dvs_port_group) loadDVPG(client *govmomi.Client, datacenter, switchName, name string, out *dvs_port_group) error {
	var pgmoObj mo.DistributedVirtualPortgroup
	pgObj, err := p.getVmomiDVPG(client, datacenter, switchName, name)
	// tokensSwitch := strings.Split(switchName, "/")
	// folder := strings.Join(tokensSwitch[:len(tokensSwitch)-1], "/")
	err = pgObj.Properties(
		context.TODO(),
		pgObj.Reference(),
		[]string{"config", "key", "portKeys"},
		&pgmoObj)
	if err != nil {
		return fmt.Errorf("Could not retrieve properties: %+v", err)
	}
	policy := pgmoObj.Config.Policy.GetDVPortgroupPolicy()
	out.description = pgmoObj.Config.Description
	out.name = pgmoObj.Config.Name
	out.numPorts = int(pgmoObj.Config.NumPorts)
	out.autoExpand = *pgmoObj.Config.AutoExpand
	out.pgType = pgmoObj.Config.Type
	out.policy.allowBlockOverride = policy.BlockOverrideAllowed
	out.policy.allowLivePortMoving = policy.LivePortMovingAllowed
	out.policy.allowNetworkRPOverride = *policy.NetworkResourcePoolOverrideAllowed
	out.policy.allowShapingOverride = policy.ShapingOverrideAllowed
	out.policy.allowTrafficFilterOverride = *policy.TrafficFilterOverrideAllowed
	out.policy.allowVendorConfigOverride = policy.VendorConfigOverrideAllowed
	out.policy.portConfigResetDisconnect = policy.PortConfigResetAtDisconnect
	out.portNameFormat = pgmoObj.Config.PortNameFormat
	out.switchId = switchName
	return nil
}

func (p *dvs_port_group) makeDVPGConfigSpec() types.DVPortgroupConfigSpec {
	a := types.DVPortgroupPolicy{
		BlockOverrideAllowed: p.policy.allowBlockOverride,
		LivePortMovingAllowed: p.policy.allowLivePortMoving,
		NetworkResourcePoolOverrideAllowed: &p.policy.allowNetworkRPOverride,
		PortConfigResetAtDisconnect: p.policy.portConfigResetDisconnect,
		ShapingOverrideAllowed: p.policy.allowShapingOverride,
		TrafficFilterOverrideAllowed: &p.policy.allowTrafficFilterOverride,
		VendorConfigOverrideAllowed: p.policy.allowVendorConfigOverride,
	}
	return types.DVPortgroupConfigSpec{
		AutoExpand: &p.autoExpand,
		Description: p.description,
		Name: p.name,
		NumPorts: int32(p.numPorts),
		PortNameFormat: p.portNameFormat,
		Type: "earlyBinding",
		Policy: &a,
	}
}

func (p *dvs_port_group) createPortgroup(c *govmomi.Client) error {
	createSpec := p.makeDVPGConfigSpec()
	switchID, err := parseDVSID(p.switchId) // here we get the datacenter ID aswell
	dvsObj := dvs{}

	err = loadDVS(c, switchID.datacenter, switchID.name, &dvsObj)
	if err != nil {
		return fmt.Errorf("Cannot loadDVS: %+v", err)
	}
	dvsMo, err := dvsObj.getDVS(c, dvsObj.name)
	if err != nil {
		return fmt.Errorf("Cannot getDVS: %+v", err)
	}
	task, err := createDVPortgroup(c, dvsMo, createSpec)

	_, err = task.WaitForResult(context.TODO(), nil)
	if err != nil {
		return fmt.Errorf("Could not create the DVPG: %+v", err)
	}
	return nil
}

func (p *dvs_port_group) deletePortgroup (c *govmomi.Client) error {
	log.Println("TODO: deletePortgroup → not implemented")
	return nil
}

func (p *dvs_port_group) getProperties(c *govmomi.Client) (*mo.DistributedVirtualPortgroup, error) {

	dvspgMo := mo.DistributedVirtualPortgroup{}
	switchID, err := parseDVSID(p.switchId)
	if err != nil {
		return nil, err
	}
	dvspgobj, err := p.getVmomiDVPG(c, switchID.datacenter, switchID.name, p.name)
	if err != nil {
		return nil, err
	}
	return &dvspgMo, dvspgobj.Properties(
		context.TODO(),
		dvspgobj.Reference(),
		[]string{"config", "key", "portKeys"},
		&dvspgMo)
}

// dvs_map_vm_dvpg methods

func (m *dvs_map_vm_dvs) loadMapVMDVPG(c *govmomi.Client, datacenter, switchName, portgroup, vmPath string) (out *dvs_map_vm_dvs, err error) {
	var errs []error
	vmObj, err := getVirtualMachine(c, datacenter, vmPath)
	if err != nil {
		errs = append(errs, err)
	}
	pgID, err := parseDVPGID(portgroup)
	if err != nil {
		errs = append(errs, err)
	}
	dvpgObj, err := loadDVPG(c, pgID.datacenter, pgID.switchName, pgID.name)
	if err != nil {
		errs = append(errs, err)
	}
	dvpgProps, err := dvpgObj.getProperties(c)
	if err != nil {
		errs = append(errs, err)
	}
	devs, err  := vmObj.Device(context.TODO())
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		goto EndStatement
	}
	for _, dev := range devs {
		switch dev.(type) {
		case (types.BaseVirtualEthernetCard):
			log.Printf("Veth")
			d2 := dev.(types.BaseVirtualEthernetCard)
			back, casted := d2.GetVirtualEthernetCard().Backing.(*types.VirtualEthernetCardDistributedVirtualPortBackingInfo)
			if !casted {
				errs = append(errs, fmt.Errorf("Cannot cast Veth, aborting…"))
				goto EndStatement
			}
			if back.Port.PortgroupKey == dvpgProps.Key {
				out.nicLabel = d2.GetVirtualEthernetCard().VirtualDevice.DeviceInfo.GetDescription().Label
				break
			}
		default:
			log.Printf("Type not implemented: %T %+v\n", dev, dev)
		}
	}
	out.portgroup = portgroup
	out.vm = vmPath

EndStatement:
	if len(errs) > 0 {
		return nil, fmt.Errorf("Errors in loadMapVMDVPG: %+v", errs)
	}
	return out, err
}

// functions

// object loading functions

// load a DVS
func loadDVS(c *govmomi.Client, datacenter, dvsName string, output *dvs) error {
	output.datacenter = datacenter
	err := output.loadDVS(c, datacenter, dvsName)
	return err
}

// load map between host and DVS
func loadMapHostDVS(switchName string, hostmember types.DistributedVirtualSwitchHostMember) (out *dvs_map_host_dvs, err error) {
	h := hostmember.Config.Host
	hostObj, casted := h.Value, true
	if !casted {
		err = fmt.Errorf("Could not cast Host to mo.HostSystem")
		return
	}
	backingInfosObj := hostmember.Config.Backing

	backingInfos, casted := backingInfosObj.(*types.DistributedVirtualSwitchHostMemberPnicBacking)
	if !casted {
		err = fmt.Errorf("Could not cast Host to mo.HostSystem")
		return
	}
	for _, pnic := range backingInfos.PnicSpec {
		out.nicName = append(out.nicName, pnic.PnicDevice)
	}
	out.switchName = switchName
	out.hostName = hostObj
	return
}

// load a DVPG
func loadDVPG(client *govmomi.Client, datacenter, switchName, name string) (*dvs_port_group, error) {
	dvpg := dvs_port_group{}
	err := dvpg.loadDVPG(client, datacenter, switchName, name, &dvpg)
	return &dvpg, err
}

// load a map between DVPG and VM
func loadMapVMDVPG(client *govmomi.Client, datacenter, switchName, portgroup, vmPath string) (out *dvs_map_vm_dvs, err error) {
	return out.loadMapVMDVPG(client, datacenter, switchName, portgroup, vmPath)
}

// Host manipulation functions

func getHost(c *govmomi.Client, datacenter, hostPath string) (*object.HostSystem, error) {
	dc, _, err := getDCAndFolders(c, datacenter)
	if err != nil {
		return nil, fmt.Errorf("Could not get DC and folders: %+v", err)
	}
	finder := find.NewFinder(c.Client, true)
	finder.SetDatacenter(dc)
	host, err := finder.HostSystem(context.TODO(), hostPath)
	if err != nil {
		return nil, fmt.Errorf("Cannot find HostSystem %s: %+v", hostPath, err)
	}
	return host, nil
}


// VM manipulation functions

func getVirtualMachine(c *govmomi.Client, datacenter, vmPath string) (*object.VirtualMachine, error) {
	var finder *find.Finder
	var errs []error
	var err error
	var dc *object.Datacenter
	//var folders *object.DatacenterFolders
	var vm *object.VirtualMachine
	dc, _, err = getDCAndFolders(c, datacenter)
	if err != nil {
		errs = append(errs, err)
		goto EndPosition
	}
	finder = find.NewFinder(c.Client, true)
	finder.SetDatacenter(dc)
	vm, err = finder.VirtualMachine(context.TODO(), vmPath)
	if err != nil {
		errs = append(errs, err)
		goto EndPosition
	}
	EndPosition:
	if len(errs) > 0 {
		err = fmt.Errorf("Errors in getVirtualMachine: %+v", errs)
	}
	return vm, err
}

// device manipulation functions

func getDeviceByName(c *govmomi.Client, vm *object.VirtualMachine, deviceName string) (*types.BaseVirtualDevice, error) {
	devices, err := vm.Device(context.TODO())
	if err != nil {
		return nil, err
	}
	out := devices.Find(deviceName)

	return &out, nil
}

// VEth manipulation functions

func getVEthByName(c *govmomi.Client, vm *object.VirtualMachine, deviceName string) (*types.VirtualEthernetCard, error) {
	dev, err := getDeviceByName(c, vm, deviceName)
	if err != nil {
		return nil, err
	}
	return (*dev).(types.BaseVirtualEthernetCard).GetVirtualEthernetCard(), nil
}


// bind a VEth and a portgroup → change the VEth so it is bound to one port in the portgroup.
func bindVEthAndPortgroup(c *govmomi.Client, vm *object.VirtualMachine, veth *types.VirtualEthernetCard, portgroup *dvs_port_group) error {
	// use a VirtualMachineConfigSpec.deviceChange (VirtualDeviceConfigSpec[])
	conf := types.VirtualMachineConfigSpec{}
	devChange := types.VirtualDeviceConfigSpec{}
	devChange.Operation = "edit"
	devChange.Device = veth
	cbk, casted := veth.Backing.(*types.VirtualEthernetCardDistributedVirtualPortBackingInfo)
	if !casted {
		return fmt.Errorf("Invalid backing info for %T", cbk)
	}
	properties, err := portgroup.getProperties(c)
	if err != nil {
		return fmt.Errorf("Cannot get portgroup properties: %+v", err)
	}

	cbk.Port = types.DistributedVirtualSwitchPortConnection{
		PortgroupKey: properties.Key,
	}
	conf.DeviceChange = []types.BaseVirtualDeviceConfigSpec{&devChange,}
	task, err := vm.Reconfigure(context.TODO(), conf)
	if err != nil {
		return err
	}
	return waitForTaskEnd(task, "Cannot complete vm.Reconfigure: %+v")
}
