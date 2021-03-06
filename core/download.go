package core

import (
	"fmt"
	"net/url"
	"path"
	"sync"

	"github.com/currycan/supkube/pkg/logger"
	"github.com/currycan/supkube/pkg/sshcmd/cmd"
	"github.com/currycan/supkube/pkg/sshcmd/md5sum"
)

//location : url
//md5
//dst: /root
//hook: cd /root && rm -rf kube && tar zxvf %s  && cd /root/kube/shell && sh init.sh
func SendPackage(location string, hosts []string, dst string, before, after *string) string {
	var md5 string
	location, md5 = downloadFile(location)
	pkg := path.Base(location)
	fullPath := fmt.Sprintf("%s/%s", dst, pkg)
	mkDstDir := fmt.Sprintf("mkdir -p %s || true", dst)
	var wm sync.WaitGroup
	for _, host := range hosts {
		wm.Add(1)
		go func(host string) {
			defer wm.Done()
			_ = SSHConfig.CmdAsync(host, mkDstDir)
			logger.Debug("[%s]please wait for mkDstDir", host)
			if before != nil {
				logger.Debug("[%s]please wait for before hook", host)
				_ = SSHConfig.CmdAsync(host, *before)
			}
			if SSHConfig.IsFileExist(host, fullPath) {
				if SSHConfig.ValidateMd5sumLocalWithRemote(host, location, fullPath) {
					logger.Info("[%s]SendPackage:  %s file is exist and ValidateMd5 success", host, fullPath)
				} else {
					rm := fmt.Sprintf("rm -f %s", fullPath)
					_ = SSHConfig.Cmd(host, rm)
					// del then copy
					if ok := SSHConfig.CopyForMD5(host, location, fullPath, md5); ok {
						logger.Info("[%s]copy file md5 validate success", host)
					} else {
						logger.Error("[%s]copy file md5 validate failed", host)
					}
				}
			} else {
				if ok := SSHConfig.CopyForMD5(host, location, fullPath, md5); ok {
					logger.Info("[%s]copy file md5 validate success", host)
				} else {
					logger.Error("[%s]copy file md5 validate failed", host)
				}
			}
			if after != nil {
				logger.Debug("[%s]please wait for after hook", host)
				_ = SSHConfig.CmdAsync(host, *after)
			}
		}(host)
	}
	wm.Wait()
	return location
}

func DownloadFile(location string) (filePATH, md5 string) {
	return downloadFile(location)
}

//
func downloadFile(location string) (filePATH, md5 string) {
	if _, isURL := isURL(location); isURL {
		absPATH := "/tmp/supkube/" + path.Base(location)
		if !cmd.IsFileExist(absPATH) {
			//generator download cmd
			dwnCmd := downloadCmd(location)
			//os exec download command
			cmd.Cmd("/bin/sh", "-c", "mkdir -p /tmp/supkube && cd /tmp/supkube && "+dwnCmd)
		}
		location = absPATH
	}
	//file md5
	md5 = md5sum.FromLocal(location)
	return location, md5
}

//??????url ??????command
func downloadCmd(url string) string {
	//only http
	u, isHTTP := isURL(url)
	var c = ""
	if isHTTP {
		param := ""
		if u.Scheme == "https" {
			param = "--no-check-certificate"
		}
		c = fmt.Sprintf(" wget -c %s %s", param, url)
	}
	return c
}

func isURL(u string) (url.URL, bool) {
	if uu, err := url.Parse(u); err == nil && uu != nil && uu.Host != "" {
		return *uu, true
	}
	return url.URL{}, false
}
