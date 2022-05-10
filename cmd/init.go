/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/currycan/supkube/cert"
	"github.com/currycan/supkube/cni"
	"github.com/currycan/supkube/core"
	"github.com/currycan/supkube/pkg/logger"
)

var exampleInit = `
	# init with password with three master one node
	supkube init --passwd your-server-password  \
	--master 192.168.0.2 --master 192.168.0.3 --master 192.168.0.4 \
	--node 192.168.0.5 --user root \
	--version v1.18.0 --pkg-url=/root/kube1.18.0.tar.gz

	# init with pk-file , when your server have different password
	supkube init --pk /root/.ssh/id_rsa \
	--master 192.168.0.2 --node 192.168.0.5 --user root \
	--version v1.18.0 --pkg-url=/root/kube1.18.0.tar.gz

	# when use multi network. set a can-reach with --interface
 	supkube init --interface 192.168.0.254 \
	--master 192.168.0.2 --master 192.168.0.3 --master 192.168.0.4 \
	--node 192.168.0.5 --user root --passwd your-server-password \
	--version v1.18.0 --pkg-url=/root/kube1.18.0.tar.gz

	# when your interface is not "eth*|en*|em*" like.
	supkube init --interface your-interface-name \
	--master 192.168.0.2 --master 192.168.0.3 --master 192.168.0.4 \
	--node 192.168.0.5 --user root --passwd your-server-password \
	--version v1.18.0 --pkg-url=/root/kube1.18.0.tar.gz
`

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Simplest way to init your kubernets HA cluster",
	Long: `supkube init --master 192.168.0.2 --master 192.168.0.3 --master 192.168.0.4 \
	--node 192.168.0.5 --user root --passwd your-server-password \
	--version v1.18.0 --pkg-url=/root/kube1.18.0.tar.gz`,
	Example: exampleInit,
	Run: func(cmd *cobra.Command, args []string) {
		c := &core.Config{}
		// 没有重大错误可以直接保存配置. 但是apiservercertsans为空. 但是不影响用户 clean
		// 如果用户指定了配置文件,并不使用--master, 这里就不dump, 需要使用load获取配置文件了.
		if cfgFile != "" && len(core.MasterIPs) == 0 {
			err := c.Load(cfgFile)
			if err != nil {
				logger.Error("load cfgFile %s err: %q", cfgFile, err)
				os.Exit(1)
			}
		} else {
			c.Dump(cfgFile)
		}
		core.BuildInit()
		// 安装完成后生成完整版
		c.Dump(cfgFile)
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		// 使用了cfgFile 就不进行preRun了
		if cfgFile == "" && core.ExitInitCase() {
			_ = cmd.Help()
			os.Exit(core.ErrorExitOSCase)
		}
	},
}

func init() {
	initCmd.AddCommand(NewInitGenerateCmd())
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.
	initCmd.Flags().StringVar(&core.SSHConfig.User, "user", "root", "servers user name for ssh")
	initCmd.Flags().StringVar(&core.SSHConfig.Password, "passwd", "", "password for ssh")
	initCmd.Flags().StringVar(&core.SSHConfig.PkFile, "pk", cert.GetUserHomeDir()+"/.ssh/id_rsa", "private key for ssh")
	initCmd.Flags().StringVar(&core.SSHConfig.PkPassword, "pk-passwd", "", "private key password for ssh")

	initCmd.Flags().StringVar(&core.KubeadmFile, "kubeadm-config", "", "kubeadm-config.yaml template file")

	initCmd.Flags().StringVar(&core.APIServer, "apiserver", "apiserver.cluster.local", "apiserver domain name")
	initCmd.Flags().StringVar(&core.VIP, "vip", "10.103.97.2", "virtual ip")
	initCmd.Flags().StringSliceVar(&core.MasterIPs, "master", []string{}, "kubernetes multi-masters ex. 192.168.0.2-192.168.0.4")
	initCmd.Flags().StringSliceVar(&core.NodeIPs, "node", []string{}, "kubernetes multi-nodes ex. 192.168.0.5-192.168.0.5")
	initCmd.Flags().StringSliceVar(&core.CertSANS, "cert-sans", []string{}, "kubernetes apiServerCertSANs ex. 47.0.0.22 supkube.com ")

	initCmd.Flags().StringVar(&core.PkgURL, "pkg-url", "", "http://store.lameleg.com/kube1.14.1.tar.gz download offline package url, or file location ex. /root/kube1.14.1.tar.gz")
	initCmd.Flags().StringVar(&core.Version, "version", "", "version is kubernetes version")
	initCmd.Flags().StringVar(&core.Repo, "repo", "k8s.gcr.io", "choose a container registry to pull control plane images from")
	initCmd.Flags().StringVar(&core.PodCIDR, "podcidr", "172.16.0.0/16", "Specify range of IP addresses for the pod network")
	initCmd.Flags().StringVar(&core.SvcCIDR, "svccidr", "10.96.0.0/12", "Use alternative range of IP address for service VIPs")
	initCmd.Flags().StringVar(&core.Interface, "interface", "eth.*|en.*|em.*", "name of network interface, when use calico IP_AUTODETECTION_METHOD, set your ipv4 with can-reach=192.168.0.1")

	initCmd.Flags().BoolVar(&core.WithoutCNI, "without-cni", false, "If true we not core cni plugin")
	initCmd.Flags().StringVar(&core.Network, "network", cni.CALICO, "cni plugin, calico..")
	initCmd.Flags().BoolVar(&core.BGP, "bgp", false, "bgp mode enable, calico..")
	initCmd.Flags().StringVar(&core.MTU, "mtu", "1440", "mtu of the ipip mode , calico..")
	initCmd.Flags().StringVar(&core.LvscareImage.Image, "lvscare-image", "fanux/lvscare", "lvscare image name")
	initCmd.Flags().StringVar(&core.LvscareImage.Tag, "lvscare-tag", "latest", "lvscare image tag name")

	initCmd.Flags().IntVar(&core.Vlog, "vlog", 0, "kubeadm log level")

	// 不像用户暴露
	// initCmd.Flags().StringVar(&core.CertPath, "cert-path", cert.GetUserHomeDir() + "/.supkube/pki", "cert file path")
	// initCmd.Flags().StringVar(&core.CertEtcdPath, "cert-etcd-path", cert.GetUserHomeDir() + "/.supkube/pki/etcd", "etcd cert file path")
}

func NewInitGenerateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "gen",
		Short: "show default supkube init config",
		Run: func(cmd *cobra.Command, args []string) {
			c := &core.Config{}
			c.ShowDefaultConfig()
		},
	}
}
