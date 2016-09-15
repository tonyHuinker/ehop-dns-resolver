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
	var slow = false
	var removeCustom = false
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

	ask3 := askForInput("Remove all existing custom names first?(yes/no)")
	if ask3 == "yes" {
		removeCustom = true
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
				if device.DNSName == "" && device.DhcpName == "" && device.CustomName == "" {
					if slow {
						time.Sleep(500 * time.Millisecond)
					}
					answers, err := net.LookupAddr(device.Ipaddr4)
					if err == nil {
						changeDeviceName(myhop, answers[0], device.DefaultName, device.ID, makeChanges)
					}
				}
			}
		}
	}
}
