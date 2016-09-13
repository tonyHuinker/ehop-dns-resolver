# ehop-dns-resolver  

  This script will grab all L3 devices found in the ExtraHop system, check to see if they don't have a custom, DNS, or DHCP name.  If it has none of those three, it will make a DNS request and update the device name to that of the answer of the DNS request.

  This is useful if you are not seeing certain segments of DNS traffic and ExtraHop contains a lot of devices similar to "VMWare 10.1.1.1250" for example.

  If you wish to compile from source, go ahead.. if you wish to download the binary, you can follow these steps.

  1. Download either Windows or OSX binary from bin folder
  2. Download keys file
  3. Replace IP/Hostname in the keys file to that of your own ExtraHop IP/hostname.
  4. Replace the api key to that of your own API key
  5. Run program

  The program will ask a series of questions...

  Question 1. What is the name of your keys file?
  (default answer would be keys here)

  Question 2. Do you actually want to make changes?  
  (enter "yes" to update the device names in ExtraHop... enter test to see the changes but not actually make them)

  Question 3. Fast or slow?
  By default.. the script will make the DNS requests as fast as it can... if you are worried about over loading the DNS server whatsoever... enter "slow" and the script will wait 500ms between each call.
