package vsphere

import "github.com/hashicorp/terraform/helper/schema"

func resourceVSphereDVS() *schema.Resource {
	return &schema.Resource{
		Create: resourceVSphereDVSCreate,
		Read:   resourceVSphereDVSRead,
		Update: resourceVSphereDVSUpdate,
		Delete: resourceVSphereDVSDelete,
		Schema: resourceVSphereDVSSchema(),
	}
}

func resourceVSphereDVPG() *schema.Resource {
	return &schema.Resource{
		Create: resourceVSphereDVPGCreate,
		Read:   resourceVSphereDVPGRead,
		Update: resourceVSphereDVPGUpdate,
		Delete: resourceVSphereDVPGDelete,
		Schema: resourceVSphereDVPGSchema(),
	}
}

func resourceVSphereMapHostDVS() *schema.Resource {
	return &schema.Resource{
		Create: resourceVSphereMapHostDVSCreate,
		Read:   resourceVSphereMapHostDVSRead,
		Delete: resourceVSphereMapHostDVSDelete,
		Schema: resourceVSphereMapHostDVSSchema(),
	}
}

func resourceVSphereMapVMDVPG() *schema.Resource {
	return &schema.Resource{
		Create: resourceVSphereMapVMDVPGCreate,
		Read:   resourceVSphereMapVMDVPGRead,
		// Update: resourceVSphereMapVMDVPGUpdate, // not needed
		Delete: resourceVSphereMapVMDVPGDelete,
		Schema: resourceVSphereMapVMDVPGSchema(),
	}
}
