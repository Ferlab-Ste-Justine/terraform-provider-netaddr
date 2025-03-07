package provider

import (
	"errors"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNetAddrRangeKeyspaceRead(d *schema.ResourceData, meta interface{}, rangeType string, prettify PrettifyAddr) error {
	conn := meta.(EtcdConnection)
	keyPrefix := d.Get("range_id").(string)

	keyspace, keyspaceErr := conn.GetAddrRangeKeyspace(keyPrefix)
	if keyspaceErr != nil {
		return errors.New(fmt.Sprintf("Error retrieving keyspace info at prefix '%s': %s", keyPrefix, keyspaceErr.Error()))
	}
	if keyspace.Type != rangeType {
		return errors.New(fmt.Sprintf("Error retrieving keyspace info at prefix '%s': Range type doesn't match", keyPrefix))
	}

	addrList := keyspace.Names
	sort.SliceStable(addrList, func(i, j int) bool {
		return addrList[i].Name < addrList[j].Name
	})

	addrSchemaList := make([]map[string]interface{}, 0)
	for _, addr := range addrList {
		addrSchemaList = append(addrSchemaList, map[string]interface{}{
			"name": addr.Name,
			"address": prettify(addr.Address),
		})
	}

	genAddrList := keyspace.GeneratedAddresses
	sort.SliceStable(genAddrList, func(i, j int) bool {
		return genAddrList[i].Name < genAddrList[j].Name
	})

	genAddrSchemaList := make([]map[string]interface{}, 0)
	for _, addr := range genAddrList {
		genAddrSchemaList = append(genAddrSchemaList, map[string]interface{}{
			"name": addr.Name,
			"address": prettify(addr.Address),
		})
	}

	hardAddrList := keyspace.HardcodedAddresses
	sort.SliceStable(hardAddrList, func(i, j int) bool {
		return hardAddrList[i].Name < hardAddrList[j].Name
	})

	hardAddrSchemaList := make([]map[string]interface{}, 0)
	for _, addr := range hardAddrList {
		hardAddrSchemaList = append(hardAddrSchemaList, map[string]interface{}{
			"name": addr.Name,
			"address": prettify(addr.Address),
		})
	}

	delAddrList := keyspace.DeletedAddresses
	sort.SliceStable(delAddrList, func(i, j int) bool {
		return delAddrList[i].Name < delAddrList[j].Name
	})

	delAddrSchemaList := make([]map[string]interface{}, 0)
	for _, addr := range delAddrList {
		delAddrSchemaList = append(delAddrSchemaList, map[string]interface{}{
			"name": addr.Name,
			"address": prettify(addr.Address),
		})
	}

	d.SetId(keyPrefix)
	d.Set("first_address", prettify(keyspace.FirstAddress))
	d.Set("last_address", prettify(keyspace.LastAddress))
	d.Set("next_address", prettify(keyspace.NextAddress))
	d.Set("addresses", addrSchemaList)
	d.Set("generated_addresses", genAddrSchemaList)
	d.Set("hardcoded_addresses", hardAddrSchemaList)
	d.Set("deleted_addresses", delAddrSchemaList)
	
	return nil
}