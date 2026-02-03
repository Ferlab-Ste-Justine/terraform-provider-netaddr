package provider

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNetAddrAddressV2Read(d *schema.ResourceData, meta interface{}, rangeType string, prettify PrettifyAddr) error {
	conn := meta.(EtcdConnection)
	name := d.Get("name").(string)

	keyPrefixes := GetRangeIdsFromResource(d)

	for _, keyPrefix := range keyPrefixes {
		addrRange, addrRangeExists, addrRangeErr := conn.GetAddrRange(keyPrefix)
		if !addrRangeExists {
			return errors.New(fmt.Sprintf("Error retrieving address range at prefix '%s': Range does not exist", keyPrefix))
		}
		if addrRangeErr != nil {
			return errors.New(fmt.Sprintf("Error retrieving address range at prefix '%s': %s", keyPrefix, addrRangeErr.Error()))
		}
		if addrRange.Type != rangeType {
			return errors.New(fmt.Sprintf("Error retrieving address range at prefix '%s': Range type doesn't match", keyPrefix))
		}

		addr, found, addrErr := conn.FindAddress(keyPrefix, name)
		if addrErr != nil {
			return addrErr
		}
		if !found {
			continue
		}

		d.SetId(name)
		d.Set("address", prettify(addr))
		d.Set("found_in_range", keyPrefix)
		
		return nil
	}
	
	return errors.New(fmt.Sprintf("Error retrieving address named '%s': Address was not found in any of the input ranges", name))
}