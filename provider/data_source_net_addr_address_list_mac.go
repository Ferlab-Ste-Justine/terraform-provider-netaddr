package provider

import (
	"github.com/Ferlab-Ste-Justine/terraform-provider-netaddr/address"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceNetAddrAddressListMac() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieves all mac addresses in a range.",
		Read: dataSourceNetAddrAddressListMacRead,
		Schema: map[string]*schema.Schema{
			"range_id": &schema.Schema{
				Description: "Identifier of the address range to get the addresses from.",
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"addresses": {
				Description: "List of addresses in the range.",
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

func dataSourceNetAddrAddressListMacRead(d *schema.ResourceData, meta interface{}) error {
	return dataSourceNetAddrAddressListRead(d, meta, "mac", address.MacBytesToString)
}