# nrf-client-emulator
This project enables Arista EOS and MikroTik RouterOS to detect when an address becomes unreachable (similar to Cisco's ip-sla icmp-echo) and shut down a specified range of interfaces, accordingly (and restore them when reachability returns).

It is meant to be used as a client to enable [NRF](https://github.com/rwinkhart/network-redundancy-fuzzer) to work on non-Cisco hardware/software.

If using only Arista hardware in your network, it would be better to manipulate Arista devices from a remote server using Arista's [goeapi](https://github.com/aristanetworks/goeapi). This project leverages eAPI in a less efficient way for the sake of compatibility with NRF/Cisco devices.

Clients for other vendors _may_ be implemented in the future.

# How-To

## Arista
1. Download (binaries not yet available) or compile nrf-client-emulator for your target hardware.
2. Install the binary on the target system at `/mnt/flash/nrf-client`. You may want to copy it over using `scp` from `bash` mode.
3. Ensure the binary is executable (`chmod +x /mnt/flash/nrf-client` in `bash` mode).
4. Enable eAPI access from `config terminal`.
    - `management api http-commands`
    - `protocol http`
    - `no shutdown`
    - `end`
    - Verify with `sh management api http-commands`
5. Configure eapi.conf.
    - A template can be generated simply by running `bash /mnt/flash/nrf-client` for the first time.
    - Modify settings in the generated `/mnt/flash/eapi.conf` as needed. You may want to use `vi` from `bash` mode.
6. Configure an event handler for nrf-client from `config terminal`.
    - `event-handler nrf`
    - `trigger on-boot`
    - `asynchronous`
    - `action bash /mnt/flash/nrf-client <nrf server ip> <interval (seconds)> <fail range, e.g. Et2-4,Et7>`
    - `end`
    - Verify with `sh event-handler`
7. `reload`; nrf-client will start automatically following the next boot.

## MikroTik
1. Copy the contents of nrf-client.mikrotik into your RouterOS CLI to store the script.
2. Specify interfaces to manage using `:global nrfFailGroup`.
   - e.g. `:global nrfFailGroup [:toarray "ether2,ether3,ether4"]`
3. Schedule the script to run at your desired interval using `/system/scheduler/add`.
    - e.g. `/system/scheduler/add name=nrfScheduler interval=00:00:05 on-event=nrf-client`
4. That's it! The script begins running immediately.
