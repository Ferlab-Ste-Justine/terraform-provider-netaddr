resource "netaddr_range_ipv4" "migration_ipv4" {
    key_prefix = "/test/migration-ipv4/"
    first_address = "192.171.0.1"
    last_address = "192.171.0.3"
}

resource "netaddr_range_ipv4" "migration_ipv4_range2" {
    key_prefix = "/test/migration-ipv4-range2/"
    first_address = "192.171.0.10"
    last_address = "192.171.0.20"
}

/*resource "netaddr_address_ipv4" "migration_ipv4_addr1" {
    range_id = netaddr_range_ipv4.migration_ipv4.id
    name = "addr1"
    retain_on_delete = true
}

output "migration_ipv4_addr1" {
  value = netaddr_address_ipv4.migration_ipv4_addr1
}

resource "netaddr_address_ipv4" "migration_ipv4_addr2" {
    range_id = netaddr_range_ipv4.migration_ipv4.id
    name = "addr2"
    retain_on_delete = true
}

output "migration_ipv4_addr2" {
  value = netaddr_address_ipv4.migration_ipv4_addr2
}

resource "netaddr_address_ipv4" "migration_ipv4_addr3" {
    range_id = netaddr_range_ipv4.migration_ipv4.id
    name = "addr3"
    hardcoded_address = "192.171.0.3"
    retain_on_delete = true
}

output "migration_ipv4_addr3" {
  value = netaddr_address_ipv4.migration_ipv4_addr3
}*/

resource "netaddr_address_ipv4_v2" "migration_ipv4_addr1" {
    range_ids = [netaddr_range_ipv4.migration_ipv4.id, netaddr_range_ipv4.migration_ipv4_range2.id]
    name = "addr1"
    manage_existing = true
}

output "migration_ipv4_addr1" {
  value = netaddr_address_ipv4_v2.migration_ipv4_addr1
}

resource "netaddr_address_ipv4_v2" "migration_ipv4_addr2" {
    range_ids = [netaddr_range_ipv4.migration_ipv4.id, netaddr_range_ipv4.migration_ipv4_range2.id]
    name = "addr2"
    manage_existing = true
}

output "migration_ipv4_addr2" {
  value = netaddr_address_ipv4_v2.migration_ipv4_addr2
}

resource "netaddr_address_ipv4_v2" "migration_ipv4_addr3" {
    range_ids = [netaddr_range_ipv4.migration_ipv4.id, netaddr_range_ipv4.migration_ipv4_range2.id]
    name = "addr3"
    hardcoded_address = "192.171.0.3"
    manage_existing = true
}

output "migration_ipv4_addr3" {
  value = netaddr_address_ipv4_v2.migration_ipv4_addr3
}