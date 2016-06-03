package vsphere

import "github.com/hashicorp/terraform/helper/schema"
import "fmt"

type dvs 					struct {
	name					string
	folder					string
	datacenter				string
	extensionKey				string
	description				string
	contact					struct {
		name				string
		infos				string
	}
	switchUsagePolicy			struct {
		autoPreinstallAllowed		bool
		autoUpgradeAllowed		bool
		partialUpgradeAllowed		bool
	}
	switchIPAddress				string
	numStandalonePorts			int
}

type dvs_map_host_dvs 				struct {
	hostName				string
	switchName				string
}

type dvs_port_group 				struct {
	name					string
	pgType					string
	description				string
	autoExpand				bool
	numPorts				int
	portNameFormat				string
	policy					struct {
		allowBlockOverride		bool
		allowLivePortMoving		bool
		allowNetworkRPOverride		bool
		portConfigResetDisconnect	bool
		allowShapingOverride		bool
		allowTrafficFilterOverride	bool
		allowVendorConfigOverride	bool
	}
}

type dvs_map_vm_dvs				struct {
	vm					string
	nicIndex				int
	portgroup				string
	port					string
}

// parse a provided Terraform config into a dvs struct
func parseDVS(d *schema.ResourceData) (*dvs, error) {
	f := dvs{};
	if v, ok := d.GetOk("name"); ok {
		f.name = v.(string)
	}
	if v, ok := d.GetOk("folder"); ok {
		f.folder = v.(string)
	}
	if v, ok := d.GetOk("datacenter"); ok {
		f.datacenter = v.(string)
	}
	if v, ok := d.GetOk("extension_key"); ok {
		f.extensionKey = v.(string)
	}
	if v, ok := d.GetOk("description"); ok {
		f.description = v.(string)
	}
	if v, ok := d.GetOk("switch_ip_address"); ok {
		f.switchIPAddress = v.(string)
	}
	if v, ok := d.GetOk("num_standalone_ports"); ok {
		f.numStandalonePorts = v.(int)
	}
	// contact
	if s, ok := d.GetOk("contact"); ok {
		vmap, casted := s.(map[string]interface{})
		if !casted {
			return nil, fmt.Errorf("Cannot cast contact as a string map. Contact: %+v", s)
		}
		f.contact.name = vmap["name"].(string)
		f.contact.infos = vmap["infos"].(string)
	}
	if s, ok := d.GetOk("switch_usage_policy"); ok {
		vmap, casted := s.(map[string]interface{})
		if !casted {
			return nil, fmt.Errorf("Cannot cast switch_usage_policy as a string map. Contact: %+v", s)
		}
		f.switchUsagePolicy.autoPreinstallAllowed = vmap["auto_preinstall_allowed"].(bool)
		f.switchUsagePolicy.autoUpgradeAllowed = vmap["auto_upgrade_allowed"].(bool)
		f.switchUsagePolicy.partialUpgradeAllowed = vmap["partial_upgrade_allowed"].(bool)
	}

	return &f, nil
}

func resourceVSphereDVSSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema {
		"name": &schema.Schema {
			Type:				schema.TypeString,
			Required:		true,
			// ForceNew:		true,
		},
		"extension_key": &schema.Schema {
			Type:				schema.TypeString,
			Optional:		true,
		},
		"description": &schema.Schema {
			Type:				schema.TypeString,
			Optional:		true,
		},
		"contact": &schema.Schema {
			Type:				schema.TypeList,
			Optional:		true,
			Elem: &schema.Resource {
				Schema: map[string]*schema.Schema{
					"name": &schema.Schema {
						Type:			schema.TypeString,
						Required: true,
					},
					"infos": &schema.Schema {
						Type:			schema.TypeString,
						Required: true,
					},
				},
			},
		},
		"switch_usage_policy": &schema.Schema {
			Type:				schema.TypeList,
			Optional:		true,
			Elem: &schema.Resource {
				Schema: map[string]*schema.Schema{
					"auto_preinstall_allowed": &schema.Schema {
						Type:			schema.TypeBool,
						Optional: true,
					},
					"auto_upgrade_allowed": &schema.Schema {
						Type:			schema.TypeBool,
						Optional: true,
					},
					"partial_upgrade_allowed": &schema.Schema {
						Type:			schema.TypeBool,
						Optional: true,
					},
				},
			},
		},
		"switch_ip_address": &schema.Schema {
			Type:				schema.TypeString,
			Optional:		true,
		},
		"num_standalone_ports": &schema.Schema {
			Type:				schema.TypeInt,
			Optional:		true,
		},
	}
}
/* functions for DistributedVirtualPortgroup */
func resourceVSphereDVSPGSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema {
		"name": &schema.Schema {
			Type:				schema.TypeString,
			Required:		true,
			ForceNew:		true,
		},
		"datacenter": &schema.Schema {
			Type:				schema.TypeString,
			Required:		true,
			ForceNew:		true,
		},
		"type": &schema.Schema {
			Type:				schema.TypeString,
			Required:		true,
		},
		"description": &schema.Schema {
			Type:				schema.TypeString,
			Optional:		true,
		},
		"auto_expand": &schema.Schema {
			Type:				schema.TypeBool,
			Optional:		true,
		},
		"num_ports": &schema.Schema {
			Type:				schema.TypeInt,
			Optional:		true,
		},
		"port_name_format": &schema.Schema {
			Type:				schema.TypeString,
			Optional:		true,
		},
		"policy": &schema.Schema {
			Type:				schema.TypeList,
			Optional:		true,
			Elem: &schema.Resource {
				Schema: map[string]*schema.Schema{
					"allow_block_override": &schema.Schema {
						Type:			schema.TypeBool,
						Optional: true,
					},
					"allow_live_port_moving": &schema.Schema {
						Type:			schema.TypeBool,
						Optional: true,
					},
					"allov_network_resources_pool_override": &schema.Schema {
						Type:			schema.TypeBool,
						Optional: true,
					},
					"port_config_reset_disconnect": &schema.Schema {
						Type:			schema.TypeBool,
						Optional: true,
					},
					"allow_shaping_override": &schema.Schema {
						Type:			schema.TypeBool,
						Optional: true,
					},
					"allow_traffic_filter_override": &schema.Schema {
						Type:			schema.TypeBool,
						Optional: true,
					},
					"allow_vendor_config_override": &schema.Schema {
						Type:			schema.TypeBool,
						Optional: true,
					},
				},
			},
		},
	}
}

/* functions for MapHostDVS */

// parse a DVSPG ResourceData to a dvs_port_group struct
func parseDVSPG(d *schema.ResourceData) (*dvs_port_group, error) {
	o := dvs_port_group{}
	if v, ok := d.GetOk("name"); ok {
		o.name = v.(string)
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
			return nil, fmt.Errorf("Cannot cast policy as a string map. See: %+v", s)
		}
		o.policy.allowBlockOverride = vmap["allow_block_override"].(bool)
		o.policy.allowLivePortMoving = vmap["allow_live_port_moving"].(bool)
		o.policy.allowNetworkRPOverride = vmap["allow_network_rp_override"].(bool)
		o.policy.portConfigResetDisconnect = vmap["port_config_reset_disconnect"].(bool)
		o.policy.allowShapingOverride = vmap["allow_shaping_override"].(bool)
		o.policy.allowTrafficFilterOverride = vmap["allow_traffic_filter_override"].(bool)
		o.policy.allowVendorConfigOverride = vmap["allow_vendor_config_override"].(bool)
	}
	return &o, nil
}
/* MapHostDVS functions */
func resourceVSphereMapHostDVSSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema {
		"host": &schema.Schema {
			Type:				schema.TypeString,
			Required:		true,
			// ForceNew:		true,
		},
		"switch": &schema.Schema {
			Type:				schema.TypeString,
			Required:		true,
			// ForceNew:		true,
		},
	}
}

/* parse a MapHostDVS to its struct */
func parseMapHostDVS(d *schema.ResourceData) (*dvs_map_host_dvs, error) {
	o := dvs_map_host_dvs{};
	if v, ok := d.GetOk("host_name"); ok {
		o.hostName = v.(string)
	}
	if v, ok := d.GetOk("switch_name"); ok {
		o.switchName = v.(string)
	}
	return &o, nil
}
/* Functions for MapVMDVS */

func resourceVSphereMapVMDVPGSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema {
		"vm": &schema.Schema {
			Type:				schema.TypeString,
			Required:		true,
			ForceNew:		true,
		},
		"nic_index": &schema.Schema {
			Type:				schema.TypeInt,
			Required:		true,
			ForceNew:		true,
		},
		"portgroup": &schema.Schema {
			Type:				schema.TypeString,
			Required:		true,
			ForceNew:		true,
		},
		"port": &schema.Schema {
			Type:				schema.TypeString,
			Required:		true,
			ForceNew:		true,
		},
	}
}

/* parse a MapVMDVS to its struct */
func parseMapVMDVS(d *schema.ResourceData) (*dvs_map_vm_dvs, error) {
	o := dvs_map_vm_dvs{};
	if v, ok := d.GetOk("vm"); ok {
		o.vm = v.(string)
	}
	if v, ok := d.GetOk("nic_index"); ok {
		o.nicIndex = v.(int)
	}
	if v, ok := d.GetOk("portgroup"); ok {
		o.portgroup = v.(string)
	}
	if v, ok := d.GetOk("port"); ok {
		o.port = v.(string)
	}
	return &o, nil
}
