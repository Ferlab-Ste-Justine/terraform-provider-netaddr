package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceNetAddrRangeMac() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieves data on an existing mac address range.",
		Read: dataSourceNetAddrRangeMacRead,
		Schema: map[string]*schema.Schema{
			"key_prefix": &schema.Schema{
				Description: "Etcd key prefix for address range.",
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
		},
	}
}

func dataSourceNetAddrRangeMacRead(d *schema.ResourceData, meta interface{}) error {
	return dataSourceNetAddrRangeRead(d, meta, "mac", MacBytesToString)
}