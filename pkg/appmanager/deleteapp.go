package appmanager

import (
	"fmt"
	"os"

	"github.com/currycan/supkube/core"
	"github.com/currycan/supkube/pkg/logger"
)

type DeleteFlags struct {
	Config     string
	PkgURL     string
	WorkDir    string
	CleanForce bool
}

func GetDeleteFlags(appURL string) *DeleteFlags {
	return &DeleteFlags{
		Config:     core.PackageConfig,
		WorkDir:    core.WorkDir,
		PkgURL:     appURL,
		CleanForce: core.CleanForce,
	}
}

func DeleteApp(flag *DeleteFlags, cfgFile string) error {
	//TODO
	c := &core.Config{}
	if err := c.Load(cfgFile); err != nil {
		logger.Error(err)
		c.ShowDefaultConfig()
		os.Exit(0)
	}
	pkgConfig, _ := LoadAppConfig(flag.PkgURL, flag.Config)
	pkgConfig.URL = flag.PkgURL
	pkgConfig.Name = nameFromURL(flag.PkgURL)
	pkgConfig.Workdir = flag.WorkDir
	pkgConfig.Workspace = fmt.Sprintf("%s/%s", flag.WorkDir, pkgConfig.Name)

	if !flag.CleanForce {
		prompt := fmt.Sprintf("delete command will del your coreed %s App , continue delete (y/n)?", pkgConfig.Name)
		result := core.Confirm(prompt)
		if !result {
			logger.Info("delete  %s App is skip, Exit", pkgConfig.Name)
			os.Exit(-1)
		}
	}

	everyNodesCmd, masterOnlyCmd := NewDeleteCommands(pkgConfig.Cmds)
	masterOnlyCmd.Run(*c, pkgConfig)
	everyNodesCmd.CleanUp(*c, pkgConfig)
	return nil
}

// return command run on every nodes and run only on master node
func NewDeleteCommands(cmds []Command) (Runner, Runner) {
	everyNodesCmd := &RunOnEveryNodes{}
	masterOnlyCmd := &RunOnMaster{}
	for _, c := range cmds {
		switch c.Name {
		case "REMOVE", "STOP":
			everyNodesCmd.Cmd = append(everyNodesCmd.Cmd, c)
		case "DELETE":
			masterOnlyCmd.Cmd = append(masterOnlyCmd.Cmd, c)
		default:
			// logger.Warn("Unknown command:%s,%s", c.Name, c.Cmd)
			// don't care other commands
		}
	}
	return everyNodesCmd, masterOnlyCmd
}
