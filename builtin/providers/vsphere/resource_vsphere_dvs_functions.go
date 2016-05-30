package vsphere

import "github.com/hashicorp/terraform/helper/schema"
import "log"
import "fmt"

const dvs_name_format = "DVS:[%s] %s"

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
	return nil
}

func resourceVSphereDVSUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Starting DVSUpdate")
	client, err := getGovmomiClient(meta)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Client: %+v", client)
	return nil
}

func resourceVSphereDVSDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Starting DVSDelete")
	client, err := getGovmomiClient(meta)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Client: %+v", client)
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

func resourceVSphereMapHostDVSPGRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceVSphereMapHostDVSPGUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceVSphereMapHostDVSPGDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}

/* Functions for MapVMDVS */

func resourceVSphereMapVMDVPGCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceVSphereMapVMDVPGPGRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceVSphereMapVMDVPGPGUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceVSphereMapVMDVPGPGDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}
