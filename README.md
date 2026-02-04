# About

This terraform providers allows you to automatically assign ipv4 and mac addresses within a specified range as terraform resources.

The following are also supported: 
- Specifying some hardcoded addresses in addition to the automatically assigned ones
- Accessing the address values from any other terraform environment by using data sources

Etcd is used as a store to manage the state of the assigned addresses.

The provider can be found at: https://registry.terraform.io/providers/Ferlab-Ste-Justine/netaddr/latest

# Local Troubleshoot

Sometimes, you need to troubleshoot changes to the provider locally before publishing your modifications to the terraform registry.

You need to have both golang 1.16 and Terraform setup on your machine for this to work. This also relies on a local minikube installation for running etcd.

## Setup Terraform to Detect Your Provider

Create a file named **.terraformrc** in your home directory with the following content:

```
provider_installation {
  dev_overrides {
    "ferlab/netaddr" = "<Path to the project's parent directory on your machine>/terraform-provider-netaddr"
  }
  direct {}
}
```

## Lauch the Etcd Server

Then, launch etcd by going to **test-environment/server** and typing:

```
terraform apply
```

Etcd will be exposed on port **32379**

Copy the **test-environment/server/certs** directory to **test-environment/provider/certs**.

## Build the Provider

Go to the root directory of this project and type:

```
go build .
```

## play with the Provider

From there, you can go to the **test-environment/provider** directory, edit the terraform scripts as you wish and experiment with the provider.

Note that you should not do **terraform init**. The provider was already setup globally in a previous step for your user.

# Etcd Address Modelization

## Key Space

All the etcd keys managed by the provider are scoped with a given key prefix passed by the user.

Addresses are modeled in etcd with the following two constructs: **address ranges** and **addresses**.

**Address ranges** have the following entries:
- **Type**: 
  - **key**: `<user prefix>info/type`
  - **description**: Identifies the type of address the range manages. Currently can be **ipv4** or **mac**. This key doesn't change.
- **FirstAddress**: 
  - **key**: `<user prefix>info/firstaddr`
  - **description**: First address in the range. This key doesn't change.
- **LastAddress**:
  - **key**: `<user prefix>info/lastaddr`
  - **description**: Last address in the range. This key doesn't change.
- **NextAddress**: rangePrefix + "data/nextaddr",
  - **key**: `<user prefix>data/nextaddr`
  - **description**: Pointer keeping track of the next generated address to return. It is monotonically increasing, starting at **FirstAddress** and never exceeding **LastAddress**.

**Addresses** have the following entries:
- **Name**:
  - **key**: `<user prefix>data/name/<user defined name>`
  - **Content**: Address
  - **description**: Entry present for all assigned addresses giving the address for a given user defined name/label for the address.
- **GeneratedAddress**: 
  - **key**: `<user prefix>data/address/generated/<address>`
  - **content**: User defined name/label for the address.
  - **description**: Entry present for all automatically assigned addresses. Used for preventive error checking when creating hardcoded addresses and integrity checking when deleting generated addresses.
- **HardcodedAddress**: 
  - **key**: `<user prefix>data/address/hardcoded/<address>`
  - **content**: User defined name/label for the address.
  - **description**: Entry present for all user specified hardcoded addresses. Beyond error/integrity checks, it is also used when returning a generated address to determine hardcoded addresses to skip over. 
- **DeletedAddress**: 
  - **key**: `<user prefix>data/address/deleted/<address>`
  - **Content**: User defined name/label for the address.
  - **description**: Entry present for all freed addresses that are behind the **NextAddress** pointer of their range. Used to keep track of freed addresses that can be reassigned.

## Workflow

All write operations by the provider are transactional (using etcd transactions to enforce this). Either the entire operation succeeds or the entire operation fails. Barring unforeseen bugs in etcd itself (or this provider), the keyspace cannot be in an inconsistent state during the course of an operation or if it fails before completing.

There are two classes of address managed by the provider which are treated differently: **generated** addresses where the user is happy to get any non-taken address (kind of like dhcp, usually for programmatically generated machines) and **hardcoded** addresses where the user specifies a hardcoded address that is taken (kind of like static ips, usually either for legacy manually provisioned machines or for boostrap machines, like the etcd cluster used by the provider for example).

### Hardcoded Addresses

When being created, the address can be taken from anywhere in the address range (as specified by the user) though integrity checks are done to ensure that the desired address is not already assigned.

When being deleted, an entry in the deleted addresses is created for the address if it has fallen behind the **NextAddress** pointer of its address range.

### Generated Addresses

When being created, a look is taken at deleted addresses first and if any is found, it is assigned (and removed from the pool of deleted addresses). Otherwise, a look is taken at the **NextAddress** pointer of the address range to determine the next address to assign. The pointer is incremented to skip over any pre-existing hardcoded addresses until an available address is found which is then assigned (and the **NextAddress** pointer is further incremented since that address is now assigned). Should the **NextAddress** pointer exceed **LastAddress** for the address range, an error is returned as there are no more addresses available to assign.

When being deleted, an entry in the deleted addresses is created for the address (since generated addresses are always behind the **NextAddress** pointer).

# V2 Version of Ipv4 Addresses

## Note on V2 and V1

We recently added a v2 version for ipv4 addresses to support multiple available address ranges for the same address (technically, we could very easily do the same for MAC addresses, but we believe their range tend to be larger which is the case for us).

The use case (based on recent real life experience), is that the network admins assign you an ip range in a network, you exhaust that ip range and they give you another additional ip range that is part of the same network.

The V2 version handle manages addresses in many ranges so that you don't have to keep track of which address is in which range.

The etcd address key structure and logic to interact with that key structure remains unchanged for v2 relative to v1. 

The main thing that the v2 does is fetch an address from several ranges on address creation (skipping over ranges that are full) and look at all those ranges when reading the address for its data source. 

For the resource, the range the address was created in is added to the resource in the terraform state as a performance optimization (and for informative purposes) and the resource act like v1 for the remainder of its lifecycle.

All that to say that if you add a network range and start adding v2 addresses while you still have v1 addresses for the first range of a network, that will work just fine and you won't have any data conflicts.

You can also have a v1 address, access that address elsewhere with a v2 data source and it won't be a problem (the v2 data source will just look at its range and find it in the range the v1 resource put it in). However, you should note use a v1 data source for a v2 address resource (the data source will only have a single range to find the address in and may be missing one of the ranges the address was assigned to by the resource).

If you want to migrate v1 addresses to v2 without losing the addresses, you can set **retain_on_delete** to **true** for the v1 resource (which won't delete the address in the range when the terraform resource is deleted), set **manage_existing** to **true** on the replacement v2 resource (which will prevent resource creation from tiggering an error when the address is found during the creation sanity check) and you will be set. 

## Note on Adding More Ranges for V2

The migration path from V1 ipv4 address to V2 ipv4 address will also work in-place for an existing V2 ipv4 address you want to add an additional extra range to.

When you change the ip ranges of the address, it will trigger a re-creation of the resource in the terraform lifecycle, but by putting **retain_on_delete** and **manage_existing** to true, the terraform resource deletion will be a no-op for the address and the terraform resource creation won't trigger an error when the address is found (and essential be a no-op also for the address). Just make sure that you are just using this technique to add ranges and no remove them, or you might get into trouble.

Note that you can also use the above technique to migrate the management of an address between different terraform pipelines without having to change it or hardcode it.