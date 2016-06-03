package vsphere
import "github.com/hashicorp/terraform/helper/schema"

func resourceVSphereDVS() *schema.Resource {
	return &schema.Resource{
		Create: resourceVSphereDVSCreate,
		Read:   resourceVSphereDVSRead,
		Update:   resourceVSphereDVSUpdate,
		Delete:   resourceVSphereDVSDelete,
		Schema: resourceVSphereDVSSchema(),
	}
}

func resourceVSphereDVSPG() *schema.Resource {
	return &schema.Resource{
		Create: resourceVSphereDVSPGCreate,
		Read:   resourceVSphereDVSPGRead,
		Update:   resourceVSphereDVSPGUpdate,
		Delete:   resourceVSphereDVSPGDelete,
		Schema: resourceVSphereDVSPGSchema(),
	}
}

func resourceVSphereMapHostDVS() *schema.Resource {
	return &schema.Resource{
		Create: resourceVSphereMapHostDVSCreate,
		Read: resourceVSphereMapHostDVSRead,
		Update: resourceVSphereMapHostDVSUpdate,
		Delete: resourceVSphereMapHostDVSDelete,
		Schema: resourceVSphereMapHostDVSSchema(),
	}
}

func resourceVSphereMapVMDVPG() *schema.Resource {
	return &schema.Resource{
		Create: resourceVSphereMapVMDVPGCreate,
		Read: resourceVSphereMapVMDVPGRead,
		// Update: resourceVSphereMapVMDVPGUpdate, // not needed
		Delete: resourceVSphereMapVMDVPGDelete,
		Schema: resourceVSphereMapVMDVPGSchema(),
	}
}
