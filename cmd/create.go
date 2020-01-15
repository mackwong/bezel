package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.bj.sensetime.com/diamond/bezel/pkg/model"
	"gitlab.bj.sensetime.com/diamond/bezel/pkg/utils"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strconv"
)

var (
	GlobalConfigFileName = "edge-config.yaml"
	GenerateSubConfig    bool
	OutPutDir            = "./"
)

func NewCreateCmd() *cobra.Command {
	createCmd := &cobra.Command{
		Use:           "create",
		Short:         "Create global config form edge cluster",
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			input := ScanConfigFields()
			globalConfigFile := filepath.Join(OutPutDir, GlobalConfigFileName)
			if err := writeEdgeConfigYaml(input, globalConfigFile); err != nil {
				return err
			}
			if GenerateSubConfig {
				if err := SplitFromGlobalConfig(globalConfigFile, OutPutDir); err != nil {
					return nil
				}
			}
			return nil
		},
	}
	createCmd.Flags().BoolVarP(&GenerateSubConfig, "sub-config", "s", true, "If sub-config(s) flag used, it will generate both global and sub config files.")
	createCmd.Flags().StringVarP(&OutPutDir, "output", "o", "./", "The dir for store configs")
	return createCmd
}

func ScanConfigFields() *model.GlobalConfig {
	machineConfigs := make([]*model.MachineConfig, 0)
	defaultDiamondConfig := model.NewDefaultDiamondConfig()

	diamond := ScanInputToStruct(defaultDiamondConfig).(*model.DiamondConfig)
	machineNum, _ := strconv.Atoi(diamond.MachineNum)
	log.Infof("Your have %d machine to configure details. \n", machineNum)
	for i := 0; i < machineNum; i++ {
		log.Printf("\nYour are configuring the machine %d: \n", i)
		defaultMachineConfig := model.NewDefaultMachineConfig()
		machineConfig := ScanInputToStruct(defaultMachineConfig).(*model.MachineConfig)
		machineConfigs = append(machineConfigs, machineConfig)
	}
	return &model.GlobalConfig{
		Diamond:  diamond,
		Machines: machineConfigs,
	}
}

func ScanInputToStruct(obj interface{}) interface{} {
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
			if err := utils.ValidateScanValue(fName, input); err != nil {
				log.Errorf("invalidate value of field %s: %s please input the right value", fName, err)
				continue
			}
			vv.Field(i).SetString(input)
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

func SplitFromGlobalConfig(cfgPath, outputDir string) (err error) {
	var cfg []byte
	cfg, err = ioutil.ReadFile(cfgPath)
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
			haPeer[machine.HostName] = machine.IP
		}
	}

	var haPeers []model.Peer
	if gc.Diamond.MasterNum == "1" {
		haPeers = []model.Peer{}
	} else {
		haPeers = model.NewHaPeer(haPeer)
	}

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
		outFile := fmt.Sprintf("sub-edge-config-%s-%s.yaml", subConfig.Role, subConfig.IP)
		dir := filepath.Join(outputDir, "sub")
		if err = WriteSubConfigYaml(subConfig, dir, outFile); err != nil {
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
