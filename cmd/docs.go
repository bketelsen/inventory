/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/spf13/cobra/doc"
)

// docsCmd represents the docs command
var docsCmd = &cobra.Command{
	Use:    "docs",
	Hidden: true,
	Short:  "Generate documentation for the project",
	Run: func(cmd *cobra.Command, args []string) {
		linkHandler := func(name string) string {
			base := strings.TrimSuffix(name, path.Ext(name))
			return "/docs/cli/" + strings.ToLower(base) + "/"
		}
		filePrepender := func(filename string) string {
			now := time.Now().Format(time.RFC3339)
			name := filepath.Base(filename)
			base := strings.TrimSuffix(name, path.Ext(name))
			url := "/docs/cli/" + strings.ToLower(base) + "/"
			return fmt.Sprintf(fmTemplate, now, strings.Replace(base, "_", " ", -1), base, url)
		}

		err := doc.GenMarkdownTreeCustom(rootCmd, "./docs/content/docs/cli/", filePrepender, linkHandler)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(docsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// docsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// docsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

const fmTemplate = `---
date: %s
title: "%s"
slug: %s
url: %s
---
`
