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
	"fmt"
	"net"
	"os"
	"strconv"
)

/*
func ScanConfigFields () {
        l := make(map[string]string)
        for field, value range field_list {
                ScanCmdline()
                ValidateScanValue(l, field, value)
                l[field] := value
        }
        return l
}

func GenerateConfigAndPreview (in map[string]string) {
        WriteEdgeConfigYaml(in, upath)
        PreviewEdgeConfigYaml(upath)
}
*/

var machine_num, master_num int
var role_master = 0

func ScanCmdline() string {
	//fmt.Println("scanning")
	input := bufio.NewScanner(os.Stdin)
	if input.Scan() {
		fmt.Println("Your input is: ", input.Text(), ".")
	}
	return input.Text()
}

func ValidateScanValue(field string, value string) bool {
	switch field {
	case "arranger":
		arranger_const := map[string]bool{"k3s": true, "ke": true, "edgesite": true}
		if arranger_const[value] == false {
			fmt.Println("Arranger MUST be one of k3s,k3,edgesite")
			return false
		} else {
			return true
		}
	case "upstreamDNS", "dockerRegistry", "k8sHAVip", "ip", "gatewayIP", "netmask":
		if !ValidateIP(value) {
			fmt.Println("Not a valid IP address")
			return false
		} else {
			return true
		}
	case "role":
		if (value != "master") && (value != "worker") {
			fmt.Println("The role must be master or worker")
			return false
		} else {
			if value == "master" {
				role_master++
				if role_master > master_num {
					fmt.Println("--------------------------------------------------------")
					fmt.Println("You have configured master role number more than one in the global config file.")
					fmt.Println("Please don`t add master role any more.")
					return false
				} else {
					return true
				}
			}
		}
	case "machine-num":
		var ch bool
		if machine_num, ch = IfNumeral(value); ch {
			return true
		} else {
			fmt.Printf("%q not a number. \n", value)
			return false
		}
	case "master-num":
		var st bool
		if master_num, st = IfNumeral(value); st {
			if (value != "1") && (value != "3") {
				fmt.Println("Master number must be 1 or 3.")
				return false
			} else {
				if master_num > machine_num {
					fmt.Println("Master number is more than the whole. Please ensure master number is less than or equal machine number.")
					return false
				} else {
					return true
				}
			}
		} else {
			fmt.Printf("%q not a number. \n", value)
			return false
		}
	default:
		return true
	}
	return true
}

func ValidateIP(ip string) bool {
	if net.ParseIP(ip) == nil {
		return false
	} else {
		return true
	}
}

func IfNumeral(s string) (int, bool) {
	if v, err := strconv.Atoi(s); err == nil {
		return v, true
	} else {
		return -1, false
	}
}

func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		//fmt.Println(err)
		return false
	}
	return true
}

func CreateFile(path string) {
	_, err := os.Create(path)
	if err != nil {
		fmt.Println("create file error,", err)
	}
}
