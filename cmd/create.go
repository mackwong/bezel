package cmd

import (
	"fmt"
	"github.com/alecthomas/gometalinter/_linters/src/gopkg.in/yaml.v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.bj.sensetime.com/diamond/bezel/pkg/model"
	"gitlab.bj.sensetime.com/diamond/bezel/pkg/utils"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strconv"
)

var (
	GlobalConfigFile = "edge-config.yaml"
	GenerateSubConfig bool
)

func NewCreateCmd() *cobra.Command {
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create global config for edge cluster",
		RunE: func(cmd *cobra.Command, args []string) error{
			pwd, _ := os.Getwd()
			gPath := path.Join(pwd, GlobalConfigFile)
			input := ScanConfigFields()
			if err := writeEdgeConfigYaml(input, gPath); err != nil {
				return err
			}
			if GenerateSubConfig {
				if err := SplitFromGlobalConfig(gPath); err != nil {
					return nil
				}
			}
			log.Info("Sub config files will be at `sub` path.")
			return nil
		},
	}
	createCmd.Flags().BoolVarP(&GenerateSubConfig, "sub-config", "s", true, "If sub-config(s) flag used, it will generate both global and sub config files.")
	return createCmd
}

func ScanConfigFields() *model.GlobalConfig {

	var diamondFields = make(map[string]string)
	var machineFields = make(map[string]string)
	var machineConfigs = make(map[int]map[string]string)

	// TODO add default value on interactive cmdline
	diamondFields["name"] = "diamond-edge-ha"
	diamondFields["arranger"] = "edgesite"
	diamondFields["upstreamDNS"] = "114.114.114.114"
	diamondFields["dockerRegistry"] = "10.5.49.73"
	diamondFields["machine-num"] = "4"
	diamondFields["master-num"] = "3"
	diamondFields["k8sMasterIP"] = "10.4.72.231"

	// TODO add default value on interactive cmdline
	machineFields["name"] = "test-00"
	machineFields["hostname"] = "test-00"
	machineFields["ip"] = "10.4.72.140"
	machineFields["netmask"] = "255.255.255.0"
	machineFields["role"] = "master"
	machineFields["gatewayIP"] = "10.4.72.1"

	df := ScanInputToMapCache(diamondFields)

	machineNum, _ := strconv.Atoi(df["machine-num"])
	log.Println("=====================================")
	log.Printf("Your have %d machine to configure details. \n", machineNum)
	for i := 0; i < machineNum; i++ {
		log.Println("=====================================")
		log.Printf("Your are configuring the machine %d. \n", i)
		log.Println("=====================================")
		mf := ScanInputToMapCache(machineFields)
		machineConfigs[i] = mf
	}
	return model.NewGlobalConfig(df, machineConfigs)
}

func ScanInputToMapCache(fieldsMap map[string]string) map[string]string {
	var l = make(map[string]string)
	var listToSort = []string{}
	for f := range fieldsMap {
		listToSort = append(listToSort, f)
	}
	sort.Strings(listToSort)
	log.Println("All fields you should configure:", listToSort)
	for _, field := range listToSort {
		for {
			log.Println("Please configure", field)
			input := utils.ScanCmdline()
			if input == "" {
				log.Infof("No input on field %s, will use the default value. ", field)
				break
			}
			if utils.ValidateScanValue(field, input) {
				l[field] = input
				break
			}
			log.Warningf("invalidate value of field %s, please input the right value", field)
		}
	}
	return l
}

func writeEdgeConfigYaml(gc *model.GlobalConfig, path string) (err error) {
	var yamlByte []byte
	yamlByte, err = yaml.Marshal(gc)
	if err != nil {
		log.Errorf("can not marshal, err: %s", err)
		return
	}
	err = ioutil.WriteFile(path, yamlByte, 0644)
	if err != nil {
		log.Errorf("write edge config file error: %s", err)
	}
	return
}

func SplitFromGlobalConfig(cfgPath string) (err error){
	cfg, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		log.Errorf("read file %s err: %s", cfgPath, err)
		return
	}

	gc := &model.GlobalConfig{}
	err = yaml.Unmarshal(cfg, gc)
	if err != nil {
		log.Errorf("Unmarshal file %s err: %s", cfgPath, err)
		return
	}

	haPeer := make(map[string]string)
	for _, machine := range gc.Machines {
		if machine.Role == "master" {
			haPeer[machine.Hostname] = machine.IP
		}
	}

	sc := make(map[string]string)
	sc["arranger"] = gc.Diamond.Arranger
	sc["upstreamDNS"] = gc.Diamond.UpstreamDNS
	sc["k8sMasterIP"] = gc.Diamond.K8sMasterIP
	sc["dockerRegistry"] = gc.Diamond.DockerRegistry
	//sc["ha_peer"] = haPeer

	var haPeers []model.Peer
	if gc.Diamond.MasterNum == "1" {
		haPeers = []model.Peer{}
	} else {
		haPeers = model.NewHaPeer(haPeer)
	}

	for _, machine := range gc.Machines {
		sc["hostname"] = machine.Hostname
		sc["ip"] = machine.IP
		sc["netmask"] = machine.Netmask
		sc["gatewayIP"] = machine.GatewayIP
		sc["role"] = machine.Role
		nsc := model.NewSubConfig(sc, haPeers)
		var outFile string
		if machine.Role == "master" {
			outFile = fmt.Sprintf("sub-edge-config-master-%s.yaml", machine.IP)
		} else {
			outFile = fmt.Sprintf("sub-edge-config-worker-%s.yaml", machine.IP)
		}
		err = WriteSubConfigYaml(nsc, "./sub", outFile)
		if err != nil {
			return
		}
	}
	return
}

func WriteSubConfigYaml(config *model.SubConfig, parentDir string, fileName string) (err error) {
	yamlByte, err := yaml.Marshal(config)
	if err != nil {
		log.Errorf("can not marshal, err: %s", err)
		return
	}
	fullPath := path.Join(parentDir, fileName)
	log.Infof("Sub config will write to %s", fullPath)

	if !utils.IsExist(parentDir) {
		err = os.MkdirAll(parentDir, 0644)
		log.Errorf("can not mkdir dir %s, err: %s", parentDir, err)
		return
	}
	err = ioutil.WriteFile(fullPath, yamlByte, 0644)
	if err != nil {
		log.Errorf("can not read file %s, err: %s", fullPath, err)
	}
	return
}

