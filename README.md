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