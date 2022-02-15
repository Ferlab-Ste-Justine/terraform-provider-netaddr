package provider

import(
	"bytes"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNetAddrAddressCreate(d *schema.ResourceData, meta interface{}, rangeType string, parse ParseAddr, prettify PrettifyAddr, incAddr IncrementAddress, addrIsGreater AddressIsGreater) error {
	conn := meta.(EtcdConnection)
	name, _ := d.GetOk("name")
	keyPrefix, _ := d.GetOk("range_id")
	hAddr, setAsHardcoded := d.GetOk("hardcoded_address")

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

	if !conn.Strict {
		addrExists, addrIsHardcoded, addr, detailsErr := conn.GetAddressDetails(keyPrefix.(string), name.(string))
		if detailsErr != nil {
			return errors.New(fmt.Sprintf("Error retrieving address details in non-strict mode at prefix '%s': %s", keyPrefix, detailsErr.Error()))
		}

		if addrExists {
			if (setAsHardcoded && !addrIsHardcoded) || (addrIsHardcoded && !setAsHardcoded) {
				return errors.New(fmt.Sprintf("Error creating address in non-strict mode at prefix '%s': Pre-existing address with the same name has an hardcoded/generated state that doesn't match the terraform configuration", keyPrefix))
			}

			if setAsHardcoded {
				setAddrAsBytes, err := parse(hAddr.(string))
				if err != nil {
					return err
				}

				if !bytes.Equal(setAddrAsBytes, addr) {
					return errors.New(fmt.Sprintf("Error creating hardcoded address in non-strict mode at prefix '%s': Pre-existing address doesn't match set address value", keyPrefix))
				}
			}

			d.SetId(name.(string))
			return resourceNetAddrAddressRead(d, meta, rangeType, prettify)
		}
	}

	if setAsHardcoded {
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
		if !conn.Strict {
			d.SetId("")
			return nil
		}
		
		return errors.New(fmt.Sprintf("Error retrieving address range at prefix '%s': Range does not exist", keyPrefix))
	}
	if addrRangeErr != nil {
		return errors.New(fmt.Sprintf("Error retrieving address range at prefix '%s': %s", keyPrefix, addrRangeErr.Error()))
	}
	if addrRange.Type != rangeType {
		return errors.New(fmt.Sprintf("Error retrieving address range at prefix '%s': Range type does not match address type", keyPrefix))
	}	

	if !conn.Strict {
		addrExists, _, _, detailsErr := conn.GetAddressDetails(keyPrefix.(string), name.(string))
		if detailsErr != nil {
			return errors.New(fmt.Sprintf("Error retrieving address details in non-strict mode at prefix '%s': %s", keyPrefix, detailsErr.Error()))
		}

		if !addrExists {
			d.SetId("")
			return nil
		}
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
	_, setAsHardcoded := d.GetOk("hardcoded_address")
	addr, addrExists := d.GetOk("address")

	if !conn.Strict {
		addrExists, _, _, detailsErr := conn.GetAddressDetails(keyPrefix.(string), name.(string))
		if detailsErr != nil {
			return errors.New(fmt.Sprintf("Error retrieving address details in non-strict mode at prefix '%s': %s", keyPrefix, detailsErr.Error()))
		}

		if !addrExists {
			return nil
		}
	}

	if !addrExists {
		return errors.New(fmt.Sprintf("Cannot delete address named '%s': Address is missing", name.(string)))
	}

	addrAsBytes, err := parse(addr.(string))
	if err != nil {
		return err
	}

	if setAsHardcoded {
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
