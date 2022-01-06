# About

This terraform providers allows you to automatically assign ipv4 and mac addresses within a specified range as terraform resources.

The following are also supported: 
- Specifying some hardcoded addresses in addition to the automatically assigned ones
- Accessing the address values from any other terraform environment by using data sources

Etcd is used as a store to manage the state of the assigned addresses.