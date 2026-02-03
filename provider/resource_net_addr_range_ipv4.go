package provider

import (
	"github.com/Ferlab-Ste-Justine/terraform-provider-netaddr/address"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceNetAddrRangeIpv4() *schema.Resource {
	return &schema.Resource{
		Description: "Address range to create ipv4 addresses on.",
		Create: resourceNetAddrRangeIpv4Create,
		Read:   resourceNetAddrRangeIpv4Read,
		Delete: resourceNetAddrRangeIpv4Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"key_prefix": {
				Description: "Etcd key prefix for all the keys related to the range.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"first_address": {
				Description: "First assignable address in the range.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"last_address": {
				Description: "Last assignable address in the range.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceNetAddrRangeIpv4Create(d *schema.ResourceData, meta interface{}) error {
	return resourceNetAddrRangeCreate(d, meta, "ipv4", address.Ipv4StringToBytes, address.Ipv4BytesToString)
}

func resourceNetAddrRangeIpv4Read(d *schema.ResourceData, meta interface{}) error {
	return resourceNetAddrRangeRead(d, meta, "ipv4", address.Ipv4BytesToString)
}

func resourceNetAddrRangeIpv4Delete(d *schema.ResourceData, meta interface{}) error {
	return resourceNetAddrRangeDelete(d, meta)
}