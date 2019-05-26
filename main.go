package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kyle11235/cli/config"
)

const (
	MainCommand = "cli"
)

func main() {
	// init
	config.Init()

	// sub commands
	subCommands := map[string]string{
		"help":    "show help",
		"version": "show version",
		"context": "show context",
	}

	// version flag
	versionCommand := flag.NewFlagSet("version", flag.ExitOnError)
	var versionScope string
	versionScopeDefault := "current" // is only set without -a
	versionScopeDesc := "show versions"
	versionCommand.StringVar(&versionScope, "scope", versionScopeDefault, versionScopeDesc)
	versionCommand.StringVar(&versionScope, "s", versionScopeDefault, versionScopeDesc)

	// context flag
	contextCommand := flag.NewFlagSet("context", flag.ExitOnError)
	var contextUse string
	contextUseDefault := ""
	contextUseDesc := "use context"
	contextCommand.StringVar(&contextUse, "use", contextUseDefault, contextUseDesc)
	contextCommand.StringVar(&contextUse, "u", contextUseDefault, contextUseDesc)

	// usage - just list subCommands
	if len(os.Args) < 2 {
		fmt.Printf("\nUsage:\n\n")
		// left justify 10 width
		fmt.Printf("%-10s%-10s\n\n", "", MainCommand+" <command> [arguments]")
		fmt.Printf("The commands are:\n\n")
		for k, v := range subCommands {
			fmt.Printf("%-10s%-10s%s\n", "", k, v)
		}
		fmt.Printf("\nUse %s help <command> for more information about a command.\n\n", MainCommand)
		return
	}

	// check sub command
	switch os.Args[1] {
	case "help":
		// list child command of help
		if len(os.Args) < 3 {
			options := make([]string, len(subCommands)*2)
			i := 0
			for k := range subCommands {
				options[i] = k
				options[i+1] = "|"
				i = i + 2
			}
			options = options[:len(options)-1]
			fmt.Printf("\n%s help %s\n\n", MainCommand, options)
			return
		}
		// result
		fmt.Println(subCommands[os.Args[2]])
	case "version":
		// check flags
		versionCommand.Parse(os.Args[2:])
		if versionCommand.Parsed() {
			switch versionScope {
			case "current":
				fmt.Printf("\n%s version=%s\n\n", MainCommand, config.GetCurrentVersion())
			case "latest":
				version, err := config.GetLatestVersion()
				if err != nil {
					fmt.Printf("\nerror=%s\n\n", err)
					return
				}
				fmt.Printf("\nLatest version=%s\n\n", version)
			default:
			}
		}
	case "context":
		// check flags
		contextCommand.Parse(os.Args[2:])
		if contextCommand.Parsed() {
			if contextUse == "" {
				config.PrintConfig()
			} else {
				config.UseContext(contextUse)
			}
		}
	default:
		flag.PrintDefaults()
	}

}
