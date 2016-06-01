package vsphere

import (
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
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

func switchBoilerplate(c *govmomi.Client) (*object.Datacenter, *object.DatacenterFolders, error) {
	datacenter, err := getDatacenter(c, d.datacenter)
	if err != nil {
		return nil, nil, fmt.Errorf("Cannot get datacenter from %+v [%+v]", d, err)
	}

	// get Network Folder from datacenter
	dcFolders, err := datacenter.Folders(context.TODO())
	if err != nil {
		return nil, nil, fmt.Errorf("Cannot get folders for datacenter %+v [%+v]", datacenter, err)
	}
	folder := dcFolders.NetworkFolder
	return datacenter, dcFolders, nil
}


func (d *dvs) createSwitch(c *govmomi.Client) error {
	_, folders, err := switchBoilerplate(c)
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
func (d *dvs) getDVS(c *govmomi.Client, dvsName string) error {
	datacenter, folders, err := switchBoilerplate(c)
	if err != nil {
		return fmt.Errorf("Could not get datacenter and  folders: %+v", err)
	}
	folder := folders.NetworkFolder
	folder.

}
