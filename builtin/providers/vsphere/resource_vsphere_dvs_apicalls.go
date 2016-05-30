package vsphere

import (
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/types"
	"golang.org/x/net/context"
	"fmt"
)

func (d *dvs) makeDVSConfigSpec() (types.DVSCreateSpec, error) {
	return types.DVSCreateSpec{
		ConfigSpec: &types.DVSConfigSpec{
			Contact: &types.DVSContactInfo{
				Contact: d.contact.infos,
				Name: d.contact.name,
			},
			ExtensionKey: d.key,
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

func (d *dvs) createSwitch(c *govmomi.Client) error {
	datacenter, err := getDatacenter(c, d.datacenter)
	if err != nil {
		return fmt.Errorf("Cannot get datacenter from %+v [%+v]", d, err)
	}

	// get Network Folder from datacenter
	dcFolders, err := datacenter.Folders(context.TODO())
	if err != nil {
		return fmt.Errorf("Cannot get folders for datacenter %+v [%+v]", datacenter, err)
	}
	folder := dcFolders.NetworkFolder

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
