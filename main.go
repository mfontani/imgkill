package main

import (
	"bytes"
	"flag"
	"fmt"
	"os/exec"
	"strings"

	"github.com/rivo/tview"
)

// Keep list of selected images (repository, tag) in a map, to make it easier
// to gauge whether an image is or isn't in the list.
var selectedImages = make(map[string]bool)

// What "FOO images" command to run (i.e. "docker images ..." or "podman ...")
var cmdType string = "docker"

// Filter to only show images whose Repository contains this string
var onlyRepository string = ""

type arrayFlags []string

func (af *arrayFlags) String() string {
	return strings.Join(*af, ", ")
}

func (af *arrayFlags) Set(value string) error {
	*af = append(*af, value)
	return nil
}

// Whether the user wants to _remove_ any images Repository matching these
// strings from the list
var skipRepositories arrayFlags

// Whether the user wants to see a _list_ or a _tree_ of images. Tree is default.
var optList = false

// Whether the user wants to be prompted to run the deletion of images.
var promptForDeletion = false

func main() {
	flag.StringVar(&cmdType, "type", cmdType, "type of cmd to run (docker, podman) 'images' with")
	flag.BoolVar(&optList, "list", optList, "display as list instead of as tree")
	flag.StringVar(&onlyRepository, "only", onlyRepository, "only show Repository matching this")
	flag.Var(&skipRepositories, "skip", "skip showing repositories matching any")
	flag.BoolVar(&promptForDeletion, "prompt", promptForDeletion, "prompt for running the deletion commands")
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
		var cmdList []string
		cmdList = append(cmdList, "rmi")
		for img := range selectedImages {
			cmdList = append(cmdList, img)
		}
		fmt.Println(cmdType, strings.Join(cmdList, " "))
		if promptForDeletion {
			fmt.Println("Run above command? [N|y]")
			var ans string
			fmt.Scanln(&ans)
			if ans == "y" || ans == "Y" || ans == "yes" {
				cmd := exec.Command(cmdType, cmdList...)
				var out bytes.Buffer
				cmd.Stdout = &out
				err := cmd.Run()
				if err != nil {
					panic(fmt.Sprintf("Executing '%s': %v", cmd, err))
				}
				fmt.Println(out.String())
			}
		}
	}
}
