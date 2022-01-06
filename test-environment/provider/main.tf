resource "netaddr_range" "ipv4" {
    key_prefix = "/test/ipv4/"
    first_address = "192.168.0.1"
    last_address = "192.168.2.254"
    type = "ipv4"
}

resource "netaddr_range" "mac" {
    key_prefix = "/test/mac/"
    first_address = "52:54:00:00:00:00"
    last_address = "52:54:00:ff:ff:ff"
    type = "mac"
}

resource "netaddr_address" "ipv4_test" {
    range_id = netaddr_range.ipv4.id
    name = "test"
}

resource "netaddr_address" "ipv4_test2" {
    range_id = netaddr_range.ipv4.id
    name = "test2"
}

resource "netaddr_address" "mac_test" {
    range_id = netaddr_range.mac.id
    name = "test3"
}

resource "netaddr_address" "mac_test2" {
    range_id = netaddr_range.mac.id
    name = "test4"
}

resource "netaddr_address" "mac_test3" {
    range_id = netaddr_range.mac.id
    name = "test5"
    hardcoded_address = "52:54:00:00:00:02"
}

resource "netaddr_address" "mac_test4" {
    range_id = netaddr_range.mac.id
    name = "test6"
}

resource "netaddr_address" "mac_test5" {
    range_id = netaddr_range.mac.id
    name = "test7"
}

output "test" {
  value = netaddr_address.ipv4_test.address
}

output "test2" {
  value = netaddr_address.ipv4_test2.address
}

output "test3" {
  value = netaddr_address.mac_test.address
}

output "test4" {
  value = netaddr_address.mac_test2.address
}

output "test5" {
  value = netaddr_address.mac_test3.address
}

output "test6" {
  value = netaddr_address.mac_test4.address
}

output "test7" {
  value = netaddr_address.mac_test5.address
}