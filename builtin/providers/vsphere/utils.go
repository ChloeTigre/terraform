package vsphere

import "github.com/vmware/govmomi"
import "fmt"

func getGovmomiClient(meta interface{}) (*govmomi.Client, error) {
       client, casted := meta.(*govmomi.Client)
       if !casted {
               return nil, fmt.Errorf("%+v is not castable as govmomi.Client", meta)
       }
       return client, nil
}
