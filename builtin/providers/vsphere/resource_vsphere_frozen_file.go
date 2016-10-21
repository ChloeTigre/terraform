package vsphere

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/soap"
	"golang.org/x/net/context"
)

type frozenFile struct {
	sourceDatacenter  string
	datacenter        string
	sourceDatastore   string
	datastore         string
	sourceFile        string
	destinationFile   string
	createDirectories bool
	copyFile          bool
	isDiskFile        bool
}

func ResourceVSphereFrozenFile() *schema.Resource {
	return &schema.Resource{
		Create: resourceVsphereFrozenFileCreate,
		Read:   resourceVsphereFrozenFileRead,
		Update: resourceVsphereFrozenFileUpdate,
		Delete: resourceVsphereFrozenFileDelete,

		Schema: map[string]*schema.Schema{
			"datacenter": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"source_datacenter": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"datastore": {
				Type:     schema.TypeString,
				Required: true,
			},

			"source_datastore": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"source_file": {
				Type:     schema.TypeString,
				Required: true,
			},

			"destination_file": {
				Type:     schema.TypeString,
				Required: true,
			},

			"create_directories": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"is_disk": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceVsphereFrozenFileCreate(d *schema.ResourceData, meta interface{}) error {

	log.Printf("[DEBUG] creating file: %#v", d)
	client := meta.(*govmomi.Client)

	f := frozenFile{}

	if v, ok := d.GetOk("source_datacenter"); ok {
		f.sourceDatacenter = v.(string)
		f.copyFile = true
	}

	if v, ok := d.GetOk("datacenter"); ok {
		f.datacenter = v.(string)
	}

	if v, ok := d.GetOk("source_datastore"); ok {
		f.sourceDatastore = v.(string)
		f.copyFile = true
	}

	if v, ok := d.GetOk("datastore"); ok {
		f.datastore = v.(string)
	} else {
		return fmt.Errorf("datastore argument is required")
	}

	if v, ok := d.GetOk("source_file"); ok {
		f.sourceFile = v.(string)
	} else {
		return fmt.Errorf("source_file argument is required")
	}

	if v, ok := d.GetOk("destination_file"); ok {
		f.destinationFile = v.(string)
	} else {
		return fmt.Errorf("destination_file argument is required")
	}

	if v, ok := d.GetOk("create_directories"); ok {
		f.createDirectories = v.(bool)
	}

	if v, ok := d.GetOk("is_disk"); ok {
		f.isDiskFile = v.(bool)
	}

	err := createFrozenFile(client, &f)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("[%v] %v/%v", f.datastore, f.datacenter, f.destinationFile))
	log.Printf("[INFO] Created file: %s", f.destinationFile)

	return resourceVsphereFrozenFileRead(d, meta)
}

func createFrozenFile(client *govmomi.Client, f *frozenFile) error {
	var err error
	finder := find.NewFinder(client.Client, true)

	dc, err := finder.Datacenter(context.TODO(), f.datacenter)
	if err != nil {
		return fmt.Errorf("error %s", err)
	}
	finder = finder.SetDatacenter(dc)

	ds, err := getDatastore(finder, f.datastore)
	if err != nil {
		return fmt.Errorf("error %s", err)
	}

	if f.copyFile {
		// Copying frozenFile from withing vSphere
		source_dc, err := finder.Datacenter(context.TODO(), f.sourceDatacenter)
		if err != nil {
			return fmt.Errorf("error %s", err)
		}
		finder = finder.SetDatacenter(dc)

		source_ds, err := getDatastore(finder, f.sourceDatastore)
		if err != nil {
			return fmt.Errorf("error %s", err)
		}

		fm := object.NewFileManager(client.Client)
		dm := object.NewVirtualDiskManager(client.Client)
		if f.createDirectories {
			directoryPathIndex := strings.LastIndex(f.destinationFile, "/")
			path := f.destinationFile[0:directoryPathIndex]
			err = fm.MakeDirectory(context.TODO(), ds.Path(path), dc, true)
			if err != nil {
				return fmt.Errorf("error %s", err)
			}
		}
		var task *object.Task
		if !f.isDiskFile {
			task, err = fm.CopyDatastoreFile(context.TODO(), source_ds.Path(f.sourceFile), source_dc, ds.Path(f.destinationFile), dc, true)
		} else {
			task, err = dm.CopyVirtualDisk(context.TODO(), source_ds.Path(f.sourceFile), source_dc, ds.Path(f.destinationFile), dc, nil, true)
		}

		if err != nil {
			return fmt.Errorf("error %s", err)
		}

		// strange race condition occurs sometimes
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[DEBUG] recovered from task.WaitForResult for copyFile")
			}
		}()
		_, err = task.WaitForResult(context.TODO(), nil)
		if err != nil {
			return fmt.Errorf("error %s", err)
		}

	} else {
		// Uploading frozenFile to vSphere
		dsurl, err := ds.URL(context.TODO(), dc, f.destinationFile)
		if err != nil {
			return fmt.Errorf("error %s", err)
		}

		p := soap.DefaultUpload
		err = client.Client.UploadFile(f.sourceFile, dsurl, &p)
		if err != nil {
			return fmt.Errorf("error %s", err)
		}
	}

	return nil
}

func resourceVsphereFrozenFileRead(d *schema.ResourceData, meta interface{}) error {

	log.Printf("[DEBUG] reading file: %#v", d)
	f := frozenFile{}

	if v, ok := d.GetOk("source_datacenter"); ok {
		f.sourceDatacenter = v.(string)
	}

	if v, ok := d.GetOk("datacenter"); ok {
		f.datacenter = v.(string)
	}

	if v, ok := d.GetOk("source_datastore"); ok {
		f.sourceDatastore = v.(string)
	}

	if v, ok := d.GetOk("datastore"); ok {
		f.datastore = v.(string)
	} else {
		return fmt.Errorf("datastore argument is required")
	}

	if v, ok := d.GetOk("source_file"); ok {
		f.sourceFile = v.(string)
	} else {
		return fmt.Errorf("source_file argument is required")
	}

	if v, ok := d.GetOk("destination_file"); ok {
		f.destinationFile = v.(string)
	} else {
		return fmt.Errorf("destination_file argument is required")
	}

	log.Printf("[INFO] Not updating frozen file")
	return nil
}

func resourceVsphereFrozenFileUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Not updating frozen file")
	return nil
}

func resourceVsphereFrozenFileDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] deleting file: %#v", d)
	f := frozenFile{}

	if v, ok := d.GetOk("datacenter"); ok {
		f.datacenter = v.(string)
	}

	if v, ok := d.GetOk("datastore"); ok {
		f.datastore = v.(string)
	} else {
		d.SetId("")
		log.Printf("[ERROR] datastore argument is required")
		return nil
	}

	if v, ok := d.GetOk("source_file"); ok {
		f.sourceFile = v.(string)
	} else {
		d.SetId("")
		log.Printf("[ERROR] source_file argument is required")
		return nil
	}

	if v, ok := d.GetOk("destination_file"); ok {
		f.destinationFile = v.(string)
	} else {
		d.SetId("")
		log.Printf("[ERROR] destination_file argument is required")
		return nil
	}

	client := meta.(*govmomi.Client)

	err := deleteFrozenFile(client, &f)
	if err != nil {
		log.Printf("[ERROR] deleteFile: %[1]T %#[1]v", err)
	}

	d.SetId("")
	return nil
}

func deleteFrozenFile(client *govmomi.Client, f *frozenFile) error {

	dc, err := getDatacenter(client, f.datacenter)
	if err != nil {
		return err
	}

	finder := find.NewFinder(client.Client, true)
	finder = finder.SetDatacenter(dc)

	ds, err := getDatastore(finder, f.datastore)
	if err != nil {
		return fmt.Errorf("error %s", err)
	}

	fm := object.NewFileManager(client.Client)
	task, err := fm.DeleteDatastoreFile(context.TODO(), ds.Path(f.destinationFile), dc)
	if err != nil {
		log.Printf("[ERROR] Could not delete file.")
		return err

	}

	_, err = task.WaitForResult(context.TODO(), nil)
	if err != nil {
		return err
	}
	return nil
}
