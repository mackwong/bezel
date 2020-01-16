package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.bj.sensetime.com/diamond/service-providers/bezel/pkg/model"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

func NewGenerateCmd() *cobra.Command {
	generateCmd := &cobra.Command{
		Use:           "gen",
		Short:         "generate bezel config file",
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath := filepath.Join(outPutDir, "demo.yaml")
			if err := generate(filePath); err != nil {
				return err
			}
			log.Infof("%s is generated successfully, Please modify it and use `%s create -c %s` to generate edge configs",
				filePath, os.Args[0], filePath)
			return nil
		},
	}
	generateCmd.Flags().StringVarP(&outPutDir, "output", "o", "./", "The dir for store configs")
	return generateCmd
}

func generate(filePath string) error {
	cfg := model.NewSampleBezelConfig()
	out, err := yaml.Marshal(cfg)
	if err != nil {
		log.Errorf("can not marshal, err: %s", err)
		return err
	}
	err = ioutil.WriteFile(filePath, out, 0644)
	if err != nil {
		log.Errorf("write file %s, err: %s", filePath, err)
	}
	return err
}
