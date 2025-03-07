package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceNetAddrRangeUsageIpv4() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieves ipv4 addresses utilisation data on an address range.",
		Read: dataSourceNetAddrRangeUsageIpv4Read,
		Schema: map[string]*schema.Schema{
			"range_id": &schema.Schema{
				Description: "Identifier of the address range to get the capacity from.",
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"capacity": {
				Description: "Number of addresses in the range.",
				Type:         schema.TypeInt,
				Computed: true,
			},
			"used_capacity": {
				Description: "Number of used addresses in the range.",
				Type:         schema.TypeInt,
				Computed: true,
			},
			"free_capacity": {
				Description: "Number of free addresses in the range.",
				Type:         schema.TypeInt,
				Computed: true,
			},
		},
	}
}


func dataSourceNetAddrRangeUsageIpv4Read(d *schema.ResourceData, meta interface{}) error {
	return dataSourceNetAddrRangeUsageRead(d, meta, "ipv4", Ipv4RangeAddressCount)
}