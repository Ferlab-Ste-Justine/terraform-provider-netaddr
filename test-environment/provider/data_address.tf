data "netaddr_range_ipv4" "basic_ipv4" {
    key_prefix = netaddr_range_ipv4.basic_ipv4.id
}

data "netaddr_range_mac" "basic_mac" {
    key_prefix = netaddr_range_mac.basic_mac.id
}

data "netaddr_address_mac" "basic_mac_addr1" {
    range_id = data.netaddr_range_mac.basic_mac.id
    name = "addr1"
    depends_on = [netaddr_address_mac.basic_mac_addr1]
}

data "netaddr_address_mac" "basic_mac_addr2" {
    range_id = data.netaddr_range_mac.basic_mac.id
    name = "addr2"
    depends_on = [netaddr_address_mac.basic_mac_addr2]
}

output "data_range_basic_ipv4" {
  value = "first_address: ${data.netaddr_range_ipv4.basic_ipv4.first_address}, last_address: ${data.netaddr_range_ipv4.basic_ipv4.last_address}"
}

output "data_range_basic_mac" {
  value = "first_address: ${data.netaddr_range_mac.basic_mac.first_address}, last_address: ${data.netaddr_range_mac.basic_mac.last_address}"
}

output "data_address_basic_mac_addr1" {
  value = data.netaddr_address_mac.basic_mac_addr1.address
}

output "data_address_basic_mac_addr2" {
  value = data.netaddr_address_mac.basic_mac_addr2.address
}