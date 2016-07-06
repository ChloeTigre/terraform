package helpers

import (
	"fmt"

	"errors"

	"github.com/davecgh/go-spew/spew"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
	"golang.org/x/net/context"
)

// TaskFailedError is an error triggered when a Task is not successful.
var TaskFailedError error

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
		return fmt.Errorf("[%T] "+message, err, err)
	}
	taskmo := mo.Task{}
	task.Properties(context.TODO(), task.Reference(), []string{"info"}, &taskmo)
	spew.Dump("Task", taskmo)
	if taskmo.Info.State != types.TaskInfoStateSuccess {
		return TaskFailedError
	}
	return nil

}

func init() {
	TaskFailedError = errors.New("Task Failed")
}
