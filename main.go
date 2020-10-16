package main

import (
	"github.com/konoui/alfred-tldr/cmd"
)

func main() {
	rootCmd := cmd.NewRootCmd()
	cmd.Execute(rootCmd)
}
