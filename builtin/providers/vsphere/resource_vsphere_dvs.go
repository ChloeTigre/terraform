package vsphere

import "github.com/hashicorp/terraform/helper/schema"

type dvs struct {
	name         string
	folder       string
	datacenter   string
	extensionKey string
	description  string
	contact      struct {
		name  string
		infos string
	}
	switchUsagePolicy struct {
		autoPreinstallAllowed bool
		autoUpgradeAllowed    bool
		partialUpgradeAllowed bool
	}
	switchIPAddress    string
	numStandalonePorts int
}

type dvs_map_host_dvs struct {
	hostName   string
	switchName string
	nicName    []string
}

type dvs_port_group struct {
	name           string
	switchId       string
	pgType         string
	description    string
	autoExpand     bool
	numPorts       int
	portNameFormat string
	policy         struct {
		allowBlockOverride         bool
		allowLivePortMoving        bool
		allowNetworkRPOverride     bool
		portConfigResetDisconnect  bool
		allowShapingOverride       bool
		allowTrafficFilterOverride bool
		allowVendorConfigOverride  bool
	}
}

type dvs_map_vm_dvs struct {
	vm        string
	nicLabel  string
	portgroup string
}

func resourceVSphereDVSSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			// ForceNew:		true,
		},
		"extension_key": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"contact": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"infos": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		"switch_usage_policy": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"auto_preinstall_allowed": &schema.Schema{
						Type:     schema.TypeBool,
						Optional: true,
					},
					"auto_upgrade_allowed": &schema.Schema{
						Type:     schema.TypeBool,
						Optional: true,
					},
					"partial_upgrade_allowed": &schema.Schema{
						Type:     schema.TypeBool,
						Optional: true,
					},
				},
			},
		},
		"switch_ip_address": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"num_standalone_ports": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
	}
}

/* functions for DistributedVirtualPortgroup */
func resourceVSphereDVPGSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"switch_id": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"datacenter": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"type": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"auto_expand": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
		"num_ports": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		"port_name_format": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"policy": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"allow_block_override": &schema.Schema{
						Type:     schema.TypeBool,
						Optional: true,
					},
					"allow_live_port_moving": &schema.Schema{
						Type:     schema.TypeBool,
						Optional: true,
					},
					"allov_network_resources_pool_override": &schema.Schema{
						Type:     schema.TypeBool,
						Optional: true,
					},
					"port_config_reset_disconnect": &schema.Schema{
						Type:     schema.TypeBool,
						Optional: true,
					},
					"allow_shaping_override": &schema.Schema{
						Type:     schema.TypeBool,
						Optional: true,
					},
					"allow_traffic_filter_override": &schema.Schema{
						Type:     schema.TypeBool,
						Optional: true,
					},
					"allow_vendor_config_override": &schema.Schema{
						Type:     schema.TypeBool,
						Optional: true,
					},
				},
			},
		},
	}
}

/* functions for MapHostDVS */

/* MapHostDVS functions */
func resourceVSphereMapHostDVSSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"host": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"switch": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"nic_names": &schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			Computed: true,
			ForceNew: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Set:      schema.HashString,
		},
	}
}

/* parse a MapHostDVS to its struct */
func parseMapHostDVS(d *schema.ResourceData) (*dvs_map_host_dvs, error) {
	o := dvs_map_host_dvs{}
	if v, ok := d.GetOk("host"); ok {
		o.hostName = v.(string)
	}
	if v, ok := d.GetOk("switch"); ok {
		o.switchName = v.(string)
	}
	if v, ok := d.GetOk("nic_names"); ok {
		o.nicName = v.([]string)
	}
	return &o, nil
}

/* Functions for MapVMDVS */

func resourceVSphereMapVMDVPGSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"vm": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"nic_label": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"portgroup": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
	}
}

/* parse a MapVMDVS to its struct */
func parseMapVMDVS(d *schema.ResourceData) (*dvs_map_vm_dvs, error) {
	o := dvs_map_vm_dvs{}
	if v, ok := d.GetOk("vm"); ok {
		o.vm = v.(string)
	}
	if v, ok := d.GetOk("nic_label"); ok {
		o.nicLabel = v.(string)
	}
	if v, ok := d.GetOk("portgroup"); ok {
		o.portgroup = v.(string)
	}
	return &o, nil
}
