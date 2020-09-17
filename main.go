package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os/exec"
	"sort"
	"strings"

	//"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

var selectedImages = make(map[string]bool)
var cmdType string = "docker"

func addFuncItem(l *tview.List, mainText, secondaryText string, shortcut rune) {
	l.AddItem(mainText, secondaryText, shortcut, func() {
		if selectedImages[mainText] {
			delete(selectedImages, mainText)
		} else {
			selectedImages[mainText] = true
		}
		// refresh "view"
		for i := 0; i < l.GetItemCount()-1; i++ {
			main, secondary := l.GetItemText(i)
			if selectedImages[main] {
				secondary = "SELECTED"
			} else {
				secondary = ""
			}
			l.SetItemText(i, main, secondary)
		}
	})
}

type imageJSONLine struct {
	Tag        string
	Repository string
	ID         string
	Size       string
	CreatedAt  string
}

func addImages(l *tview.List) {
	cmd := exec.Command(cmdType, "images", "--format", `{"Tag":"{{.Tag}}","Repository":"{{.Repository}}","ID":"{{.ID}}","Size":"{{.Size}}","CreatedAt":"{{.CreatedAt}}"}`)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		panic(fmt.Sprintf("Executing '%s': %v", cmd, err))
	}
	images := make([]imageJSONLine, 0)
	// fmt.Printf("Executed: '%s'\n", cmd)
	// fmt.Printf("Got: '%v'\n", out.String())
	for _, line := range bytes.Split(out.Bytes(), []byte{'\n'}) {
		if len(line) > 0 {
			var v imageJSONLine
			if err := json.Unmarshal(line, &v); err != nil {
				panic(fmt.Sprintf("Unmarshalling '%v': %v", line, err))
			}
			if v.Repository != "<none>" && v.Tag != "<none>" && v.Tag != "latest" {
				if (cmdType == "docker" && strings.HasPrefix(v.Repository, "local/mfontani/") && !strings.HasPrefix(v.Repository, "local/mfontani/base/")) ||
					(cmdType == "podman" && strings.HasPrefix(v.Repository, "localhost/mfontani/") && !strings.HasPrefix(v.Repository, "localhost/mfontani/base/")) {
					images = append(images, v)
					// addFuncItem(l, fmt.Sprintf("%s:%s", v.Repository, v.Tag), "", 0)
					// fmt.Printf("Added: '%s:%s'\n", v.Repository, v.Tag)
				}
			}
		}
	}
	sort.SliceStable(images, func(i, j int) bool {
		if images[i].Repository == images[j].Repository {
			return images[i].Tag > images[j].Tag
		}
		return images[i].Repository < images[j].Repository
	})
	for _, v := range images {
		addFuncItem(l, fmt.Sprintf("%s:%s", v.Repository, v.Tag), "", 0)
	}
}

func main() {
	flag.StringVar(&cmdType, "type", cmdType, "type of cmd to run (docker, podman) 'images' with")
	flag.Parse()
	app := tview.NewApplication()
	list := tview.NewList()
	addImages(list)
	list.AddItem("Quit", "Press to exit", 'q', func() {
		app.Stop()
	})
	if err := app.SetRoot(list, true).SetFocus(list).Run(); err != nil {
		panic(err)
	}

	if len(selectedImages) > 0 {
		fmt.Print(cmdType, " rmi")
		for img := range selectedImages {
			fmt.Print(" ", img)
		}
		fmt.Println("")
	}
}
