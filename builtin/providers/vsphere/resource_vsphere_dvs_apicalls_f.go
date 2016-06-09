package vsphere

import (
	"fmt"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	"golang.org/x/net/context"
)

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
	conf.DeviceChange = []types.BaseVirtualDeviceConfigSpec{&devChange}
	task, err := vm.Reconfigure(context.TODO(), conf)
	if err != nil {
		return err
	}
	return waitForTaskEnd(task, "Cannot complete vm.Reconfigure: %+v")
}
