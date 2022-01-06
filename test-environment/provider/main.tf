//Ipv4 Validation
resource "netaddr_range_ipv4" "test" {
    key_prefix = "/test/ipv4/"
    first_address = "192.168.0.1"
    last_address = "192.168.2.254"
}

resource "netaddr_address_ipv4" "test" {
    range_id = netaddr_range_ipv4.test.id
    name = "test"
}

resource "netaddr_address_ipv4" "test2" {
    range_id = netaddr_range_ipv4.test.id
    name = "test2"
}

output "ipv4_test" {
  value = netaddr_address_ipv4.test.address
}

output "ipv4_test2" {
  value = netaddr_address_ipv4.test2.address
}

//Mac Validation
resource "netaddr_range_mac" "test" {
    key_prefix = "/test/mac/"
    first_address = "52:54:00:00:00:00"
    last_address = "52:54:00:ff:ff:ff"
}

resource "netaddr_address_mac" "test" {
    range_id = netaddr_range_mac.test.id
    name = "test"
}

resource "netaddr_address_mac" "test2" {
    range_id = netaddr_range_mac.test.id
    name = "test2"
}

resource "netaddr_address_mac" "test3" {
    range_id = netaddr_range_mac.test.id
    name = "test3"
    hardcoded_address = "52:54:00:00:00:02"
}

resource "netaddr_address_mac" "test4" {
    range_id = netaddr_range_mac.test.id
    name = "test4"
}

resource "netaddr_address_mac" "test5" {
    range_id = netaddr_range_mac.test.id
    name = "test5"
}

output "mac_test" {
  value = netaddr_address_mac.test.address
}

output "mac_test2" {
  value = netaddr_address_mac.test2.address
}

output "mac_test3" {
  value = netaddr_address_mac.test3.address
}

output "mac_test4" {
  value = netaddr_address_mac.test4.address
}

output "mac_test5" {
  value = netaddr_address_mac.test5.address
}