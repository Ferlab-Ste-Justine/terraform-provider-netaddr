package provider

import (
	"github.com/Ferlab-Ste-Justine/terraform-provider-netaddr/address"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceNetAddrRangeMac() *schema.Resource {
	return &schema.Resource{
		Description: "Address range to create mac addresses on.",
		Create: resourceNetAddrRangeMacCreate,
		Read:   resourceNetAddrRangeMacRead,
		Delete: resourceNetAddrRangeMacDelete,
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

func resourceNetAddrRangeMacCreate(d *schema.ResourceData, meta interface{}) error {
	return resourceNetAddrRangeCreate(d, meta, "mac", address.MacStringToBytes, address.MacBytesToString)
}

func resourceNetAddrRangeMacRead(d *schema.ResourceData, meta interface{}) error {
	return resourceNetAddrRangeRead(d, meta, "mac", address.MacBytesToString)
}

func resourceNetAddrRangeMacDelete(d *schema.ResourceData, meta interface{}) error {
	return resourceNetAddrRangeDelete(d, meta)
}