# nrf-client-emulator
This project enables Arista EOS to detect when an address becomes unreachable (similar to Cisco's ip-sla icmp-echo) and shut down a specified range of interfaces, accordingly (and restore them when reachability returns).

It is meant to be used as a client to enable [NRF](https://github.com/rwinkhart/network-redundancy-fuzzer) to work on non-Cisco hardware/software.

If using only Arista hardware in your network, it would be better to manipulate Arista devices from a remote server using Arista's [goeapi](https://github.com/aristanetworks/goeapi). This project leverages eAPI in a less efficient way for the sake of compatibility with NRF/Cisco devices.

Clients for other vendors may be implemented in the future.
