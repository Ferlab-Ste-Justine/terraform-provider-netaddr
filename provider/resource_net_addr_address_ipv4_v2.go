package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceNetAddrAddressIpv4V2() *schema.Resource {
	return &schema.Resource{
		Description: "Ipv4 address. Version 2 adds support for assignment from multiple ranges (useful if you get an extra range of ips from the same subnet later on).",
		Create: resourceNetAddrAddressIpv4V2Create,
		Read:   resourceNetAddrAddressIpv4V2Read,
		Delete: resourceNetAddrAddressIpv4V2Delete,
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

func resourceNetAddrAddressIpv4V2Create(d *schema.ResourceData, meta interface{}) error {
	return resourceNetAddrAddressV2Create(d, meta, "ipv4", Ipv4StringToBytes, Ipv4BytesToString, IncAddressBy1, AddressGreaterThan)
}

func resourceNetAddrAddressIpv4V2Read(d *schema.ResourceData, meta interface{}) error {
	return resourceNetAddrAddressV2Read(d, meta, "ipv4", Ipv4BytesToString)
}

func resourceNetAddrAddressIpv4V2Delete(d *schema.ResourceData, meta interface{}) error {
	return resourceNetAddrAddressV2Delete(d, meta, Ipv4StringToBytes, Ipv4BytesToString, AddressLessThan)
}