/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	gp "github.com/Never-M/MyGossip/pkg/gossiper"
	"github.com/spf13/cobra"
)

var name string
var ip string

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start a gossiper node",
	Long: `start command is used for starting a gossiper node`,
	Run: func(cmd *cobra.Command, args []string) {

		g := gp.NewGossiper(name, ip)
		logger := g.GetLogger()
		logger.Info("Start " + name + " on " + ip)
		g.Start()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().StringVarP(&name, "name", "n", "", "node name")
	startCmd.Flags().StringVarP(&ip, "ip", "i", "", "node ip")

	startCmd.MarkFlagRequired("name")
	startCmd.MarkFlagRequired("ip")



	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
