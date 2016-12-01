package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tonyHuinker/ehop"
)

func askForInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(prompt)
	response, _ := reader.ReadString('\n')
	fmt.Println("\nThank You")
	return strings.TrimSpace(response)
}

func changeDeviceName(myhop *ehop.EDA, newName string, oldName string, ID int, makeChange bool, ask bool) {
	if makeChange {
		if ask{
			answer := askForInput("Would you like to overwrite Custom Name " + oldName + " to " + newName + "?\ny/n\n>>")
			switch answer {
			case "y":
				body := `{"custom_name": "` + newName + `"}`
				_, err := ehop.CreateEhopRequest("PATCH", "devices/"+strconv.Itoa(ID), body, myhop)
				if err != nil {
					fmt.Println("Error making device update call")
				} else {
					fmt.Println("Updated device succesfully from " + oldName + " to " + newName)
				}
			case "n":
				fmt.Println("Ok..  skipping this one")
			default:
				fmt.Println("Incorrect response... skipping this one")
			}
		}else{
				body := `{"custom_name": "` + newName + `"}`
				_, err := ehop.CreateEhopRequest("PATCH", "devices/"+strconv.Itoa(ID), body, myhop)
				if err != nil {
					fmt.Println("Error making device update call")
				} else {
					fmt.Println("Updated device succesfully from " + oldName + " to " + newName)
				}
		}
	}else{
		fmt.Println("**Test** would have changed " + oldName + " to " + newName)
	}
}
func showStats(myhop *ehop.EDA){
	//variables used for this function
	numL3, numDNSName, numCustomName, numDhcpName := 0, 0, 0, 0

	//Get all devices from the system
	resp, _ := ehop.CreateEhopRequest("GET", "devices?active_from=-86400000", "null", myhop)
	defer resp.Body.Close()
	var devices []ehop.Device
	//Put into struct
	error := json.NewDecoder(resp.Body).Decode(&devices)
	if error != nil {
		fmt.Println(error.Error())
		os.Exit(-1)
	}


	for _, device := range devices {
		if device.IsL3 {
			numL3++
			if device.CustomName != "" && device.DNSName == "" {
				//fmt.Println(device.CustomName + " is a custom name for this IP " + device.Ipaddr4)
				numCustomName++
			}
			if device.DNSName != "" {
				numDNSName++
			}
			if device.DhcpName != "" {
				numDhcpName++
			}
		}
	}
	fmt.Printf("Total L3 Devices = %d\nDevices with DNS Names = %d\nDevices with DHCP Names = %d\nDevices with CustomNames and no DNS name = %d\n", numL3, numDNSName, numDhcpName, numCustomName)
}

func showChoices(){
	fmt.Printf("Would you like to\n [%d] -- Resolve names for all devices that do not have DNS names, WITHOUT overwriting existing CustomNames\n [%d] -- Resolve names for all devices that do not have DNS names, WHILE overwriting existing CustomNames\n [%d] -- Resolve names for all devices that do not have DNS names, and ASK before overwriting CustomNames\n [%d] -- Do a dry run, and show results without actually making any changes\n [%d] -- Delete all custom names\n", 1, 2, 3, 4, 5)
}
func removeCustomName(myhop *ehop.EDA, newName string, oldName string, ID int, makeChange bool) {
	if makeChange {
		body := `{"custom_name": ""}`
		_, err := ehop.CreateEhopRequest("PATCH", "devices/"+strconv.Itoa(ID), body, myhop)
		if err != nil {
			fmt.Println("Error making device update call")
		} else {
			fmt.Println("Removed custom name for" + oldName)
		}
	} else {
		fmt.Println("**Test** would have removed name from " + oldName)
	}
}

func main() {
	var makeChanges = false
	var overWriteCustom = false
	var removeCustom = false
	var askForCustom = false
	//Specify Key File

	keyFile := askForInput("What is the name of your keyFile?")
	myhop := ehop.NewEDAfromKey(keyFile)

	//Get time to wait between DNS requests
	time2 := askForInput("How much time (in ms) to wait between each DNS call? Recommend 500ms or greater")
	timeInMs,_ := strconv.Atoi(time2)

	showStats(myhop)
	showChoices()
	ask := askForInput(">>")
	switch ask {
	case "1":
		overWriteCustom = false
		makeChanges = true

	case "2":
		overWriteCustom = true
		makeChanges = true
		askForCustom = false
	case "3":
		askForCustom = true
		makeChanges = true
		overWriteCustom = false
	case "4":
		makeChanges = false
	case "5":
		removeCustom = true
	default:
		fmt.Println("Exiting...")
		os.Exit(-1)

	}

	//Get all devices from the system
	resp, _ := ehop.CreateEhopRequest("GET", "devices?active_from=-86400000", "null", myhop)
	defer resp.Body.Close()
	var devices []ehop.Device

	//Put into struct
	error := json.NewDecoder(resp.Body).Decode(&devices)
	if error != nil {
		fmt.Println(error.Error())
		os.Exit(-1)
	}

	//Grab all L3 devices... put into an array for the req.body
	for _, device := range devices {
		if device.IsL3 {
			if removeCustom {
				removeCustomName(myhop, "", device.DefaultName, device.ID, makeChanges)
			} else {
				if device.DNSName == "" {
					if device.CustomName == "" {
						time.Sleep(time.Duration(timeInMs) * time.Millisecond)
						answers, err := net.LookupAddr(device.Ipaddr4)
						if err == nil {
							changeDeviceName(myhop, answers[0], device.DefaultName, device.ID, makeChanges, false)
						}else{
							fmt.Println("DNS Lookup failed for "+ device.Ipaddr4)
						}

					}else if(overWriteCustom){
						time.Sleep(500 * time.Millisecond)
						answers, err := net.LookupAddr(device.Ipaddr4)
						if err == nil {
							changeDeviceName(myhop, answers[0], device.DefaultName, device.ID, makeChanges, false)
						}else {
							fmt.Println("DNS Lookup failed for "+ device.Ipaddr4)
						}
					}else if(askForCustom){
						time.Sleep(500 * time.Millisecond)
						answers, err := net.LookupAddr(device.Ipaddr4)
						if err == nil {
							changeDeviceName(myhop, answers[0], device.CustomName, device.ID, makeChanges, true)
						}else {
							fmt.Println("DNS Lookup failed for "+ device.Ipaddr4)
						}
					}
				}
			}
		}
	}
}
