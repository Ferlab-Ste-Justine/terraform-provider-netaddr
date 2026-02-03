package provider

import (
	"github.com/Ferlab-Ste-Justine/terraform-provider-netaddr/address"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceNetAddrRangeIpv4() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieves data on an existing ipv4 address range.",
		Read: dataSourceNetAddrRangeIpv4Read,
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

func dataSourceNetAddrRangeIpv4Read(d *schema.ResourceData, meta interface{}) error {
	return dataSourceNetAddrRangeRead(d, meta, "ipv4", address.Ipv4BytesToString)
}