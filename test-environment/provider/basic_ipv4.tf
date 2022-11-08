//Ipv4 Validation
resource "netaddr_range_ipv4" "basic_ipv4" {
    key_prefix = "/test/basic-ipv4/"
    first_address = "192.168.0.1"
    last_address = "192.168.2.254"
}

resource "netaddr_address_ipv4" "basic_ipv4_addr1" {
    range_id = netaddr_range_ipv4.basic_ipv4.id
    name = "addr1"
}

resource "netaddr_address_ipv4" "basic_ipv4_addr2" {
    range_id = netaddr_range_ipv4.basic_ipv4.id
    name = "addr2"
}

output "basic_ipv4_addr1" {
  value = netaddr_address_ipv4.basic_ipv4_addr1.address
}

output "basic_ipv4_addr2" {
  value = netaddr_address_ipv4.basic_ipv4_addr2.address
}
