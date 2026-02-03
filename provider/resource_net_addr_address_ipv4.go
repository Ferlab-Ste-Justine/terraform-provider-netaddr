package provider

import (
	"github.com/Ferlab-Ste-Justine/terraform-provider-netaddr/address"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceNetAddrAddressIpv4() *schema.Resource {
	return &schema.Resource{
		Description: "Ipv4 address.",
		Create: resourceNetAddrAddressIpv4Create,
		Read:   resourceNetAddrAddressIpv4Read,
		Delete: resourceNetAddrAddressIpv4Delete,
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

func resourceNetAddrAddressIpv4Create(d *schema.ResourceData, meta interface{}) error {
	return resourceNetAddrAddressCreate(d, meta, "ipv4", address.Ipv4StringToBytes, address.Ipv4BytesToString, address.IncAddressBy1, address.AddressGreaterThan)
}

func resourceNetAddrAddressIpv4Read(d *schema.ResourceData, meta interface{}) error {
	return resourceNetAddrAddressRead(d, meta, "ipv4", address.Ipv4BytesToString)
}

func resourceNetAddrAddressIpv4Delete(d *schema.ResourceData, meta interface{}) error {
	return resourceNetAddrAddressDelete(d, meta, address.Ipv4StringToBytes, address.Ipv4BytesToString, address.AddressLessThan)
}