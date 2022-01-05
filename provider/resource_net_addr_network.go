package provider

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceNetAddrNetwork() *schema.Resource {
	return &schema.Resource{
		Description: "Network to create ipv4 and mac addresses on.",
		Create: resourceNetAddrNetworkCreate,
		Read:   resourceNetAddrNetworkRead,
		Delete: resourceNetAddrNetworkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"key_prefix": {
				Description: "Etcd key prefix for all the keys related to the network.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
            "ipv4": {
                Type:     schema.TypeMap,
                Required: true,
                ForceNew: true,
                ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
                    v := val.(map[string]interface{})

                    _, firstOk := v["first"]
                    _, lastOk := v["last"]

                    if (!firstOk) || (!lastOk) {
                        return []string{}, []error{errors.New("For the ipv4 field, the following keys need to be defined: first, last")}
                    }

                    return []string{}, []error{}
                },
            },
            "ipv6": {
                Type:     schema.TypeMap,
                Required: true,
                ForceNew: true,
                ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
                    v := val.(map[string]interface{})

                    _, firstOk := v["first"]
                    _, lastOk := v["last"]

                    if (!firstOk) || (!lastOk) {
                        return []string{}, []error{errors.New("For the ipv6 field, the following keys need to be defined: first, last")}
                    }

                    return []string{}, []error{}
                },
            },
            "mac": {
                Type:     schema.TypeMap,
                Required: true,
                ForceNew: true,
                ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
                    v := val.(map[string]interface{})

                    _, firstOk := v["first"]
                    _, lastOk := v["last"]

                    if (!firstOk) || (!lastOk) {
                        return []string{}, []error{errors.New("For the mac field, the following keys need to be defined: first, last")}
                    }

                    return []string{}, []error{}
                },
            },
		},
	}
}

func networkSchemaToModel(d *schema.ResourceData) (Network, error) {
	ipv4, _ := d.GetOk("ipv4")
	ipv4Cast := ipv4.(map[string]interface{})
	ipv6, _ := d.GetOk("ipv6")
	ipv6Cast := ipv6.(map[string]interface{})
	mac, _ := d.GetOk("mac")
	macCast := mac.(map[string]interface{})

	firstIpv4, _ := ipv4Cast["first"].(string)
	firstIpv4Bytes, firstIpv4Err := Ipv4StringToBytes(firstIpv4)
	if firstIpv4Err != nil {
		return Network{}, firstIpv4Err
	}

	lastIpv4, _ := ipv4Cast["last"].(string)
	lastIpv4Bytes, lastIpv4Err := Ipv4StringToBytes(lastIpv4)
	if lastIpv4Err != nil {
		return Network{}, lastIpv4Err
	}

	firstIpv6, _ := ipv6Cast["first"].(string)
	firstIpv6Bytes, firstIpv6Err := Ipv6StringToBytes(firstIpv6)
	if firstIpv6Err != nil {
		return Network{}, firstIpv6Err
	}

	lastIpv6, _ := ipv6Cast["last"].(string)
	lastIpv6Bytes, lastIpv6Err := Ipv6StringToBytes(lastIpv6)
	if lastIpv6Err != nil {
		return Network{}, lastIpv6Err
	}

	firstMac, _ := macCast["first"].(string)
	firstMacBytes, firstMacErr := MacStringToBytes(firstMac)
	if firstMacErr != nil {
		return Network{}, firstMacErr
	}

	lastMac, _ := macCast["last"].(string)
	lastMacBytes, lastMacErr := MacStringToBytes(lastMac)
	if lastMacErr != nil {
		return Network{}, lastMacErr
	}

	return Network{
		Ipv4: AddressRange{
			First: firstIpv4Bytes,
			Last: lastIpv4Bytes,
		},
		Ipv6: AddressRange{
			First: firstIpv6Bytes,
			Last: lastIpv6Bytes,
		},
		Mac: AddressRange{
			First: firstMacBytes,
			Last: lastMacBytes,
		},
	}, nil
}

func resourceNetAddrNetworkCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(EtcdConnection)
	keyPrefix, _ := d.GetOk("key_prefix")

	network, err := networkSchemaToModel(d)
	if err != nil {
		return errors.New(fmt.Sprintf("Error creating network: %s", err.Error()))
	}

	err = conn.CreateNetwork(keyPrefix.(string), network)
	if err != nil {
		return errors.New(fmt.Sprintf("Error creating network: %s", err.Error()))
	}

	d.SetId(keyPrefix.(string))
	return resourceNetAddrNetworkRead(d, meta)
}

func resourceNetAddrNetworkRead(d *schema.ResourceData, meta interface{}) error {
	keyPrefix := d.Id()
	conn := meta.(EtcdConnection)

	network, exists, err := conn.GetNetworkInfo(keyPrefix)

	if !exists {
		return errors.New(fmt.Sprintf("Error retrieving network at prefix '%s': Network does not exist", keyPrefix))
	}

	if err != nil {
		return errors.New(fmt.Sprintf("Error retrieving network at prefix '%s': %s", keyPrefix, err.Error()))
	}

	if Ipv4BytesToString(network.Ipv4.First) == "<nil>" {
		return errors.New(fmt.Sprintf("Error retrieving network at prefix '%s': First ipv4 address is malformed", keyPrefix))
	}

	if Ipv4BytesToString(network.Ipv4.Last) == "<nil>" {
		return errors.New(fmt.Sprintf("Error retrieving network at prefix '%s': Last ipv4 address is malformed", keyPrefix))
	}

	if Ipv6BytesToString(network.Ipv6.First) == "<nil>" {
		return errors.New(fmt.Sprintf("Error retrieving network at prefix '%s': First ipv6 address is malformed", keyPrefix))
	}

	if Ipv6BytesToString(network.Ipv6.Last) == "<nil>" {
		return errors.New(fmt.Sprintf("Error retrieving network at prefix '%s': Last ipv6 address is malformed", keyPrefix))
	}

	d.Set("key_prefix", keyPrefix)
	d.Set("ipv4", map[string]string{
		"first": Ipv4BytesToString(network.Ipv4.First),
		"last": Ipv4BytesToString(network.Ipv4.Last),
	})
	d.Set("ipv6", map[string]string{
		"first": Ipv6BytesToString(network.Ipv6.First),
		"last": Ipv6BytesToString(network.Ipv6.Last),
	})
	d.Set("mac", map[string]string{
		"first": MacBytesToString(network.Mac.First),
		"last": MacBytesToString(network.Mac.Last),
	})

	return nil
}

func resourceNetAddrNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(EtcdConnection)
	keyPrefix, _ := d.GetOk("key_prefix")

	err := conn.DestroyNetwork(keyPrefix.(string))
	if err != nil {
		return errors.New(fmt.Sprintf("Error destroying network: %s", err.Error()))
	}

	return nil
}