package main

import (
	"fmt"
	"os"

	"github.com/skx/subcommands"
)

// Recovery is good
func recoverPanic() {
	if os.Getenv("DEBUG") != "" {
		return
	}

	if r := recover(); r != nil {
		fmt.Printf("recovered from panic while running %v\n%s\n", os.Args, r)
		fmt.Printf("To see the panic run 'export DEBUG=on' and repeat.\n")
	}
}

// Register the subcommands, and run the one the user chose.
func main() {

	//
	// Catch errors
	//
	defer recoverPanic()

	//
	// Register each of our subcommands.
	//
	subcommands.Register(&SSLExpiryCommand{})
	subcommands.Register(&calcCommand{})
	subcommands.Register(&chooseFileCommand{})
	subcommands.Register(&chooseSTDINCommand{})
	subcommands.Register(&chronicCommand{})
	subcommands.Register(&collapseCommand{})
	subcommands.Register(&commentsCommand{})
	subcommands.Register(&cppCommand{})
	subcommands.Register(&envTemplateCommand{})
	subcommands.Register(&execSTDINCommand{})
	subcommands.Register(&expectCommand{})
	subcommands.Register(&feedsCommand{})
	subcommands.Register(&findCommand{})
	subcommands.Register(&fingerdCommand{})
	subcommands.Register(&html2TextCommand{})
	subcommands.Register(&httpGetCommand{})
	subcommands.Register(&httpdCommand{})
	subcommands.Register(&ipsCommand{})
	subcommands.Register(&markdownTOCCommand{})
	subcommands.Register(&passwordCommand{})
	subcommands.Register(&rssCommand{})
	subcommands.Register(&runDirectoryCommand{})
	subcommands.Register(&splayCommand{})
	subcommands.Register(&timeoutCommand{})
	subcommands.Register(&todoCommand{})
	subcommands.Register(&treeCommand{})
	subcommands.Register(&urlsCommand{})
	subcommands.Register(&validateJSONCommand{})
	subcommands.Register(&validateXMLCommand{})
	subcommands.Register(&validateYAMLCommand{})
	subcommands.Register(&versionCommand{})
	subcommands.Register(&watchCommand{})
	subcommands.Register(&withLockCommand{})

	//
	// Execute the one the user chose.
	//
	os.Exit(subcommands.Execute())
}
