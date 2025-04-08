/*
Copyright © 2025 Brian Ketelsen <bketelsen@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/bketelsen/toolbox/cobra"
	"github.com/bketelsen/toolbox/cobra/doc"
	"github.com/spf13/viper"
)

// NewGendocsCommand creates a new command to generate documentation for the project
func NewGendocsCommand(config *viper.Viper) *cobra.Command {
	// gendocsCmd represents the gendocs command
	gendocsCmd := &cobra.Command{
		Use:    "gendocs",
		Hidden: true,
		Short:  "Generates documentation for the project",
		Long: `Generates documentation for the command using the cobra doc generator.
The documentation is generated in the ./content/docs/cli directory and
is in markdown format.`,
		Run: func(cmd *cobra.Command, args []string) {
			bp := config.GetString("docs.basepath")
			cmd.Logger.Info("Base path for documentation", "basepath", bp)
			linkHandler := func(name string) string {
				base := strings.TrimSuffix(name, path.Ext(name))
				return bp + "/docs/cli/" + strings.ToLower(base) + "/"
			}
			filePrepender := func(filename string) string {
				now := time.Now().Format(time.RFC3339)
				name := filepath.Base(filename)
				base := strings.TrimSuffix(name, path.Ext(name))
				url := "/docs/cli/" + strings.ToLower(base) + "/"
				return fmt.Sprintf(fmTemplate, now, strings.Replace(base, "_", " ", -1), base, url)
			}
			err := os.MkdirAll("./content/docs/cli/", 0755)
			if err != nil {
				log.Fatal(err)
			}
			err = doc.GenMarkdownTreeCustom(cmd.Root(), "./content/docs/cli/", filePrepender, linkHandler)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	//	gendocsCmd.Flags().StringP("basepath", "b", "inventory", "Base path for the documentation (default is /inventory)")

	gendocsCmd.PreRunE = func(cmd *cobra.Command, _ []string) error {
		_ = config.BindPFlag("docs.basepath", cmd.Flags().Lookup("basepath"))
		return nil

	}
	// Define cobra flags, the default value has the lowest (least significant) precedence
	gendocsCmd.Flags().StringP("basepath", "b", "inventory", "Base path for the documentation (default is /inventory)")

	return gendocsCmd
}

const fmTemplate = `---
date: %s
title: "%s"
slug: %s
url: %s
---
`
