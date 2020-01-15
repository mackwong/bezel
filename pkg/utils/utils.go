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
	"regexp"
)

var machineNum, masterNum, roleMaster int64

func ScanCmdline() string {
	input := bufio.NewScanner(os.Stdin)
	if input.Scan() {
		log.Debugf("Your input is: %s", input.Text())
	}
	return input.Text()
}

func ValidateValue(field string, value interface{}) (err error) {
	switch field {
	case "Name", "HostName":
		if match, _ := regexp.MatchString("[A-Za-z][A-Za-z0-9_]*", value.(string)); !match {
			return errors.New("name and hostname shoule match [A-Za-z][A-Za-z0-9_]*")
		}
	case "Arranger":
		arrangerConst := map[string]bool{"k3s": true, "ke": true, "edgesite": true}
		if !arrangerConst[value.(string)] {
			return errors.New("Arranger MUST be one of k3s, k3, edgesite")
		}
	case "UpstreamDNS", "DockerRegistry", "K8sMasterIP", "IP", "GatewayIP", "Netmask":
		if !ValidateIP(value.(string)) {
			return errors.New("Not a valid IP address")
		}
	case "Role":
		if value != "master" && value.(string) != "worker" {
			return errors.New("The role must be master or worker")
		}
		if value.(string) == "master" {
			roleMaster++
			if roleMaster > masterNum {
				return errors.New("You have configured master role number more than one in the global config file.")
			}
		}
	case "MachineNum":
		machineNum = value.(int64)
		if machineNum > 1000 {
			return fmt.Errorf("%q can not large than 1000. \n", value)
		}
	case "MasterIP":
		if len(value.([]string)) > int(masterNum) {
			return fmt.Errorf("Master IP list are more than masterNum")
		}
	case "MasterNum":
		masterNum = value.(int64)
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

func IsInCIDR(tip, cidr string) (bool, error) {
	ips, err := GetAllIPS(cidr)
	if err != nil {
		return false, err
	}
	// remove network address and broadcast address?
	for _, p := range ips {
		if p == tip {
			return true, nil
		}
	}
	return false, nil
}

func GetAllIPS(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		log.Errorf("parse cidr err: %s", err)
		return nil, err
	}
	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	return ips, nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
