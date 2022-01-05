package provider

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceNetAddrRange() *schema.Resource {
	return &schema.Resource{
		Description: "Address range to create ipv4 and mac addresses on.",
		Create: resourceNetAddrRangeCreate,
		Read:   resourceNetAddrRangeRead,
		Delete: resourceNetAddrRangeDelete,
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
			"type": {
				Description: "Type of address the range will contain. Only 'ipv4' or 'mac' are currently supported.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
                    v := val.(string)

					if v != "ipv4" && v != "mac" {
						return []string{}, []error{errors.New("For the type field of netaddr_range, only the following values are supported: ipv4, mac")}
					}

                    return []string{}, []error{}
                },
			},
		},
	}
}

func resourceNetAddrRangeCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(EtcdConnection)
	keyPrefix, _ := d.GetOk("key_prefix")
	rangeType, _ := d.GetOk("type")

	var parse ParseAddr
	if rangeType.(string) == "ipv4" {
		parse = Ipv4StringToBytes
	} else {
		parse = MacStringToBytes
	}

	firstAddr, _ := d.GetOk("first_address")
	firstAddrBytes, firstAddrErr := parse(firstAddr.(string))
	if firstAddrErr != nil {
		return errors.New(fmt.Sprintf("Error creating address range: %s", firstAddrErr.Error()))
	}

	lastAddr, _ := d.GetOk("last_address")
	lastAddrBytes, lastAddrErr := parse(lastAddr.(string))
	if lastAddrErr != nil {
		return errors.New(fmt.Sprintf("Error creating address range: %s", lastAddrErr.Error()))
	}

	addrRange := AddressRange{
		Type: rangeType.(string),
		FirstAddress: firstAddrBytes,
		LastAddress: lastAddrBytes,
	}

	creationErr := conn.CreateAddrRange(keyPrefix.(string), addrRange)
	if creationErr != nil {
		return errors.New(fmt.Sprintf("Error creating address range: %s", creationErr.Error()))
	}

	d.SetId(keyPrefix.(string))
	return resourceNetAddrRangeRead(d, meta)
}

func resourceNetAddrRangeRead(d *schema.ResourceData, meta interface{}) error {
	keyPrefix := d.Id()
	conn := meta.(EtcdConnection)

	addrRange, addrRangeExists, addrRangeErr := conn.GetAddrRange(keyPrefix)
	if !addrRangeExists {
		return errors.New(fmt.Sprintf("Error retrieving address range at prefix '%s': Range does not exist", keyPrefix))
	}
	if addrRangeErr != nil {
		return errors.New(fmt.Sprintf("Error retrieving address range at prefix '%s': %s", keyPrefix, addrRangeErr.Error()))
	}

	var prettify PrettifyAddr
	if addrRange.Type == "ipv4" {
		prettify = Ipv4BytesToString
	} else {
		prettify = MacBytesToString
	}
	

	d.Set("key_prefix", keyPrefix)
	d.Set("first_address", prettify(addrRange.FirstAddress))
	d.Set("last_address", prettify(addrRange.LastAddress))

	return nil
}

func resourceNetAddrRangeDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(EtcdConnection)
	keyPrefix, _ := d.GetOk("key_prefix")

	err := conn.DestroyAddrRange(keyPrefix.(string))
	if err != nil {
		return errors.New(fmt.Sprintf("Error destroying address range: %s", err.Error()))
	}

	return nil
}