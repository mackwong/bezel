/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package utils

import (
	"bufio"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	"strconv"
)

var machineNum, masterNum int
var roleMaster int

func ScanCmdline() string {
	input := bufio.NewScanner(os.Stdin)
	if input.Scan() {
		log.Debugf("Your input is: %s", input.Text())
	}
	return input.Text()
}

func ValidateScanValue(field string, value string) (err error) {
	switch field {
	case "Name", "HostName":
		if len(value) > 64 {
			return errors.New("name is too long")
		}
	case "Arranger":
		arrangerConst := map[string]bool{"k3s": true, "ke": true, "edgesite": true}
		if !arrangerConst[value] {
			return errors.New("Arranger MUST be one of k3s, k3, edgesite")
		}
	case "UpstreamDNS", "DockerRegistry", "K8sMasterIP", "IP", "GatewayIP", "Netmask":
		if !ValidateIP(value) {
			return errors.New("Not a valid IP address")
		}
	case "Role":
		if value != "master" && value != "worker" {
			return errors.New("The role must be master or worker")
		}
		if value == "master" {
			roleMaster++
			if roleMaster > masterNum {
				return errors.New("You have configured master role number more than one in the global config file.")
			}
		}
	case "MachineNum":
		if machineNum, err = strconv.Atoi(value); err != nil {
			return fmt.Errorf("%q not a number. \n", value)
		}
	case "MasterNum":
		if masterNum, err = strconv.Atoi(value); err != nil {
			return fmt.Errorf("%q not a number. \n", value)
		}
		if masterNum != 1 && masterNum != 3 {
			return fmt.Errorf("Master number must be 1 or 3.")
		}
		if masterNum > machineNum {
			return fmt.Errorf("Master number is more than the whole. Please ensure master number is less than or equal machine number.")
		}
	}
	return nil
}

func ValidateIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

func IsExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
