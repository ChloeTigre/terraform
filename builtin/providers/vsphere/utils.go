package vsphere

import (
	"fmt"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"golang.org/x/net/context"
)

func getGovmomiClient(meta interface{}) (*govmomi.Client, error) {
	client, casted := meta.(*govmomi.Client)
	if !casted {
		return nil, fmt.Errorf("%+v is not castable as govmomi.Client", meta)
	}
	return client, nil
}

// getDatacenter gets datacenter object
func getDatacenter(c *govmomi.Client, dc string) (*object.Datacenter, error) {
	finder := find.NewFinder(c.Client, true)
	if dc != "" {
		d, err := finder.Datacenter(context.TODO(), dc)
		return d, err
	}
	d, err := finder.DefaultDatacenter(context.TODO())
	return d, err
}

// parse ID to components (DVS)
func parseDVSID(id string) (out *dvsID, err error) {
	_, err = fmt.Sscanf(id, dvs_name_format, out.datacenter, out.name)
	return
}

// parse ID to components (DVPG)
func parseDVPGID(id string) (out *dvPGID, err error) {
	_, err = fmt.Sscanf(id, dvpg_name_format, out.datacenter, out.switchName, out.name)
	return
}

// parse ID to components (MapHostDVS)
func parseMapHostDVSID(id string) (out *mapHostDVSID, err error) {
	_, err = fmt.Sscanf(id, maphostdvs_name_format, out.datacenter, out.switchName, out.hostName)
	return
}

// parse ID to components (MapHostDVS)
func parseMapVMDVPGID(id string) (out *mapVMDVPGID, err error) {
	_, err = fmt.Sscanf(id, mapvmdvpg_name_format, out.datacenter, out.switchName, out.portgroupName, out.vmName)
	return
}

// wait for task end
func waitForTaskEnd(task *object.Task, message string) error {
	if _, err := task.WaitForResult(context.TODO(), nil); err != nil {
		return fmt.Errorf(message, err)
	}
	return nil

}

func getDCAndFolders(c *govmomi.Client, datacenter string) (*object.Datacenter, *object.DatacenterFolders, error) {
	dvso := dvs{
		datacenter: datacenter,
	}
	return dvso.getDCAndFolders(c)
}
