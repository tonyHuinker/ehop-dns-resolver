# ehop-dns-resolver  

Script will start by asking you for a key file (see example key file in this repo or providing your ExtraHop hostname and api key)

Next it will ask how much tie your would like to wait (in ms) between every DNS call... This is to avoid putting any stress on your DNS server to avoid setting off any security alarms.

It will then present you with some summary statistics about devices on your system... Like so.

>Total L3 Devices = 1616
>Devices with DNS Names = 1039
>Devices with DHCP Names = 738
>Devices with CustomNames and no DNS name = 5

Next it will ask you some questions about how you would like it to run...

>Would you like to  

> [1] -- Resolve names for all devices that do not have DNS names, WITHOUT overwriting existing CustomNames  

> [2] -- Resolve names for all devices that do not have DNS names, WHILE overwriting existing CustomNames  

> [3] -- Resolve names for all devices that do not have DNS names, and ASK before overwriting CustomNames  

> [4] -- Do a dry run, and show results without actually making any changes  

> [5] -- Delete all custom names

Make your choice, and let it run.
