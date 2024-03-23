/*
Copyright © 2024 Tokimine Sakuha <sakuha@tsuitachi.net>

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
	"emojitool/lib"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"io"
	"strings"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [name]",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			zap.S().Fatal("name is required")
		}

		var fileId string
		token := viper.GetString("token")
		host := viper.GetString("host")

		if cmd.Flags().Changed("file") {

			res, _ := lib.UploadExample(host, token, cmd.Flags().Lookup("file").Value.String())
			defer res.Body.Close()

			decoder := json.NewDecoder(res.Body)
			var data map[string]interface{}
			err := decoder.Decode(&data)
			if err != nil {
				panic(err)
			}
			fileId = data["id"].(string)
		}
		zap.S().Debugf("fileId: %s", fileId)

		// RequiredなParameterだけで初期化
		body := map[string]interface{}{
			"name":   args[0],
			"fileId": fileId,
			"i":      token,
		}

		if cmd.Flags().Changed("category") {
			body["category"], _ = cmd.Flags().GetString("category")
		}
		if cmd.Flags().Changed("license") {
			body["license"] = cmd.Flags().Lookup("license").Value.String()
		}
		if cmd.Flags().Changed("alias") {
			if cmd.Flags().Changed("delimiter") {
				var aliases []string
				delimiter, err := cmd.Flags().GetString("delimiter")
				if err != nil {
					zap.S().Fatal(err)
				}
				aliasStrings, err := cmd.Flags().GetStringArray("alias")
				if err != nil {
					zap.S().Fatal(err)
				}
				for _, aliasString := range aliasStrings {
					alias := strings.Split(aliasString, delimiter)
					aliases = append(aliases, alias...)
				}
				body["aliases"] = aliases
			} else {
				aliases, err := cmd.Flags().GetStringArray("alias")
				if err != nil {
					zap.S().Fatal(err)
				}
				body["aliases"] = aliases
			}
		}

		res, err := lib.Request(fmt.Sprintf("https://%s/api/admin/emoji/add", viper.GetString("host")), body)
		if err != nil {
			zap.S().Fatal(err)
		}
		defer res.Body.Close()
		zap.S().Debugf("%+v", res)
		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			zap.S().Debugf("failed to read response body: %+v", err)
		}
		zap.S().Debugf("%s", resBody)
		resMap := map[string]interface{}{}
		err = json.Unmarshal(resBody, &resMap)
		if err != nil {
			zap.S().Debugf("failed to decode response body: %+v", err)
		}
		if res.StatusCode == 200 {
			zap.S().Info("success")
		} else {
			zap.S().Warnf("emoji add failed: %s", resMap["error"].(map[string]interface{})["message"])
		}

	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")
	addCmd.Flags().String("file", "", "emoji local file")
	//addCmd.Flags().String("url", "", "emoji url")
	addCmd.Flags().String("category", "", "emoji category")
	addCmd.Flags().StringArray("alias", []string{}, "emoji alias (multiple)")
	addCmd.Flags().String("delimiter", "", "emoji alias delimiter")
	addCmd.Flags().String("license", "", "emoji license")
	addCmd.Flags().Bool("sensitive", false, "emoji is sensitive")
	addCmd.Flags().Bool("local-Only", false, "emoji is local only")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
