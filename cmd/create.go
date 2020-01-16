package cmd

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.bj.sensetime.com/diamond/service-providers/bezel/pkg/model"
	"gitlab.bj.sensetime.com/diamond/service-providers/bezel/pkg/utils"
	"gopkg.in/yaml.v2"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

const GlobalConfigFileName = "edge-config.yaml"

var (
	generateSubConfig bool
	outPutDir         string
	bezelConfig       string
)

func NewCreateCmd() *cobra.Command {
	createCmd := &cobra.Command{
		Use:           "create",
		Short:         "Create global config form edge cluster",
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			var gc *model.GlobalConfig
			var err error
			if bezelConfig != "" {
				if gc, err = GetGlobalConfigByConfig(bezelConfig); err != nil {
					return err
				}
			} else {
				if gc, err = ScanConfigFields(); err != nil {
					return err
				}
			}
			globalConfigFile := filepath.Join(outPutDir, GlobalConfigFileName)
			if err := writeEdgeConfigYaml(gc, globalConfigFile); err != nil {
				return err
			}
			log.Infof("%s generated successfully\n", globalConfigFile)
			if generateSubConfig {
				if err := SplitFromGlobalConfig(gc, outPutDir); err != nil {
					return nil
				}
				log.Infof("sub files %s generated successfully\n", outPutDir)
			}
			return nil
		},
	}
	createCmd.Flags().BoolVarP(&generateSubConfig, "sub-config", "s", true, "If sub-config(s) flag used, it will generate both global and sub config files.")
	createCmd.Flags().StringVarP(&outPutDir, "output", "o", "./", "The dir for store configs")
	createCmd.Flags().StringVarP(&bezelConfig, "config", "c", "", "bezel config file")
	return createCmd
}

func GetGlobalConfigByConfig(configFile string) (*model.GlobalConfig, error) {
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Errorf("read %s err: %s", configFile, err)
		return nil, err
	}

	config := &model.BezelConfig{}
	if err = yaml.Unmarshal(content, config); err != nil {
		log.Errorf("unmarshal err: %s", err)
		return nil, err
	}
	if err = validateBezelConfig(config); err != nil {
		return nil, err
	}

	machines := make([]*model.MachineConfig, 0)
	diamond := generateDiamondConfig(config)

	masters, err := generateMasterMachine(config)
	if err != nil {
		return nil, err
	}

	worker, err := generateWorkerMachine(masters, config)
	if err != nil {
		return nil, err
	}
	machines = append(machines, masters...)
	machines = append(machines, worker...)

	return &model.GlobalConfig{
		Diamond:  diamond,
		Machines: machines,
	}, nil
}

func generateMasterMachine(config *model.BezelConfig) ([]*model.MachineConfig, error) {
	var masterIndex int
	role := "master"
	out := make([]*model.MachineConfig, 0)
	for _, ip := range config.MasterIP {
		for _, ipr := range config.IPRange {
			isIn, err := utils.IsInCIDR(ip, ipr.IPRange)
			if err != nil {
				return nil, err
			}
			if isIn {
				name, hostname, err := generateNameAndHostFormat(config.NameFormat, config.HostNameFormat, role, ip, masterIndex)
				if err != nil {
					return nil, err
				}
				out = append(out, &model.MachineConfig{
					Name:      name,
					HostName:  hostname,
					Role:      role,
					IP:        ip,
					GatewayIP: ipr.GatewayIP,
					Netmask:   ipr.Netmask,
				})
				masterIndex++
			}

		}
	}
	num := int(config.MasterNum) - len(config.MasterIP)
	for i := 0; i < num; i++ {
	loop:
		for _, ipr := range config.IPRange {
			ips, err := utils.GetAllIPS(ipr.IPRange)
			if err != nil {
				return nil, err
			}
		findIP:
			for _, ip := range ips {
				if strings.HasSuffix(ip, ".0") || strings.HasSuffix(ip, ".255") {
					continue
				}
				for _, o := range out {
					if o.IP == ip {
						continue findIP
					}
				}
				name, hostname, err := generateNameAndHostFormat(config.NameFormat, config.HostNameFormat, role, ip, masterIndex)
				if err != nil {
					return nil, err
				}
				out = append(out, &model.MachineConfig{
					Name:      name,
					HostName:  hostname,
					Role:      role,
					IP:        ip,
					GatewayIP: ipr.GatewayIP,
					Netmask:   ipr.Netmask,
				})
				masterIndex++
				break loop
			}
		}
	}
	return out, nil
}

func generateNameAndHostFormat(namePattern, hostNamePattern, role, ip string, index int) (string, string, error) {
	tmplName, err := template.New("name").Parse(namePattern)
	if err != nil {
		log.Errorf("parse name err: %s", err)
		return "", "", err
	}
	tmplHost, err := template.New("hostname").Parse(hostNamePattern)
	if err != nil {
		log.Errorf("parse hostname err: %s", err)
		return "", "", err
	}
	f := model.Formatter{
		Role:  role,
		IP:    ip,
		Index: index,
	}
	var name, hostname bytes.Buffer
	err = tmplName.Execute(&name, f)
	if err != nil {
		log.Errorf("parse name err: %s", err)
		return "", "", err
	}
	err = tmplHost.Execute(&hostname, f)
	if err != nil {
		log.Errorf("parse name err: %s", err)
		return "", "", err
	}
	return name.String(), hostname.String(), nil
}

func generateWorkerMachine(masters []*model.MachineConfig, config *model.BezelConfig) ([]*model.MachineConfig, error) {
	var workerIndex int
	role := "worker"
	out := make([]*model.MachineConfig, 0)
	for i := 0; i < int(config.MachineNum-config.MasterNum); i++ {
	loop:
		for _, ipr := range config.IPRange {
			ips, err := utils.GetAllIPS(ipr.IPRange)
			if err != nil {
				return nil, err
			}
		findIP:
			for _, ip := range ips {
				if strings.HasSuffix(ip, ".0") || strings.HasSuffix(ip, ".255") {
					continue
				}
				for _, m := range masters {
					if m.IP == ip {
						continue findIP
					}
				}
				for _, o := range out {
					if o.IP == ip {
						continue findIP
					}
				}
				name, hostname, err := generateNameAndHostFormat(config.NameFormat, config.HostNameFormat, role, ip, workerIndex)
				if err != nil {
					return nil, err
				}
				out = append(out, &model.MachineConfig{
					Name:      name,
					HostName:  hostname,
					Role:      role,
					IP:        ip,
					GatewayIP: ipr.GatewayIP,
					Netmask:   ipr.Netmask,
				})
				workerIndex++
				break loop
			}
		}
	}
	return out, nil
}

func generateDiamondConfig(config *model.BezelConfig) *model.DiamondConfig {
	return &model.DiamondConfig{
		Name:           config.Name,
		Arranger:       config.Arranger,
		UpstreamDNS:    config.UpstreamDNS,
		DockerRegistry: config.DockerRegistry,
		MachineNum:     config.MachineNum,
		MasterNum:      config.MachineNum,
		K8sMasterIP:    config.K8sMasterIP,
	}
}

func validateBezelConfig(config *model.BezelConfig) error {
	v := reflect.ValueOf(*config)
	t := reflect.TypeOf(*config)
	for i := 0; i < v.NumField(); i++ {
		fName := t.Field(i).Name
		log.Debug(fName)
		if err := utils.ValidateValue(fName, v.Field(i).Interface()); err != nil {
			log.Errorf("invalidate value of field %s: %s please input the right value", fName, err)
			return err
		}
	}
	return nil
}

func ScanConfigFields() (*model.GlobalConfig, error) {
	machines := make([]*model.MachineConfig, 0)
	defaultDiamondConfig := model.NewDefaultDiamondConfig()

	diamond := ScanInputToStruct(defaultDiamondConfig).(*model.DiamondConfig)
	log.Infof("Your have %d machine to configure details. \n", diamond.MachineNum)
	for i := 0; i < int(diamond.MachineNum); i++ {
		log.Printf("\nYour are configuring the machine %d: \n", i)
		defaultMachineConfig := model.NewDefaultMachineConfig()
		machine := ScanInputToStruct(defaultMachineConfig).(*model.MachineConfig)
		machines = append(machines, machine)
	}
	return &model.GlobalConfig{
		Diamond:  diamond,
		Machines: machines,
	}, nil
}

func ScanInputToStruct(obj interface{}) interface{} {
	var value interface{}
	var err error
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr {
		log.Errorf("only support point")
		return nil
	}
	vv := v.Elem()
	t := vv.Type()

	if vv.Kind() != reflect.Struct {
		log.Errorf("cat not scan input to nonStruct")
		return nil
	}
	for i := 0; i < vv.NumField(); i++ {
		fName := t.Field(i).Name
		for {
			log.Infof("Please configure %s:\n", fName)
			input := utils.ScanCmdline()
			if input == "" {
				log.Infof("No input on field %s, will use the default value. ", fName)
				break
			}
			if fName == "MachineNum" || fName == "MasterNum" {
				value, err = strconv.ParseInt(input, 0, 64)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				value = input
			}
			if err := utils.ValidateValue(fName, value); err != nil {
				log.Errorf("invalidate value of field %s: %s please input the right value", fName, err)
				continue
			}

			if fName == "MachineNum" || fName == "MasterNum" {
				vv.Field(i).SetInt(value.(int64))
			} else {
				vv.Field(i).SetString(value.(string))
			}
			break
		}
	}
	return obj
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

func SplitFromGlobalConfig(gc *model.GlobalConfig, outDir string) (err error) {
	subConfigs := getSubConfigs(gc)
	return writeSubConfigYaml(subConfigs, outDir)
}

func getSubConfigs(gc *model.GlobalConfig) []*model.SubConfig{
	haPeer := make(map[string]string)
	for _, machine := range gc.Machines {
		if machine.Role == "master" {
			haPeer[machine.HostName] = machine.IP
		}
	}

	var haPeers []model.Peer
	if gc.Diamond.MasterNum == 1 {
		haPeers = []model.Peer{}
	} else {
		haPeers = model.NewHaPeer(haPeer)
	}

	subConfigs := make([]*model.SubConfig, 0)
	for _, machine := range gc.Machines {
		subConfig := &model.SubConfig{
			Arranger:       gc.Diamond.Arranger,
			UpstreamDNS:    gc.Diamond.UpstreamDNS,
			K8sMasterIP:    gc.Diamond.K8sMasterIP,
			DockerRegistry: gc.Diamond.DockerRegistry,
			Hostname:       machine.HostName,
			IP:             machine.IP,
			Netmask:        machine.Netmask,
			GatewayIP:      machine.GatewayIP,
			Role:           machine.Role,
			HaPeer:         haPeers,
		}
		subConfigs = append(subConfigs, subConfig)
	}
	return subConfigs
}

func writeSubConfigYaml(configs []*model.SubConfig, outDir string) error {
	for _, config := range configs {
		outFile := fmt.Sprintf("sub-edge-config-%s-%s.yaml", config.Role, config.IP)
		dir := filepath.Join(outDir, "sub")
		if err := doWriteSubConfigYaml(config, dir, outFile); err != nil {
			return err
		}
	}
	return nil
}

func doWriteSubConfigYaml(config *model.SubConfig, parentDir string, fileName string) (err error) {
	yamlByte, err := yaml.Marshal(config)
	if err != nil {
		log.Errorf("can not marshal, err: %s", err)
		return
	}
	fullPath := path.Join(parentDir, fileName)
	log.Infof("Sub config will write to %s", fullPath)

	if exist, err := utils.IsExist(parentDir); err != nil {
		log.Errorf("%s", err)
		return err
	} else if !exist {
		if err = os.MkdirAll(parentDir, 0755); err != nil {
			log.Errorf("can not mkdir dir %s, err: %s", parentDir, err)
			return err
		}
	}
	err = ioutil.WriteFile(fullPath, yamlByte, 0644)
	if err != nil {
		log.Errorf("can not read file %s, err: %s", fullPath, err)
	}
	return
}
