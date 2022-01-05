package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceNetAddrMac() *schema.Resource {
	return &schema.Resource{
		Description: "Mac address.",
		Create: resourceNetAddrMacCreate,
		Read:   resourceNetAddrMacRead,
		Delete: resourceNetAddrMacDelete,
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
			"network_id": {
				Description: "Identifier of the network the address is tied to.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"address": {
				Description: "The address in xx-xx-xx-xx-xx-xx format.",
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceNetAddrMacCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceNetAddrMacRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceNetAddrMacDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
