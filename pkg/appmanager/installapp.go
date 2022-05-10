package appmanager

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/currycan/supkube/core"
	"github.com/currycan/supkube/pkg/logger"
)

type coreFlags struct {
	Envs    []string
	Config  string
	Values  string
	PkgURL  string
	WorkDir string
}

func GetcoreFlags(appURL string) *coreFlags {
	return &coreFlags{
		Config:  core.PackageConfig,
		PkgURL:  appURL,
		WorkDir: core.WorkDir,
		Values:  core.Values,
	}
}

func coreApp(flag *coreFlags, cfgFile string) error {
	c := &core.Config{}
	if err := c.Load(cfgFile); err != nil {
		logger.Error("%s", err)
		c.ShowDefaultConfig()
		os.Exit(0)
	}

	pkgConfig, err := LoadAppConfig(flag.PkgURL, flag.Config)
	if err != nil {
		logger.Error("Load App config from tarball err: ", err)
		os.Exit(0)
	}
	pkgConfig.URL = flag.PkgURL
	pkgConfig.Name = nameFromURL(flag.PkgURL)
	pkgConfig.Workdir = flag.WorkDir
	pkgConfig.Workspace = fmt.Sprintf("%s/%s", flag.WorkDir, pkgConfig.Name)
	s, err := getValuesContent(flag.Values)
	if err != nil {
		logger.Error("get values err:", err)
		os.Exit(-1)
	}
	pkgConfig.ValuesContent = s
	everyNodesCmd, masterOnlyCmd := NewcoreCommands(pkgConfig.Cmds)
	everyNodesCmd.Send(*c, pkgConfig)
	everyNodesCmd.Run(*c, pkgConfig)
	masterOnlyCmd.Send(*c, pkgConfig)
	masterOnlyCmd.Run(*c, pkgConfig)
	return nil
}

// return command run on every nodes and run only on master node
func NewcoreCommands(cmds []Command) (Runner, Runner) {
	everyNodesCmd := &RunOnEveryNodes{}
	masterOnlyCmd := &RunOnMaster{}
	for _, c := range cmds {
		switch c.Name {
		case "START", "LOAD":
			everyNodesCmd.Cmd = append(everyNodesCmd.Cmd, c)
		case "APPLY":
			masterOnlyCmd.Cmd = append(masterOnlyCmd.Cmd, c)
		default:
			// logger.Warn("Unknown command:%s,%s", c.Name, c.Cmd)
			// don't care other commands
		}
	}
	return everyNodesCmd, masterOnlyCmd
}

// getValuesContent is
func getValuesContent(s string) (valuesContent []byte, err error) {
	if s == "-" {
		// deal with stdin
		return ReadFromStdin()
	} else if s == "" {
		// use default and do nothing
		return nil, nil
	} else {
		// use -f file
		return ioutil.ReadFile(s)
	}
}

// ReadFromStdin is
func ReadFromStdin() (bt []byte, err error) {
	var b bytes.Buffer
	_, err = b.ReadFrom(os.Stdin)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
