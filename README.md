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
