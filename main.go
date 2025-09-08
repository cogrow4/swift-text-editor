package main

import (
	"flag"
	"log"
	"os"
)

// SWIFT - Streamlined Workflow, Increased Focus Typography
// A modern, intuitive CLI text editor

func main() {
	var filePath string
	flag.StringVar(&filePath, "f", "", "File to edit")
	flag.Parse()

	// If no file specified via flag, check first argument
	if filePath == "" && len(os.Args) > 1 {
		filePath = os.Args[1]
	}

	editor := NewTextEditor(filePath)
	if err := editor.Run(); err != nil {
		log.Fatal(err)
	}
}
