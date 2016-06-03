package vsphere

import "github.com/vmware/govmomi"
import "github.com/vmware/govmomi/find"
import "fmt"
import "github.com/vmware/govmomi/object"
import "golang.org/x/net/context"

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
	} else {
		d, err := finder.DefaultDatacenter(context.TODO())
		return d, err
	}
}
