package main

import "github.com/konoui/tldr/cmd"

func main() {
	rootCmd := cmd.NewRootCmd()
	cmd.Execute(rootCmd)
}
