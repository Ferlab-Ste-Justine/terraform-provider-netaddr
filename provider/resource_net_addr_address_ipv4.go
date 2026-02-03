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
		Update: resourceNetAddrAddressIpv4Update,
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
			"retain_on_delete": &schema.Schema{
				Description: "Whether to retain the address in etcd when the resource is deleted. Useful to set to true if you wish to migrate the address to another terraform project or migrate to the v2 version of the resource.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    false,
			},
			"manage_existing": &schema.Schema{
				Description: "Whether the address is possibly present when the resource is created. Setting this to true allows you to import the existing address without error.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    false,
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

func resourceNetAddrAddressIpv4Update(d *schema.ResourceData, meta interface{}) error {
	return resourceNetAddrAddressRead(d, meta, "ipv4", address.Ipv4BytesToString)
}

func resourceNetAddrAddressIpv4Delete(d *schema.ResourceData, meta interface{}) error {
	return resourceNetAddrAddressDelete(d, meta, address.Ipv4StringToBytes, address.Ipv4BytesToString, address.AddressLessThan)
}