package vsphere

import (
	"fmt"
	"log"
	"strings"

	"github.com/davecgh/go-spew/spew" // debug dependency
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
	"golang.org/x/net/context"
)

// functions

// object loading functions

// load a DVS
func loadDVS(c *govmomi.Client, datacenter, dvsPath string, output *dvs) error {
	output.datacenter = datacenter
	err := output.loadDVS(c, datacenter, dvsPath)
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
func loadMapVMDVPG(client *govmomi.Client, datacenter, switchName, portgroup, vmPath string) (out *dvs_map_vm_dvpg, err error) {
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
	if out == nil {
		return nil, fmt.Errorf("Could not get device named %v\n", deviceName)
	}

	return &out, nil
}

// VEth manipulation functions

func getVEthByName(c *govmomi.Client, vm *object.VirtualMachine, deviceName string) (*types.VirtualEthernetCard, string, error) {
	dev, err := getDeviceByName(c, vm, deviceName)
	if err != nil {
		return nil, "", err
	}
	if dev == nil {
		return nil, "", fmt.Errorf("Cannot return VEth: %T:%+v", err, err)
	}
	vc := (*dev).(types.BaseVirtualEthernetCard).GetVirtualEthernetCard()
	if vmx2, casted := (*dev).(*types.VirtualVmxnet2); casted {
		return vmx2.GetVirtualEthernetCard(), "vmx2", nil
	} else if vmx3, casted := (*dev).(*types.VirtualVmxnet3); casted {
		return vmx3.GetVirtualEthernetCard(), "vmx3", nil
	} else if vmx, casted := (*dev).(*types.VirtualVmxnet); casted {
		return vmx.GetVirtualEthernetCard(), "vmx", nil
	} else if e1000e, casted := (*dev).(*types.VirtualE1000e); casted {
		return e1000e.GetVirtualEthernetCard(), "e1000e", nil
	} else if e1000, casted := (*dev).(*types.VirtualE1000); casted {
		return e1000.GetVirtualEthernetCard(), "e1000", nil
	} else if pcnet, casted := (*dev).(*types.VirtualPCNet32); casted {
		return pcnet.GetVirtualEthernetCard(), "pcnet32", nil
	} else if sriov, casted := (*dev).(*types.VirtualSriovEthernetCard); casted {
		return sriov.GetVirtualEthernetCard(), "sriov", nil
	}
	return vc, "unknown", nil
}

func buildVEthDeviceChange(c *govmomi.Client, veth *types.VirtualEthernetCard, portgroup *dvs_port_group, optype types.VirtualDeviceConfigSpecOperation) (*types.VirtualDeviceConfigSpec, error) {
	devChange := types.VirtualDeviceConfigSpec{}
	cbk := types.VirtualEthernetCardDistributedVirtualPortBackingInfo{}

	properties, err := portgroup.getProperties(c)
	if err != nil {
		return nil, fmt.Errorf("Cannot get portgroup properties: %+v", err)
	}
	switchID, err := parseDVSID(portgroup.switchId)
	if err != nil {
		return nil, fmt.Errorf("Cannot parse switchID: %+v", err)
	}
	dvsObj := dvs{}
	err = loadDVS(c, switchID.datacenter, switchID.path, &dvsObj)
	if err != nil {
		return nil, fmt.Errorf("Cannot get switch: %+v", err)
	}
	dvsProps, err := dvsObj.getProperties(c)
	if err != nil {
		return nil, fmt.Errorf("Cannot get dvs properties: %+v", err)
	}

	cbk.Port = types.DistributedVirtualSwitchPortConnection{
		PortgroupKey: properties.Key,
		// PortKey:    "372",
		SwitchUuid: dvsProps.Uuid,
	}
	log.Printf("\n\n[DEBUG] cbk.Port: %+v\nswitchName: %s\nswitchUuid: %s; \nproperties: %+v\n\n", cbk, dvsProps.Name, dvsProps.Uuid, properties)
	veth2 := types.VirtualEthernetCard{}
	veth2.Key = veth.Key
	veth2.Backing = &cbk
	devChange.Operation = optype // `should be add, remove or edit`
	devChange.Device = &veth2
	return &devChange, nil
}

// bind a VEth and a portgroup → change the VEth so it is bound to one port in the portgroup.
func bindVEthAndPortgroup(c *govmomi.Client, vm *object.VirtualMachine, veth *types.VirtualEthernetCard, portgroup *dvs_port_group) error {
	// use a VirtualMachineConfigSpec.deviceChange (VirtualDeviceConfigSpec[])
	conf := types.VirtualMachineConfigSpec{}
	devspec, err := buildVEthDeviceChange(c, veth, portgroup, types.VirtualDeviceConfigSpecOperationEdit)
	if err != nil {
		return err
	}
	//devspec.Device.GetVirtualDevice().Connectable.Connected = true

	conf.DeviceChange = []types.BaseVirtualDeviceConfigSpec{devspec}

	log.Printf("\n\n\nHere comes the debug\n")
	spew.Dump("conf for reconfigure", conf)
	spew.Dump("VM to be reconfigured", vm)
	task, err := vm.Reconfigure(context.TODO(), conf)
	if err != nil {
		spew.Dump("Error\n\n", err, "\n\n")
		return err
	}
	return waitForTaskEnd(task, "Cannot complete vm.Reconfigure: %+v")
}

func unbindVEthAndPortgroup(c *govmomi.Client, vm *object.VirtualMachine, veth *types.VirtualEthernetCard, portgroup *dvs_port_group) error {
	// use a VirtualMachineConfigSpec.deviceChange (VirtualDeviceConfigSpec[])
	conf := types.VirtualMachineConfigSpec{}
	devspec, err := buildVEthDeviceChange(c, veth, portgroup, types.VirtualDeviceConfigSpecOperationEdit)
	if err != nil {
		return err
	}
	//devspec.Device.GetVirtualDevice().Connectable.Connected = false

	conf.DeviceChange = []types.BaseVirtualDeviceConfigSpec{devspec}

	task, err := vm.Reconfigure(context.TODO(), conf)
	if err != nil {
		return err
	}
	return waitForTaskEnd(task, "Cannot complete vm.Reconfigure: %+v")
}

func vmDebug(c *govmomi.Client, vm *object.VirtualMachine) {
	return
	log.Println(strings.Repeat("*", 80))
	props := mo.VirtualMachine{}
	err := vm.Properties(context.TODO(), vm.Reference(), []string{"config"}, &props)
	if err != nil {
		log.Printf("\n\n[DEBUG] Cannot get VM properties: [%+v]\n\n", err)
	}
	log.Printf("\n\nConfig: %+v\n\n", props.Config)
	veth, vtype, err := getVEthByName(c, vm, "ethernet-0")
	if err != nil {
		log.Printf("\n\n[DEBUG] Cannot get ethernet-0: [%+v]\n\n", err)
	}
	log.Printf("Ethernet type: %s", vtype)
	bk := veth.Backing.(*types.VirtualEthernetCardDistributedVirtualPortBackingInfo)
	log.Printf("[DEBUG] Backing: %+v", bk)
	log.Println(strings.Repeat("-", 80))

}
