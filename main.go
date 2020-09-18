package main

import (
	"flag"
	"fmt"

	"github.com/rivo/tview"
)

// Keep list of selected images (repository, tag) in a map, to make it easier
// to gauge whether an image is or isn't in the list.
var selectedImages = make(map[string]bool)

// What "FOO images" command to run (i.e. "docker images ..." or "podman ...")
var cmdType string = "docker"

// Whether the user wants to see a _list_ or a _tree_ of images. Tree is default.
var optList = false

func main() {
	flag.StringVar(&cmdType, "type", cmdType, "type of cmd to run (docker, podman) 'images' with")
	flag.BoolVar(&optList, "list", optList, "display as list instead of as tree")
	flag.Parse()
	app := tview.NewApplication()
	images := grabImages()

	if optList {
		displayList(app, images)
	} else {
		displayTree(app, images)
	}

	outputSelectedImages()
}

func outputSelectedImages() {
	if len(selectedImages) > 0 {
		fmt.Print(cmdType, " rmi")
		for img := range selectedImages {
			fmt.Print(" ", img)
		}
		fmt.Println("")
	}
}
