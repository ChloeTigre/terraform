package dvs

import "github.com/hashicorp/terraform/helper/schema"

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
		"default_vlan": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		"datacenter": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"type": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "earlyBinding|ephemeral",
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

/* Functions for MapVMDVPG */

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
