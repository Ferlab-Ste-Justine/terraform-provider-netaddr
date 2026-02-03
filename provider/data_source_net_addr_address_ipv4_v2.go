package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceNetAddrAddressIpv4V2() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieves data on an existing ipv4 address. Version 2 adds support for an ip address assigned from multiple ranges (useful if you get an extra range of ips from the same subnet later on).",
		Read: dataSourceNetAddrAddressIpv4V2Read,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Name of the address.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"range_ids": {
				Description: "Identifiers of the address ranges the address is tied to.",
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},		
			},
			"found_in_range": {
				Description: "Id of the range the address is in.",
				Type:         schema.TypeString,
				Computed:     true,
			},
			"address": {
				Description: "The address that got assigned to the resource.",
				Type:         schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNetAddrAddressIpv4V2Read(d *schema.ResourceData, meta interface{}) error {
	return dataSourceNetAddrAddressV2Read(d, meta, "ipv4", Ipv4BytesToString)
}