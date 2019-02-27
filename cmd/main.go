/* ====================================================
#   Copyright (C)2019 All rights reserved.
#
#   Author        : domchan
#   Email         : 814172254@qq.com
#   File Name     : main.go
#   Created       : 2019-01-29 11:53:54
#   Describe      :
#
# ====================================================*/
package main

import (
	"fmt"
	"github.com/peterh/liner"
	"os"
	"path"
	"regexp"

	"github.com/hiruok/etcd-cli/cmd/opts"

	"github.com/hiruok/etcd-cli/pkg/cmd"
	"github.com/hiruok/etcd-cli/pkg/version"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:          "etcd-cli",
		Long:         "",
		SilenceUsage: true,
		RunE:         rootE,
	}

	saveCmd = &cobra.Command{
		Use:          "upload",
		Short:        "Use to upload the configuration file in binary form to the legal path specified by etcd",
		RunE:         uploadE,
		SilenceUsage: true,
	}

	downloadCmd = &cobra.Command{
		Use:          "download",
		Short:        "Download the path-mapped file locally.",
		RunE:         downloadE,
		SilenceUsage: true,
	}

	etcdHost    string
	etcdPort    int32
	historyPath = path.Join(os.Getenv("HOME"), ".etcdcli_history")
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&etcdHost, "host", "s", "127.0.0.1", "Etcd connection host.")
	rootCmd.PersistentFlags().Int32VarP(&etcdPort, "port", "p", 2379, "Etcd connection port.")

	cmd.AddFlags(rootCmd)
	rootCmd.AddCommand(saveCmd)
	rootCmd.AddCommand(downloadCmd)
	rootCmd.AddCommand(version.Command())
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}

func rootE(_ *cobra.Command, _ []string) error {
	// set line
	line := liner.NewLiner()
	defer line.Close()
	line.SetCtrlCAborts(true)

	// load history from file
	loadHistory(line)
	defer saveHistory(line)

	r, err := opts.NewRoot(etcdHost, etcdPort)
	if err != nil {
		return err
	}
	defer r.Close()

	reg, _ := regexp.Compile(`'.*?'|".*?"|\S+`)
	for {
		prompt := fmt.Sprintf("%s[%d]  %s>", etcdHost, etcdPort, opts.PWD)

		c, err := line.Prompt(prompt)

		if err != nil {
			return err
		}
		line.AppendHistory(c)

		fields := reg.FindAllString(c, -1)

		err = r.DoScan(fields)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	return nil
}

func uploadE(_ *cobra.Command, args []string) error {
	r, err := opts.NewRoot(etcdHost, etcdPort)
	if err != nil {
		return err
	}
	defer r.Close()
	if len(args) != 2 {
		return opts.ErrInvalidParamNum
	}
	return r.Upload(args[0], args[1])
}

func downloadE(_ *cobra.Command, args []string) error {
	r, err := opts.NewRoot(etcdHost, etcdPort)
	if err != nil {
		return err
	}
	defer r.Close()
	if len(args) != 2 && len(args) != 1 {
		return opts.ErrInvalidParamNum
	}
	var localP = "./"
	if len(args) == 2 {
		localP = args[1]
	}
	return r.Download(args[0], localP)
}

func saveHistory(line *liner.State) {
	if f, err := os.OpenFile(historyPath, os.O_CREATE|os.O_RDWR, 0666); err != nil {
		fmt.Errorf("Error writing history file: %s", err.Error())
	} else {
		line.WriteHistory(f)
		f.Close()
	}
}

func loadHistory(line *liner.State) {
	if f, err := os.OpenFile(historyPath, os.O_RDONLY, 0444); err == nil {
		line.ReadHistory(f)
		f.Close()
	}
}
