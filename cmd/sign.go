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
	"fmt"
	"github.com/spf13/viper"
	"github.com/xykong/ApkChannels/sign"

	"github.com/spf13/cobra"
)

// signCmd represents the sign command
var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("sign called")
		sign.Sign()
	},
}

func init() {
	rootCmd.AddCommand(signCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// signCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// signCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	signCmd.PersistentFlags().StringP("in", "i", "", "Input APK file to sign.")
	_ = viper.BindPFlag("in", signCmd.PersistentFlags().Lookup("in"))

	signCmd.PersistentFlags().StringP("out", "o", "", "File into which to output the signed APK.")
	_ = viper.BindPFlag("out", signCmd.PersistentFlags().Lookup("out"))

	signCmd.PersistentFlags().BoolP("v1-signing-enabled", "1", true,
		"Whether to enable signing using JAR signing scheme (aka v1 signing scheme)")
	_ = viper.BindPFlag("v1-signing-enabled", signCmd.PersistentFlags().Lookup("v1-signing-enabled"))

	signCmd.PersistentFlags().BoolP("v2-signing-enabled", "2", false,
		"Whether to enable signing using APK Signature Scheme v2 (aka v2 signing scheme)")
	_ = viper.BindPFlag("v2-signing-enabled", signCmd.PersistentFlags().Lookup("v2-signing-enabled"))

	signCmd.PersistentFlags().StringP("channel", "c", "", "Channel info")
	_ = viper.BindPFlag("channel", signCmd.PersistentFlags().Lookup("channel"))
}
