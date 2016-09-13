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

func changeDeviceName(myhop *ehop.EDA, newName string, oldName string, ID int, makeChange bool) {
	if makeChange {
		body := `{"custom_name": "` + newName + `"}`
		_, err := ehop.CreateEhopRequest("PATCH", "devices/"+strconv.Itoa(ID), body, myhop)
		if err != nil {
			fmt.Println("Error making device update call")
		} else {
			fmt.Println("Updated device succesfully from " + oldName + " to " + newName)
		}
	} else {
		fmt.Println("**Test** would have changed " + oldName + " to " + newName)
	}
}

func main() {
	var makeChanges = false
	var slow = false
	//Specify Key File
	keyFile := askForInput("What is the name of your keyFile?")
	myhop := ehop.NewEDAfromKey(keyFile)

	ask := askForInput("Actually make changes (yes) or just a test run (test)?")
	if ask == "yes" {
		makeChanges = true
	}

	ask2 := askForInput("Make DNS requests safely (enter>>slow) or as fast as possible(enter>>fast)?")
	if ask2 == "slow" {
		slow = true
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
			//fmt.Println("\nCustom Name: " + device.CustomName + "\nDefault Name: " + device.DefaultName + "\nDisplay Name: " + device.DisplayName + "\nDNS Name: " + device.DNSName + "\nDHCP Name: " + device.DhcpName)
			if !((device.DNSName != "") && (device.DhcpName != "") && (device.CustomName != "")) {
				answers, err := net.LookupAddr(device.Ipaddr4)
				if err == nil {
					if slow {
						time.Sleep(500 * time.Millisecond)
					}
					changeDeviceName(myhop, answers[0], device.DefaultName, device.ID, makeChanges)
				}
			}
		}
	}
}
