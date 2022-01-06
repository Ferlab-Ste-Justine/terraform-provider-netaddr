package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceNetAddrAddressMac() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieves data on an existing mac address.",
		Read: dataSourceNetAddrAddressMacRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Name of the address.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"range_id": {
				Description: "Identifier of the address range the address is tied to.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"address": {
				Description: "The address that got assigned to the resource.",
				Type:         schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNetAddrAddressMacRead(d *schema.ResourceData, meta interface{}) error {
	return dataSourceNetAddrAddressRead(d, meta, "mac", MacBytesToString)
}