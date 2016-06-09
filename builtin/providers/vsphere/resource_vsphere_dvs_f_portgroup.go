package vsphere

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

// name format for DVPG: datacenter, switch name, name

const dvpg_name_format = "vSphere::DVPG::%s---%s---%s"

type dvPGID struct {
	datacenter string
	switchName string
	name       string
}

/* functions for DistributedVirtualPortgroup */

func resourceVSphereDVPGCreate(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] Starting DVPGCreate")
	client, err := getGovmomiClient(meta)
	if err != nil {
		return err
	}
	item := dvs_port_group{}
	err = parseDVPG(d, &item)
	if err != nil {
		return fmt.Errorf("Cannot parseDVPG %+v: %+v", d, err)
	}
	err = item.createPortgroup(client)
	if err != nil {
		return fmt.Errorf("Cannot createPortgroup: %+v", err)
	}
	switchID, err := parseDVSID(item.switchId)
	if err != nil {
		return fmt.Errorf("Could not parse DVSID %s: %+v", item.switchId, err)
	}
	d.SetId(fmt.Sprint(dvpg_name_format, switchID.datacenter, switchID.name, item.name))
	return nil
}

func resourceVSphereDVPGRead(d *schema.ResourceData, meta interface{}) error {
	var errs []error
	log.Println("[DEBUG] Starting DVPGRead")
	client, err := getGovmomiClient(meta)
	if err != nil {
		errs = append(errs, err)
	}
	log.Printf("[DEBUG] Client: %+v", client)
	// load the state from vSphere and provide the hydrated object.
	resourceID, err := parseDVPGID(d.Id())
	if err != nil {
		errs = append(errs, fmt.Errorf("Cannot parse DVSPGID… %+v", err))
	}
	if len(errs) > 0 {
		return fmt.Errorf("There are errors in DVPGRead. Cannot proceed.\n%+v", errs)
	}
	dvspgObject := dvs_port_group{}
	err = dvspgObject.loadDVPG(client, resourceID.datacenter, resourceID.switchName, resourceID.name, &dvspgObject)
	if err != nil {
		errs = append(errs, fmt.Errorf("Cannot read DVPG %+v: %+v", resourceID, err))
	}
	if len(errs) > 0 { // we cannot load the DVPG for a reason
		log.Printf("[ERROR] Cannot load DVPG %+v", resourceID)
		return fmt.Errorf("Errors in DVPGRead: %+v", errs)
	}
	return nil
}

func resourceVSphereDVPGUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
	/*
		// now populate the object
		if err:=unparseDVPG(d, &dvspgObject); err != nil {
			log.Printf("[ERROR] Cannot populate DVPG: %+v", err)
			return err
		}
	*/
}

func resourceVSphereDVPGDelete(d *schema.ResourceData, meta interface{}) error {
	//var errs []error
	log.Println("[DEBUG] Starting DVPGDelete")
	/*client, err := getGovmomiClient(meta)
	if err != nil {
		errs = append(errs, err)
	}
	*/
	// use Destroy_Task
	d.SetId("")
	return nil
}

// parse a DVPG ResourceData to a dvs_port_group struct
func parseDVPG(d *schema.ResourceData, out *dvs_port_group) error {
	o := out
	if v, ok := d.GetOk("name"); ok {
		o.name = v.(string)
	}
	if v, ok := d.GetOk("switch_id"); ok {
		o.switchId = v.(string)
	}
	if v, ok := d.GetOk("description"); ok {
		o.description = v.(string)
	}
	if v, ok := d.GetOk("auto_expand"); ok {
		o.autoExpand = v.(bool)
	}
	if v, ok := d.GetOk("num_ports"); ok {
		o.numPorts = v.(int)
	}
	if v, ok := d.GetOk("port_name_format"); ok {
		o.portNameFormat = v.(string)
	}
	if s, ok := d.GetOk("policy"); ok {
		vmap, casted := s.(map[string]interface{})
		if !casted {
			return fmt.Errorf("Cannot cast policy as a string map. See: %+v", s)
		}
		o.policy.allowBlockOverride = vmap["allow_block_override"].(bool)
		o.policy.allowLivePortMoving = vmap["allow_live_port_moving"].(bool)
		o.policy.allowNetworkRPOverride = vmap["allow_network_rp_override"].(bool)
		o.policy.portConfigResetDisconnect = vmap["port_config_reset_disconnect"].(bool)
		o.policy.allowShapingOverride = vmap["allow_shaping_override"].(bool)
		o.policy.allowTrafficFilterOverride = vmap["allow_traffic_filter_override"].(bool)
		o.policy.allowVendorConfigOverride = vmap["allow_vendor_config_override"].(bool)
	}
	return nil
}

// fill a ResourceData using the provided DVPG
func unparseDVPG(d *schema.ResourceData, in *dvs_port_group) error {
	var errs []error
	// define the contents - this means map the stuff to what Terraform expects
	fieldsMap := map[string]interface{}{
		"name":             in.name,
		"switch_id":        in.switchId,
		"description":      in.description,
		"auto_expand":      in.autoExpand,
		"num_ports":        in.numPorts,
		"port_name_format": in.portNameFormat,
		"policy": map[string]bool{
			"allow_block_override":          in.policy.allowBlockOverride,
			"allow_live_port_moving":        in.policy.allowLivePortMoving,
			"allow_network_rp_override":     in.policy.allowNetworkRPOverride,
			"port_config_reset_disconnect":  in.policy.portConfigResetDisconnect,
			"allow_shaping_override":        in.policy.allowShapingOverride,
			"allow_traffic_filter_override": in.policy.allowTrafficFilterOverride,
			"allow_vendor_config_override":  in.policy.allowVendorConfigOverride,
		},
	}
	// set values
	for fieldName, fieldValue := range fieldsMap {
		if err := d.Set(fieldName, fieldValue); err != nil {
			errs = append(errs, fmt.Errorf("%s invalid: %s", fieldName, fieldValue))
		}
	}
	// handle errors
	if len(errs) > 0 {
		return fmt.Errorf("Errors in unparseDVPG: invalid resource definition!\n%+v", errs)
	}
	return nil
}
