package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/tonyHuinker/ehop"
)

func askForInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(prompt)
	response, _ := reader.ReadString('\n')
	fmt.Println("\nThank You")
	return strings.TrimSpace(response)
}

func main() {
	//Specify Key File
	keyFile := askForInput("What is the name of your keyFile?")
	myhop := ehop.NewEDAfromKey(keyFile)

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
			if device.DNSName != "" {
				fmt.Println("Already has a name")
			} else {
				answers, err := net.LookupAddr(device.Ipaddr4)
				if err != nil {
					fmt.Println("DNS error")
				} else {
					for _, answer := range answers {
						fmt.Println(answer)
					}
				}
			}
		}
	}
}
