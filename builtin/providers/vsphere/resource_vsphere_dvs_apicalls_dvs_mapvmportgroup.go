package vsphere

import (
	"fmt"
	"log"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/types"
	"golang.org/x/net/context"
)

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
	devs, err := vmObj.Device(context.TODO())
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
