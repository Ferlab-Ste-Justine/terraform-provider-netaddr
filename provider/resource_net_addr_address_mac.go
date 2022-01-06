package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceNetAddrAddressMac() *schema.Resource {
	return &schema.Resource{
		Description: "Mac address.",
		Create: resourceNetAddrAddressMacCreate,
		Read:   resourceNetAddrAddressMacRead,
		Delete: resourceNetAddrAddressMacDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Name to associate with the address.",
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
			"hardcoded_address": {
				Description: "An optional input to fixate the address to a specific value.",
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"address": {
				Description: "The address that got assigned to the resource.",
				Type:         schema.TypeString,
				Computed:     true,
			},
		},
	}
}

func resourceNetAddrAddressMacCreate(d *schema.ResourceData, meta interface{}) error {
	return resourceNetAddrAddressCreate(d, meta, "mac", MacStringToBytes, MacBytesToString, IncAddressBy1, AddressGreaterThan)
}

func resourceNetAddrAddressMacRead(d *schema.ResourceData, meta interface{}) error {
	return resourceNetAddrAddressRead(d, meta, "mac", MacBytesToString)
}

func resourceNetAddrAddressMacDelete(d *schema.ResourceData, meta interface{}) error {
	return resourceNetAddrAddressDelete(d, meta, MacStringToBytes, MacBytesToString, AddressLessThan)
}