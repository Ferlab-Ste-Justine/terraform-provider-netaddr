package provider

import (
	"github.com/Ferlab-Ste-Justine/terraform-provider-netaddr/address"

	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNetAddrRangeUsageRead(d *schema.ResourceData, meta interface{}, rangeType string, rangeAddrCount address.RangeAddressCount) error {
	conn := meta.(address.EtcdConnection)
	keyPrefix := d.Get("range_id").(string)

	addrRange, addrRangeExists, addrRangeErr := conn.GetAddrRange(keyPrefix)
	if !addrRangeExists {
		return errors.New(fmt.Sprintf("Error retrieving address range at prefix '%s': Range does not exist", keyPrefix))
	}
	if addrRangeErr != nil {
		return addrRangeErr
	}
	if addrRange.Type != rangeType {
		return errors.New(fmt.Sprintf("Error retrieving address range at prefix '%s': Range type doesn't match", keyPrefix))
	}

	usage, usageErr := conn.GetAddrRangeUsage(keyPrefix, address.Ipv4RangeAddressCount)
	if usageErr != nil {
		return usageErr
	}

	d.SetId(keyPrefix)
	d.Set("capacity", usage.Capacity)
	d.Set("used_capacity", usage.UsedCapacity)
	d.Set("free_capacity", usage.FreeCapacity)

	return nil
}