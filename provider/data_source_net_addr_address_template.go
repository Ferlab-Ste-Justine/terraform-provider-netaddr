package provider

import (
	"github.com/Ferlab-Ste-Justine/terraform-provider-netaddr/address"

	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNetAddrAddressRead(d *schema.ResourceData, meta interface{}, rangeType string, prettify address.PrettifyAddr) error {
	conn := meta.(address.EtcdConnection)
	name := d.Get("name").(string)
	keyPrefix := d.Get("range_id").(string)

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

	addr, addrErr := conn.GetAddress(keyPrefix, name)
	if addrErr != nil {
		return addrErr
	}

	d.SetId(name)
	d.Set("address", prettify(addr))
	
	return nil
}