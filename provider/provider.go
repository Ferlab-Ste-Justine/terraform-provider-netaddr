package provider

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func init() {
	schema.DescriptionKind = schema.StringMarkdown
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Description: "Name of the etcd user that will be used to access etcd. Can alternatively be set with the ETCDCTL_USERNAME environment variable. Can also be omitted if tls certificate authentication will be used instead as the username will be infered from the certificate.",
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ETCDCTL_USERNAME", ""),
			},
			"password": &schema.Schema{
				Description: "Password of the etcd user that will be used to access etcd. Can alternatively be set with the ETCDCTL_PASSWORD environment variable. Can also be omitted if tls certificate authentication will be used instead.",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("ETCDCTL_PASSWORD", ""),
			},
			"ca_cert": &schema.Schema{
				Description: "File that contains the CA certificate that signed the etcd servers' certificates. Can alternatively be set with the ETCDCTL_CACERT environment variable. Can also be omitted.",
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ETCDCTL_CACERT", ""),
			},
			"cert": &schema.Schema{
				Description: "File that contains the client certificate used to authentify the user. Can alternatively be set with the ETCDCTL_CERT environment variable. Can be omitted if password authentication is used.",
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ETCDCTL_CERT", ""),
			},
			"key": &schema.Schema{
				Description: "File that contains the client encryption key used to authentify the user. Can alternatively be set with the ETCDCTL_KEY environment variable. Can be omitted if password authentication is used.",
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ETCDCTL_KEY", ""),
			},
			"endpoints": &schema.Schema{
				Description: "Endpoints of the etcd servers. The entry of each server should follow the ip:port format and be coma separated. Can alternatively be set with the ETCDCTL_ENDPOINTS environment variable.",
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ETCDCTL_ENDPOINTS", ""),
			},
			"connection_timeout": &schema.Schema{
				Description: "Timeout to establish the etcd servers connection in seconds. Defaults to 10.",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10,
			},
			"request_timeout": &schema.Schema{
				Description: "Timeout for individual requests the provider makes on the etcd servers in seconds. Defaults to 10.",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10,
			},
			"retries": &schema.Schema{
				Description: "Number of times operations that result in retriable errors should be re-attempted. Defaults to 10.",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10,
			},
			"strict": &schema.Schema{
				Description: "Whether the provider should trigger a failure if resources are already existing during their creation, already absent during their deletion or otherwise absent during reads. Setting this value to false is convenient, but it might not alert you of bad failure situations (like resource name duplicates or the etcd state being tampered outside of terraform) so we recommend using this setting only to recover for failure situations that are well understood like Terraform having failed to persist its state in a previous apply.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"netaddr_address_ipv4": resourceNetAddrAddressIpv4(),
			"netaddr_address_mac": resourceNetAddrAddressMac(),
			"netaddr_range_ipv4": resourceNetAddrRangeIpv4(),
			"netaddr_range_mac": resourceNetAddrRangeMac(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"netaddr_address_list_ipv4": dataSourceNetAddrAddressListIpv4(),
			"netaddr_address_list_mac": dataSourceNetAddrAddressListMac(),
			"netaddr_address_ipv4": dataSourceNetAddrAddressIpv4(),
			"netaddr_address_mac": dataSourceNetAddrAddressMac(),
			"netaddr_range_ipv4": dataSourceNetAddrRangeIpv4(),
			"netaddr_range_mac": dataSourceNetAddrRangeMac(),
			"netaddr_range_usage_ipv4": dataSourceNetAddrRangeUsageIpv4(),
		},
		ConfigureFunc: providerConfigure,
		//Should implement close once this issue is resolved: https://github.com/hashicorp/terraform-plugin-sdk/issues/63
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	endpoints, _ := d.Get("endpoints").(string)
	username, _ := d.Get("username").(string)
	password, _ := d.Get("password").(string)
	caCert, _ := d.Get("ca_cert").(string)
	cert, _ := d.Get("cert").(string)
	key, _ := d.Get("key").(string)
	connectionTimeout, _ := d.Get("connection_timeout").(int)
	requestTimeout, _ := d.Get("request_timeout").(int)
	retries, _ := d.Get("retries").(int)
	strict, _ := d.Get("strict").(bool)
	tlsConf := &tls.Config{}

	if cert != "" {
		certData, err := tls.LoadX509KeyPair(cert, key)
		if err != nil {
			return nil, err
		}
		(*tlsConf).Certificates = []tls.Certificate{certData}
		(*tlsConf).InsecureSkipVerify = false
	}

	if caCert != "" {
		caCertContent, err := ioutil.ReadFile(caCert)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Failed to read root certificate file: %s", err.Error()))
		}
		roots := x509.NewCertPool()
		ok := roots.AppendCertsFromPEM(caCertContent)
		if !ok {
			return nil, errors.New("Failed to parse root certificate authority")
		}
		(*tlsConf).RootCAs = roots
	}

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   strings.Split(endpoints, ","),
		Username:    username,
		Password:    password,
		TLS:         tlsConf,
		DialTimeout: time.Duration(connectionTimeout) * time.Second,
	})

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to connect to etcd servers: %s", err.Error()))
	}

	return EtcdConnection{
		Client:  cli,
		Timeout: requestTimeout,
		Retries: retries,
		Strict: strict,
	}, nil
}
