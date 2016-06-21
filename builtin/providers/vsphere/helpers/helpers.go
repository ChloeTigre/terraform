package helpers

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/mo"
	"golang.org/x/net/context"
)

// GetDatacenter gets datacenter object - meant for internal use
func GetDatacenter(c *govmomi.Client, dc string) (*object.Datacenter, error) {
	finder := find.NewFinder(c.Client, true)
	if dc != "" {
		d, err := finder.Datacenter(context.TODO(), dc)
		return d, err
	}
	d, err := finder.DefaultDatacenter(context.TODO())
	return d, err
}

// WaitForTaskEnd waits for a vSphere task to end
func WaitForTaskEnd(task *object.Task, message string) error {
	//time.Sleep(time.Second * 5)
	if err := task.Wait(context.TODO()); err != nil {
		spew.Dump("Error in waitForTaskEnd", err)

		taskmo := mo.Task{}
		task.Properties(context.TODO(), task.Reference(), []string{"info"}, &taskmo)
		spew.Dump("Task", taskmo)
		return fmt.Errorf("[%T] → "+message, err, err)
	}
	return nil

}
