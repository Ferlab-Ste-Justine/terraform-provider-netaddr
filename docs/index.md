---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "netaddr Provider"
subcategory: ""
description: |-
  
---

# netaddr Provider



## Example Usage

```terraform
provider "etcd" {
  endpoints = "127.0.0.1:32379"
  ca_cert = "${path.module}/certs/ca.pem"
  cert = "${path.module}/certs/root.pem"
  key = "${path.module}/certs/root.key"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- **ca_cert** (String) File that contains the CA certificate that signed the etcd servers' certificates. Can alternatively be set with the ETCDCTL_CACERT environment variable. Can also be omitted.
- **cert** (String) File that contains the client certificate used to authentify the user. Can alternatively be set with the ETCDCTL_CERT environment variable. Can be omitted if password authentication is used.
- **connection_timeout** (Number) Timeout to establish the etcd servers connection in seconds. Defaults to 10.
- **endpoints** (String) Endpoints of the etcd servers. The entry of each server should follow the ip:port format and be coma separated. Can alternatively be set with the ETCDCTL_ENDPOINTS environment variable.
- **key** (String) File that contains the client encryption key used to authentify the user. Can alternatively be set with the ETCDCTL_KEY environment variable. Can be omitted if password authentication is used.
- **password** (String, Sensitive) Password of the etcd user that will be used to access etcd. Can alternatively be set with the ETCDCTL_PASSWORD environment variable. Can also be omitted if tls certificate authentication will be used instead.
- **request_timeout** (Number) Timeout for individual requests the provider makes on the etcd servers in seconds. Defaults to 10.
- **retries** (Number) Number of times operations that result in retriable errors should be re-attempted. Defaults to 10.
- **username** (String) Name of the etcd user that will be used to access etcd. Can alternatively be set with the ETCDCTL_USERNAME environment variable. Can also be omitted if tls certificate authentication will be used instead as the username will be infered from the certificate.