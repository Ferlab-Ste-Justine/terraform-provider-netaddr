package provider

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceNetAddrAddress() *schema.Resource {
	return &schema.Resource{
		Description: "Ipv4 or mac address.",
		Create: resourceNetAddrAddressCreate,
		Read:   resourceNetAddrAddressRead,
		Delete: resourceNetAddrAddressDelete,
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

func resourceNetAddrAddressCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(EtcdConnection)
	name, _ := d.GetOk("name")
	keyPrefix, _ := d.GetOk("range_id")
	hAddr, hAddrExists := d.GetOk("hardcoded_address")

	addrRange, addrRangeExists, addrRangeErr := conn.GetAddrRange(keyPrefix.(string))
	if !addrRangeExists {
		return errors.New(fmt.Sprintf("Error retrieving address range at prefix '%s': Range does not exist", keyPrefix))
	}
	if addrRangeErr != nil {
		return errors.New(fmt.Sprintf("Error retrieving address range at prefix '%s': %s", keyPrefix, addrRangeErr.Error()))
	}

	var parse ParseAddr
	var prettify PrettifyAddr
	var addrIsGreater AddressIsGreater
	var incAddr IncrementAddress
	if addrRange.Type == "ipv4" {
		parse = Ipv4StringToBytes
		prettify = Ipv4BytesToString
		addrIsGreater = AddressGreaterThan
		incAddr = IncAddressBy1
	} else {
		parse = MacStringToBytes
		prettify = MacBytesToString
		addrIsGreater = AddressGreaterThan
		incAddr = IncAddressBy1
	}

	if hAddrExists {
		addrAsBytes, err := parse(hAddr.(string))
		if err != nil {
			return err
		}

		err = conn.CreateHardcodedAddress(keyPrefix.(string), name.(string), addrAsBytes, prettify)
		if err != nil {
			return err
		}
	} else {
		_, err := conn.CreateGeneratedAddress(keyPrefix.(string), name.(string), addrIsGreater, incAddr)
		if err != nil {
			return err
		}
	}

	d.SetId(name.(string))
	return resourceNetAddrAddressRead(d, meta)
}

func resourceNetAddrAddressRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(EtcdConnection)
	name, _ := d.GetOk("name")
	keyPrefix, _ := d.GetOk("range_id")

	addrRange, addrRangeExists, addrRangeErr := conn.GetAddrRange(keyPrefix.(string))
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

	addr, err := conn.GetAddress(keyPrefix.(string), name.(string))
	if err != nil {
		return err
	}
	d.Set("address", prettify(addr))

	return nil
}

func resourceNetAddrAddressDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(EtcdConnection)
	name, _ := d.GetOk("name")
	keyPrefix, _ := d.GetOk("range_id")
	_, hAddrExists := d.GetOk("hardcoded_address")
	addr, addrExists := d.GetOk("address")

	if !addrExists {
		return errors.New(fmt.Sprintf("Cannot delete address named '%s': Address is missing", name.(string)))
	}

	addrRange, addrRangeExists, addrRangeErr := conn.GetAddrRange(keyPrefix.(string))
	if !addrRangeExists {
		return errors.New(fmt.Sprintf("Error retrieving address range at prefix '%s': Range does not exist", keyPrefix))
	}
	if addrRangeErr != nil {
		return errors.New(fmt.Sprintf("Error retrieving address range at prefix '%s': %s", keyPrefix, addrRangeErr.Error()))
	}

	var parse ParseAddr
	var prettify PrettifyAddr
	var addrIsLess AddressIsLess
	if addrRange.Type == "ipv4" {
		parse = Ipv4StringToBytes
		prettify = Ipv4BytesToString
		addrIsLess = AddressLessThan
	} else {
		parse = MacStringToBytes
		prettify = MacBytesToString
		addrIsLess = AddressLessThan
	}

	addrAsBytes, err := parse(addr.(string))
	if err != nil {
		return err
	}

	if hAddrExists {
		err := conn.DeleteHardcodedAddress(keyPrefix.(string), name.(string), addrAsBytes, prettify, addrIsLess)
		if err != nil {
			return err
		}
	} else {
		err := conn.DeleteGeneratedAddress(keyPrefix.(string), name.(string), addrAsBytes, prettify)
		if err != nil {
			return err
		}
	}

	return nil
}
