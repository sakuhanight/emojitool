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
	"bytes"
	"emojitool/misskey"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		filename := filepath.Base(path)
		token := viper.GetString("token")
		host := viper.GetString("host")
		endpoint := fmt.Sprintf("https://%s/api/drive/files/create", host)

		file, err := os.Open(path)
		if err != nil {
			zap.S().Panicf("Could not open '%s': %v", path, err)
		}

		body := &bytes.Buffer{}

		mw := multipart.NewWriter(body)

		{
			tokenPart, _ := mw.CreateFormField("i")
			_, err = tokenPart.Write([]byte(token))
			if err != nil {
				zap.S().Panicf("CreatePart Failed")
			}
		}
		{
			filePart, _ := mw.CreateFormFile("file", filename)
			_, err = io.Copy(filePart, file)
			if err != nil {
				zap.S().Panicf("CreatePart Failed")
			}
		}

		mw.Close()

		res, _ := misskey.RequestRaw(endpoint, mw.FormDataContentType(), body)
		defer res.Body.Close()

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			panic(fmt.Sprintf("%+v", err))
		}

		fmt.Printf("%s\n", resBody)

	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uploadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uploadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
