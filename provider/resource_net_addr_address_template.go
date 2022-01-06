package provider

import(
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNetAddrAddressCreate(d *schema.ResourceData, meta interface{}, rangeType string, parse ParseAddr, prettify PrettifyAddr, incAddr IncrementAddress, addrIsGreater AddressIsGreater) error {
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
	if addrRange.Type != rangeType {
		return errors.New(fmt.Sprintf("Error retrieving address range at prefix '%s': Range type does not match address type", keyPrefix))
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
	return resourceNetAddrAddressRead(d, meta, rangeType, prettify)
}

func resourceNetAddrAddressRead(d *schema.ResourceData, meta interface{}, rangeType string, prettify PrettifyAddr) error {
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
	if addrRange.Type != rangeType {
		return errors.New(fmt.Sprintf("Error retrieving address range at prefix '%s': Range type does not match address type", keyPrefix))
	}	

	addr, err := conn.GetAddress(keyPrefix.(string), name.(string))
	if err != nil {
		return err
	}
	d.Set("address", prettify(addr))

	return nil
}

func resourceNetAddrAddressDelete(d *schema.ResourceData, meta interface{}, parse ParseAddr, prettify PrettifyAddr, addrIsLess AddressIsLess) error {
	conn := meta.(EtcdConnection)
	name, _ := d.GetOk("name")
	keyPrefix, _ := d.GetOk("range_id")
	_, hAddrExists := d.GetOk("hardcoded_address")
	addr, addrExists := d.GetOk("address")

	if !addrExists {
		return errors.New(fmt.Sprintf("Cannot delete address named '%s': Address is missing", name.(string)))
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
