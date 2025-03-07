data "netaddr_range_ipv4" "test" {
    key_prefix = "/test/ipv4/"
}

data "netaddr_range_keyspace_ipv4" "test" {
  range_id = data.netaddr_range_ipv4.test.id
}

output "first_address" {
  description = "The first address of the range."
  value       = data.netaddr_range_keyspace_ipv4.test.first_address
}

output "last_address" {
  description = "The last address of the range."
  value       = data.netaddr_range_keyspace_ipv4.test.last_address
}

output "next_address" {
  description = "The next never used before address that will be assigned from the range. Note that any deleted address will be assigned first"
  value       = data.netaddr_range_keyspace_ipv4.test.next_address
}

output "addresses" {
  description = "List of all addresses (taken from the keyspace listing addresses by name)"
  value       = data.netaddr_range_keyspace_ipv4.test.addresses
}

output "generated_addresses" {
  description = "List of all addresses that we dynamically allocated from the range (taken from the keyspace reserved to that effect)"
  value       = data.netaddr_range_keyspace_ipv4.test.generated_addresses
}

output "hardcoded_addresses" {
  description = "List of all addresses whose value was hardcoded in the resource's input from the range (taken from the keyspace reserved to that effect)"
  value       = data.netaddr_range_keyspace_ipv4.test.hardcoded_addresses
}

output "deleted_addresses" {
  description = "List of all previously allocated freed addresses that can now be allocated (taken from the keyspace reserved to that effect)"
  value       = data.netaddr_range_keyspace_ipv4.test.deleted_addresses
}