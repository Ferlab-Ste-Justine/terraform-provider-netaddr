package provider

import (
	"github.com/Ferlab-Ste-Justine/terraform-provider-netaddr/address"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceNetAddrRangeKeyspaceMac() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieves the lower level keyspace details of an mac addresses space. See github repo README for details about the keyspace",
		Read: dataSourceNetAddrRangeKeyspaceMacRead,
		Schema: map[string]*schema.Schema{
			"range_id": &schema.Schema{
				Description: "Identifier of the address range to get the key space from.",
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"first_address": {
				Description: "First assignable address in the range.",
				Type:         schema.TypeString,
				Computed: true,
			},
			"last_address": {
				Description: "Last assignable address in the range.",
				Type:         schema.TypeString,
				Computed: true,
			},
			"next_address": {
				Description: "Next assignable new address in the range.",
				Type:         schema.TypeString,
				Computed: true,
			},
			"addresses": {
				Description: "List of all addresses in the range.",
				Type:         schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description:  "Name assigned to the adress",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"address": {
							Description:  "The address",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
			"generated_addresses": {
				Description: "List of all addresses that are flagged as generated in the range.",
				Type:         schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description:  "Name assigned to the adress",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"address": {
							Description:  "The address",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
			"hardcoded_addresses": {
				Description: "List of all addresses that are flagged as hardcoded in the range.",
				Type:         schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description:  "Name assigned to the adress",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"address": {
							Description:  "The address",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
			"deleted_addresses": {
				Description: "List of all addresses that were deleted and are available to be reclaimed in the range.",
				Type:         schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description:  "Name assigned to the adress",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"address": {
							Description:  "The address",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
		},
	}
}


func dataSourceNetAddrRangeKeyspaceMacRead(d *schema.ResourceData, meta interface{}) error {
	return dataSourceNetAddrRangeKeyspaceRead(d, meta, "mac", address.MacBytesToString)
}