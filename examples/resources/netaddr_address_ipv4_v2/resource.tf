resource "netaddr_range_ipv4" "range1" {
    key_prefix = "/test/ipv4/"
    first_address = "192.168.0.1"
    last_address = "192.168.0.64"
}

resource "netaddr_range_ipv4" "range2" {
    key_prefix = "/test/ipv4-extras/"
    first_address = "192.168.0.128"
    last_address = "192.168.0.254"
}

resource "netaddr_address_ipv4_v2" "test" {
    range_ids = [netaddr_range_ipv4.range1.id, netaddr_range_ipv4.range2.id]
    name = "test"
    hardcoded_address = "192.168.0.5"
}

resource "netaddr_address_ipv4_v2" "test2" {
    range_ids = [netaddr_range_ipv4.range1.id, netaddr_range_ipv4.range2.id]
    name = "test2"
}

output "test_addr" {
  value = netaddr_address_ipv4.test.address
}

output "test2_addr" {
  value = netaddr_address_ipv4.test2.address
}

//v1 to v2 Migration example
//Note that you can also use the retain_on_delete / manage_existing notation to add ranges in-place to v2 addresses

resource "netaddr_range_ipv4" "original" {
    key_prefix = "/original/ipv4/"
    first_address = "192.168.0.1"
    last_address = "192.168.0.64"
}

resource "netaddr_range_ipv4" "extra" {
    key_prefix = "/extra/ipv4-extras/"
    first_address = "192.168.0.128"
    last_address = "192.168.0.254"
}

//Removed content before migration
/*resource "netaddr_address_ipv4" "ip" {
    range_id = netaddr_range_ipv4.original.id
    name = "ip"
    hardcoded_address = "192.168.0.5"
    retain_on_delete = true
}

resource "netaddr_address_ipv4" "ip" {
    range_id = netaddr_range_ipv4.original.id
    name = "ip2"
    retain_on_delete = true
}*/

//Content after migration
resource "netaddr_address_ipv4_v2" "ip" {
    range_ids = [netaddr_range_ipv4.original.id, netaddr_range_ipv4.extra.id]
    name = "ip"
    hardcoded_address = "192.168.0.5"
    manage_existing = true
}

resource "netaddr_address_ipv4_v2" "ip" {
    range_ids = [netaddr_range_ipv4.original.id, netaddr_range_ipv4.extra.id]
    name = "ip2"
    manage_existing = true
}