package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.bj.sensetime.com/diamond/service-providers/bezel/pkg/model"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
)

var (
	sourceFile  string
	templateDir string
)

func NewParseCmd() *cobra.Command {
	parseCmd := &cobra.Command{
		Use:           "parse",
		Short:         "Load sub config infos, render templates to destination files",
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			sc, err := loadSubConfig(sourceFile)
			if err != nil {
				return err
			}
			err = parseAllTemplates(templateDir, outPutDir, sc)
			log.Infof("parse templates successfully")
			return err
		},
	}
	parseCmd.Flags().StringVarP(&sourceFile, "source", "s", "", "Source files.")
	parseCmd.Flags().StringVarP(&templateDir, "template-dir", "t", "", "Target template path.")
	parseCmd.Flags().StringVarP(&outPutDir, "output", "o", "./", "The dir for store configs")

	_ = parseCmd.MarkFlagRequired("source")
	_ = parseCmd.MarkFlagRequired("template")
	_ = parseCmd.MarkFlagRequired("output")

	return parseCmd
}

func loadSubConfig(configPath string) (*model.SubConfig, error) {
	config, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Errorf("load sub config %s err: %s", configPath, err)
		return nil, err
	}

	sc := &model.SubConfig{}
	if err = yaml.Unmarshal(config, sc); err != nil {
		log.Error("Unmarshal sub config info failed.\n")
		return nil, err
	}

	if len(sc.HaPeer) == 0 {
		sc.HaPeer = make([]model.Peer, 3)
	}
	return sc, nil
}

func parseAllTemplates(templateDir, outputDir string, sc *model.SubConfig) error {
	files := make([]string, 0)
	listFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		log.Debugf("read %s", path)
		files = append(files, path)
		return nil
	}
	err := filepath.Walk(templateDir, listFunc)

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		log.Errorf("parse files err: %s", err)
		return err
	}

	for _, t := range tmpl.Templates() {
		filePath := filepath.Join(outputDir, t.Name())
		f, err := os.Create(filePath)
		if err != nil {
			log.Errorf("create file %s err: %s", filePath, err)
			return err
		}
		defer f.Close()
		if err = t.Execute(f, sc); err != nil {
			log.Errorf("execute file %s err : %s", filePath, err)
			return err
		}
	}
	return nil
}
