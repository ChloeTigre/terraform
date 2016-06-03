package vsphere

import "github.com/hashicorp/terraform/helper/schema"
import "log"
import "fmt"

// name format for DVS: datacenter, name
const dvs_name_format = "DVS:[%s] %s"
// name format for DVPG: datacenter, switch name, name
const dvpg_name_format = "DVPG:[%s.%s] %s"
// name format for MapHostDVS: datacenter, DVS name, Host name
const maphostdvs_name_format = "MapHostDVS:[%s] %s-%s"
// name format for MapVMDVPG: datacenter, switch name, port name, vm name
const mapvmdvpg_name_format = "MapVMDVPG:[%s] %s.%s-%s"

/* functions for DistributedVirtualSwitch */
func resourceVSphereDVSCreate(d *schema.ResourceData, meta interface{}) error {
	// this creates a DVS
	log.Printf("[DEBUG] Starting DVSCreate")
	client, err := getGovmomiClient(meta)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Client: %+v", client)
	item, err := parseDVS(d)
	if err != nil {
		return fmt.Errorf("Cannot parseDVS %+v: %+v", d, err)
	}
	err = item.createSwitch(client)
	if err != nil {
		return fmt.Errorf("Cannot createSwitch: %+v", err)
	}
	d.SetId(fmt.Sprint(dvs_name_format, item.datacenter, item.name))
	return nil
}

func resourceVSphereDVSRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Starting DVSRead")
	client, err := getGovmomiClient(meta)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Client: %+v", client)
	// load the state from vSphere and provide the hydrated object.
	return nil
}

func resourceVSphereDVSUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Starting DVSUpdate")
	client, err := getGovmomiClient(meta)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Client: %+v", client)
	// detect the different changes in the object and perform needed updates

	return nil
}

func resourceVSphereDVSDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Starting DVSDelete")
	client, err := getGovmomiClient(meta)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Client: %+v", client)
	// remove the object and its dependencies in vSphere

	// then remove object from the datastore.
	d.SetId("")
	return nil
}

/* functions for DistributedVirtualPortgroup */

func resourceVSphereDVSPGCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceVSphereDVSPGRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceVSphereDVSPGUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceVSphereDVSPGDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}

/* functions for MapHostDVS */

func resourceVSphereMapHostDVSCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceVSphereMapHostDVSRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceVSphereMapHostDVSUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceVSphereMapHostDVSDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}

/* Functions for MapVMDVS */

func resourceVSphereMapVMDVPGCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceVSphereMapVMDVPGRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceVSphereMapVMDVPGUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceVSphereMapVMDVPGDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}
