package provider

import(
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNetAddrAddressV2Create(d *schema.ResourceData, meta interface{}, rangeType string, parse ParseAddr, prettify PrettifyAddr, incAddr IncrementAddress, addrIsGreater AddressIsGreater) error {
	conn := meta.(EtcdConnection)
	name, _ := d.GetOk("name")
	hAddr, setAsHardcoded := d.GetOk("hardcoded_address")

	keyPrefixes := GetRangeIdsFromResource(d)

	if setAsHardcoded {
		addrAsBytes, err := parse(hAddr.(string))
		if err != nil {
			return err
		}

		exists, prefix, genErr := conn.GenerateHardcodedAddressWithValidation(name.(string), keyPrefixes, addrAsBytes, rangeType, !conn.Strict, prettify)
		if genErr != nil {
			return genErr
		}

		if exists {
			log.Printf(fmt.Sprintf(
				"[WARN] Creating resource for pre-existing hardcoded address of type '%s', name '%s' and address '%s' in range '%s'", 
				rangeType,
				name.(string),
				hAddr.(string),
				prefix,
			))
		} else {
			log.Printf(fmt.Sprintf(
				"[DEBUG] Created hardcoded address of type '%s', name '%s' and address '%s' in range '%s'", 
				rangeType,
				name.(string),
				hAddr.(string),
				prefix,
			))
		}

		d.Set("found_in_range", prefix)
	} else {
		exists, addr, prefix, genErr := conn.GenerateGeneratedAddressWithValidation(name.(string), keyPrefixes, rangeType, !conn.Strict, addrIsGreater, incAddr)
		if genErr != nil {
			return genErr
		}

		if exists {
			log.Printf(fmt.Sprintf(
				"[WARN] Creating resource for pre-existing generated address of type '%s', name '%s' and address '%s' in range '%s'", 
				rangeType,
				name.(string),
				prettify(addr),
				prefix,
			))
		} else {
			log.Printf(fmt.Sprintf(
				"[DEBUG] Created generated address of type '%s', name '%s' and address '%s' in range '%s'", 
				rangeType,
				name.(string),
				prettify(addr),
				prefix,
			))
		}

		d.Set("found_in_range", prefix)
	}
	
	d.SetId(name.(string))
	return resourceNetAddrAddressRead(d, meta, rangeType, prettify)
}

func resourceNetAddrAddressV2Read(d *schema.ResourceData, meta interface{}, rangeType string, prettify PrettifyAddr) error {
	conn := meta.(EtcdConnection)
	name, _ := d.GetOk("name")
	keyPrefix := d.Get("found_in_range")

	addr, found, err := conn.GetAddressWithValidation(name.(string), keyPrefix.(string), rangeType, !conn.Strict)
	if err != nil {
		return err
	}

	if !found {
		log.Printf(fmt.Sprintf(
			"[WARN] Tried to read non-existent address of type '%s' and name '%s' in range '%s'", 
			rangeType,
			name.(string),
			keyPrefix.(string),
		))

		d.SetId("")
		return nil
	}

	prettyAddr := prettify(addr)
	d.Set("address", prettyAddr)

	log.Printf(fmt.Sprintf(
		"[DEBUG] Read address of type '%s', name '%s' and address '%s' in range '%s'", 
		rangeType,
		name.(string),
		prettyAddr,
		keyPrefix.(string),
	))

	return nil
}

func resourceNetAddrAddressV2Delete(d *schema.ResourceData, meta interface{}, parse ParseAddr, prettify PrettifyAddr, addrIsLess AddressIsLess) error {
	conn := meta.(EtcdConnection)
	name := d.Get("name")
	keyPrefix := d.Get("found_in_range")
	_, setAsHardcoded := d.GetOk("hardcoded_address")
	addr := d.Get("address")

	addrAsBytes, err := parse(addr.(string))
	if err != nil {
		return err
	}

	exists, err := conn.DeleteAddressWithValidation(name.(string), keyPrefix.(string), setAsHardcoded, addrAsBytes, !conn.Strict, prettify, addrIsLess)
	if err != nil {
		return err
	}

	if !exists {
		log.Printf(fmt.Sprintf(
			"[WARN] Deleting resource for non-existent address with name '%s' and address '%s' in range '%s'", 
			name.(string),
			addr.(string),
			keyPrefix.(string),
		))
	} else if setAsHardcoded {
		log.Printf(fmt.Sprintf(
			"[DEBUG] Deleted hardcoded address with name '%s' and address '%s' in range '%s'", 
			name.(string),
			addr,
			keyPrefix.(string),
		))
	} else {
		log.Printf(fmt.Sprintf(
			"[DEBUG] Deleted generated address with name '%s' and address '%s' in range '%s'", 
			name.(string),
			addr,
			keyPrefix.(string),
		))
	}

	return nil
}