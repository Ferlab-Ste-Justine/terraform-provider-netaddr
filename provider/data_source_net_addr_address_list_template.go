package provider

import (
	"errors"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNetAddrAddressListRead(d *schema.ResourceData, meta interface{}, rangeType string, prettify PrettifyAddr) error {
	conn := meta.(EtcdConnection)
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

	addrList, addrListErr := conn.GetAddressList(keyPrefix)
	if addrListErr != nil {
		return addrListErr
	}

	sort.SliceStable(addrList, func(i, j int) bool {
		return addrList[i].Name < addrList[j].Name
	})

	schemaList := make([]map[string]interface{}, 0)
	for _, addr := range addrList {
		schemaList = append(schemaList, map[string]interface{}{
			"name": addr.Name,
			"address": prettify(addr.Address),
		})
	}

	d.SetId(keyPrefix)
	d.Set("addresses", schemaList)
	
	return nil
}