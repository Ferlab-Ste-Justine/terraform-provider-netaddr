package provider

import(
	"github.com/Ferlab-Ste-Justine/terraform-provider-netaddr/address"

	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNetAddrAddressCreate(d *schema.ResourceData, meta interface{}, rangeType string, parse address.ParseAddr, prettify address.PrettifyAddr, incAddr address.IncrementAddress, addrIsGreater address.AddressIsGreater) error {
	conn := meta.(address.EtcdConnection)
	name := d.Get("name")
	keyPrefix := d.Get("range_id")
	hAddr, setAsHardcoded := d.GetOk("hardcoded_address")

	if setAsHardcoded {
		addrAsBytes, err := parse(hAddr.(string))
		if err != nil {
			return err
		}

		exists, _, genErr := conn.GenerateHardcodedAddressWithValidation(name.(string), []string{keyPrefix.(string)}, addrAsBytes, rangeType, !conn.Strict, prettify)
		if genErr != nil {
			return genErr
		}

		if exists {
			log.Printf(fmt.Sprintf(
				"[WARN] Creating resource for pre-existing hardcoded address of type '%s', name '%s' and address '%s' in range '%s'", 
				rangeType,
				name.(string),
				hAddr.(string),
				keyPrefix.(string),
			))
		} else {
			log.Printf(fmt.Sprintf(
				"[DEBUG] Created hardcoded address of type '%s', name '%s' and address '%s' in range '%s'", 
				rangeType,
				name.(string),
				hAddr.(string),
				keyPrefix.(string),
			))
		}
	} else {
		exists, addr, _, genErr := conn.GenerateGeneratedAddressWithValidation(name.(string), []string{keyPrefix.(string)}, rangeType, !conn.Strict, addrIsGreater, incAddr)
		if genErr != nil {
			return genErr
		}

		if exists {
			log.Printf(fmt.Sprintf(
				"[WARN] Creating resource for pre-existing generated address of type '%s', name '%s' and address '%s' in range '%s'", 
				rangeType,
				name.(string),
				prettify(addr),
				keyPrefix.(string),
			))
		} else {
			log.Printf(fmt.Sprintf(
				"[DEBUG] Created generated address of type '%s', name '%s' and address '%s' in range '%s'", 
				rangeType,
				name.(string),
				prettify(addr),
				keyPrefix.(string),
			))
		}
	}
	
	d.SetId(name.(string))
	return resourceNetAddrAddressRead(d, meta, rangeType, prettify)
}

func resourceNetAddrAddressRead(d *schema.ResourceData, meta interface{}, rangeType string, prettify address.PrettifyAddr) error {
	conn := meta.(address.EtcdConnection)
	name := d.Get("name")
	keyPrefix := d.Get("range_id")

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

func resourceNetAddrAddressDelete(d *schema.ResourceData, meta interface{}, parse address.ParseAddr, prettify address.PrettifyAddr, addrIsLess address.AddressIsLess) error {
	conn := meta.(address.EtcdConnection)
	name := d.Get("name")
	keyPrefix := d.Get("range_id")
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
