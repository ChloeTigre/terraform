package vsphere

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

// name format for MapVMDVPG: datacenter, switch name, port name, vm name

type mapVMDVPGID struct {
	datacenter    string
	switchName    string
	portgroupName string
	vmName        string
}

/* Functions for MapVMDVS */

func resourceVSphereMapVMDVPGCreate(d *schema.ResourceData, meta interface{}) error {
	var errs []error
	var err error
	var params *dvs_map_vm_dvs
	// start by getting the DVS
	client, err := getGovmomiClient(meta)
	if err != nil {
		errs = append(errs, err)
	}
	params, err = parseMapVMDVS(d)
	if err != nil {
		errs = append(errs, err)
	}
	portgroupID, err := parseDVPGID(params.portgroup)
	if err != nil {
		errs = append(errs, err)
	}
	// get VM and NIC
	vm, err := getVirtualMachine(client, portgroupID.datacenter, params.vm)
	if err != nil {
		errs = append(errs, err)
	}
	veth, err := getVEthByName(client, vm, params.nicLabel)
	if err != nil {
		errs = append(errs, err)
	}
	portgroup, err := loadDVPG(client, portgroupID.datacenter, portgroupID.switchName, portgroupID.name)
	if err != nil {
		errs = append(errs, err)
	}

	// update backing informations of the VEth so it connects to the Portgroup
	err = bindVEthAndPortgroup(client, vm, veth, portgroup)
	if err != nil {
		errs = append(errs, err)
	}
	// end
	if len(errs) > 0 {
		return fmt.Errorf("Errors in MapVMDVPG.Create: %+v", errs)
	}
	d.SetId(params.getID())
	return nil
}

func resourceVSphereMapVMDVPGRead(d *schema.ResourceData, meta interface{}) error {
	// read the state of said DVPG using just its Id and set the d object
	// values accordingly

	var errs []error
	log.Println("[DEBUG] Starting MapVMDVPG")
	client, err := getGovmomiClient(meta)
	if err != nil {
		errs = append(errs, err)
	}
	// load the state from vSphere and provide the hydrated object.
	idObj, err := parseMapVMDVPGID(d.Id())
	if err != nil {
		errs = append(errs, fmt.Errorf("Cannot parse MapVMDVSPGID… %+v", err))
	}
	if len(errs) > 0 {
		return fmt.Errorf("There are errors in MapVMDVPGRead. Cannot proceed.\n%+v", errs)
	}

	mapdvspgObject, err := loadMapVMDVPG(client, idObj.datacenter, idObj.switchName, idObj.portgroupName, idObj.vmName)
	if err != nil {
		errs = append(errs, fmt.Errorf("Cannot load MapVMDVPG %+v: %+v", err, err))
	}
	if len(errs) > 0 { // we cannot load the DVPG for a reason
		log.Printf("[ERROR] Cannot load MapVMDVPG %+v", mapdvspgObject)
		return fmt.Errorf("Errors in MapVMDVPGRead: %+v", errs)
	}
	// now just populate the ResourceData
	return unparseMapHostDVPG(d, mapdvspgObject)
}

func resourceVSphereMapVMDVPGUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceVSphereMapVMDVPGDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}

func (d *dvs_map_vm_dvs) getID() string {
	portgroupID, err := parseDVPGID(d.portgroup)
	if err != nil {
		return "!!ERROR!!"
	}
	return fmt.Sprintf(
		mapvmdvpg_name_format, portgroupID.datacenter,
		portgroupID.switchName, portgroupID.name, d.vm)
}

// take a dvs_map_vm_dvs and put its contents into the ResourceData.
func unparseMapHostDVPG(d *schema.ResourceData, in *dvs_map_vm_dvs) error {
	var errs []error
	fieldsMap := map[string]interface{}{
		"nic_label": in.nicLabel,
		"portgroup": in.portgroup,
		"vm":        in.vm,
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
