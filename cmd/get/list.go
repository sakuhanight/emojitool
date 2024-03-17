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
package get

import (
	"emojitool/lib"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"os"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		host := viper.GetString("host")

		endpoint := fmt.Sprintf("https://%s/api/emojis", host)
		body := map[string]interface{}{}

		res, err := lib.Request(endpoint, body)
		defer res.Body.Close()

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			panic(fmt.Sprintf("%+v", err))
		}

		if cmd.Flags().Lookup("type").Value.String() == "json" {
			fmt.Printf("%s\n", resBody)
		} else if cmd.Flags().Lookup("type").Value.String() == "csv" {
			type emoji struct {
				Aliases                                 []string `json:"aliases"`
				Name                                    string   `json:"name"`
				Category                                string   `json:"category"`
				Url                                     string   `json:"url"`
				LocalOnly                               bool     `json:"localOnly"`
				IsSensitive                             bool     `json:"isSensitive"`
				RoleIdsThatCanBeUsedThisEmojiAsReaction []string `json:"roleIdsThatCanBeUsedThisEmojiAsReaction"`
			}
			var result struct {
				Emojis []emoji `json:"emojis"`
			}

			err = json.Unmarshal(resBody, &result)
			if err != nil {
				panic(fmt.Sprintf("%+v", err))
			}

			emojiFile, err := os.Create("emojis.csv")
			if err != nil {
				panic(fmt.Sprintf("emojiFile create failed:  %+v", err))
			}
			emojiWriter := csv.NewWriter(emojiFile)
			defer emojiWriter.Flush()

			aliasesFile, err := os.Create("aliases.csv")
			if err != nil {
				panic(fmt.Sprintf("aliasesFile create failed:  %+v", err))
			}
			aliasesWriter := csv.NewWriter(aliasesFile)
			defer aliasesWriter.Flush()

			rolesFile, err := os.Create("roles.csv")
			if err != nil {
				panic(fmt.Sprintf("rolesFile create failed:  %+v", err))
			}
			rolesWriter := csv.NewWriter(rolesFile)
			defer rolesWriter.Flush()

			// emoji.csv header
			err = emojiWriter.Write([]string{
				"Name",
				"Category",
				"Url",
				"LocalOnly",
				"IsSensitive",
			})
			if err != nil {
				panic(fmt.Sprintf("emoji.csv header write failed:  %+v", err))
			}

			// aliases.csv header
			err = aliasesWriter.Write([]string{
				"Name",
				"Alias",
			})
			if err != nil {
				panic(fmt.Sprintf("aliases.csv header write failed:  %+v", err))
			}

			// roles.csv header
			err = rolesWriter.Write([]string{
				"Name",
				"Role",
			})
			if err != nil {
				panic(fmt.Sprintf("roles.csv header write failed:  %+v", err))
			}

			// slice processing
			for _, emoji := range result.Emojis {
				emojiWriter.Write([]string{
					emoji.Name,
					emoji.Category,
					emoji.Url,
					fmt.Sprintf("%v", emoji.LocalOnly),
					fmt.Sprintf("%v", emoji.IsSensitive),
				})

				for _, alias := range emoji.Aliases {
					aliasesWriter.Write([]string{
						emoji.Name,
						alias,
					})
				}

				for _, role := range emoji.RoleIdsThatCanBeUsedThisEmojiAsReaction {
					rolesWriter.Write([]string{
						emoji.Name,
						role,
					})
				}
			}

		}
	},
}

func init() {
	Cmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")
	listCmd.Flags().String("type", "json", "output type(json/csv)")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
