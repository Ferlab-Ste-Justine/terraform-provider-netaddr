package provider

import(
	"bytes"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNetAddrRangeCreate(d *schema.ResourceData, meta interface{}, rangeType string, parse ParseAddr, prettify PrettifyAddr) error {
	conn := meta.(EtcdConnection)
	keyPrefix, _ := d.GetOk("key_prefix")

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
		Type: rangeType,
		FirstAddress: firstAddrBytes,
		LastAddress: lastAddrBytes,
	}

	if !conn.Strict {
		addrRange, addrRangeExists, addrRangeErr := conn.GetAddrRange(keyPrefix.(string))
		if addrRangeErr != nil {
			return errors.New(fmt.Sprintf("Error retrieving address range details in non-strict mode: %s", addrRangeErr.Error()))
		}

		if addrRangeExists {
			if (!bytes.Equal(firstAddrBytes, addrRange.FirstAddress)) || (!bytes.Equal(lastAddrBytes, addrRange.LastAddress)) {
				return errors.New(fmt.Sprintf("Error creating address range in non-strict mode: Pre-existing address range doesn't match specified address range"))
			}
			d.SetId(keyPrefix.(string))
			return resourceNetAddrRangeRead(d, meta, rangeType, prettify)
		}
	}

	creationErr := conn.CreateAddrRange(keyPrefix.(string), addrRange)
	if creationErr != nil {
		return errors.New(fmt.Sprintf("Error creating address range: %s", creationErr.Error()))
	}

	d.SetId(keyPrefix.(string))
	return resourceNetAddrRangeRead(d, meta, rangeType, prettify)
}

func resourceNetAddrRangeRead(d *schema.ResourceData, meta interface{}, rangeType string, prettify PrettifyAddr) error {
	keyPrefix := d.Id()
	conn := meta.(EtcdConnection)

	addrRange, addrRangeExists, addrRangeErr := conn.GetAddrRange(keyPrefix)
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
		return errors.New(fmt.Sprintf("Error retrieving address range at prefix '%s': Range type doesn't match", keyPrefix))
	}

	d.Set("key_prefix", keyPrefix)
	d.Set("first_address", prettify(addrRange.FirstAddress))
	d.Set("last_address", prettify(addrRange.LastAddress))

	return nil
}

func resourceNetAddrRangeDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(EtcdConnection)
	keyPrefix, _ := d.GetOk("key_prefix")

	if !conn.Strict {
		_, addrRangeExists, addrRangeErr := conn.GetAddrRange(keyPrefix.(string))
		if addrRangeErr != nil {
			return errors.New(fmt.Sprintf("Error retrieving address range details in non-strict mode: %s", addrRangeErr.Error()))
		}

		if !addrRangeExists {
			return nil
		}
	}

	err := conn.DestroyAddrRange(keyPrefix.(string))
	if err != nil {
		return errors.New(fmt.Sprintf("Error destroying address range: %s", err.Error()))
	}

	return nil
}