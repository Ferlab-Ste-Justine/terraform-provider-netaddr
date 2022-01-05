package provider

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceNetAddrIpv4() *schema.Resource {
	return &schema.Resource{
		Description: "Ipv4 address.",
		Create: resourceNetAddrIpv4Create,
		Read:   resourceNetAddrIpv4Read,
		Delete: resourceNetAddrIpv4Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Name to associate with the address.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"network_id": {
				Description: "Identifier of the network the address is tied to.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"hardcoded_address": {
				Description: "An optional input to fixate the address to a specific value.",
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"address": {
				Description: "The address that got assigned to the resource.",
				Type:         schema.TypeString,
				Computed:     true,
			},
		},
	}
}

func resourceNetAddrIpv4Create(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(EtcdConnection)
	name, _ := d.GetOk("name")
	networkPrefix, _ := d.GetOk("network_id")
	hAddr, hAddrExists := d.GetOk("hardcoded_address")

	if hAddrExists {
		addrAsBytes, err := Ipv4StringToBytes(hAddr.(string))
		if err != nil {
			return err
		}

		err = conn.CreateHardcodedIpv4Address(networkPrefix.(string), name.(string), addrAsBytes)
		if err != nil {
			return err
		}
	} else {
		_, err := conn.CreateGeneratedIpv4Address(networkPrefix.(string), name.(string))
		if err != nil {
			return err
		}
	}

	d.SetId(name.(string))
	return resourceNetAddrIpv4Read(d, meta)
}

func resourceNetAddrIpv4Read(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(EtcdConnection)
	name, _ := d.GetOk("name")
	networkPrefix, _ := d.GetOk("network_id")

	addr, err := conn.GetIpv4Address(networkPrefix.(string), name.(string))
	if err != nil {
		return err
	}
	d.Set("address", Ipv4BytesToString(addr))

	return nil
}

func resourceNetAddrIpv4Delete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(EtcdConnection)
	name, _ := d.GetOk("name")
	networkPrefix, _ := d.GetOk("network_id")
	_, hAddrExists := d.GetOk("hardcoded_address")
	addr, addrExists := d.GetOk("address")

	if !addrExists {
		return errors.New(fmt.Sprintf("Cannot delete address named '%s': Address is missing", name.(string)))
	}

	addrAsBytes, err := Ipv4StringToBytes(addr.(string))
	if err != nil {
		return err
	}

	if hAddrExists {
		err := conn.DeleteHardcodedIpv4Address(networkPrefix.(string), name.(string), addrAsBytes)
		if err != nil {
			return err
		}
	} else {
		err := conn.DeleteGeneratedIpv4Address(networkPrefix.(string), name.(string), addrAsBytes)
		if err != nil {
			return err
		}
	}

	return nil
}
