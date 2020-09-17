package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os/exec"
	"sort"
	"strings"

	"github.com/gdamore/tcell"
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

func grabImages() []imageJSONLine {
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
	return images
}

var clrNormal = tcell.ColorWhite
var clrSelected = tcell.ColorGreen

var optList = false

func main() {
	flag.StringVar(&cmdType, "type", cmdType, "type of cmd to run (docker, podman) 'images' with")
	flag.BoolVar(&optList, "list", optList, "display as list instead of as tree")
	flag.Parse()
	app := tview.NewApplication()
	images := grabImages()

	if optList {
		list := tview.NewList()
		for _, v := range images {
			addFuncItem(list, fmt.Sprintf("%s:%s", v.Repository, v.Tag), "", 0)
		}
		list.AddItem("Quit", "Press to exit", 'q', func() {
			app.Stop()
		})
		list.SetDoneFunc(func() {
			app.Stop()
		})
		if err := app.SetRoot(list, true).SetFocus(list).Run(); err != nil {
			panic(err)
		}
	} else {
		root := tview.NewTreeNode(cmdType)
		tree := tview.NewTreeView().
			SetRoot(root).
			SetCurrentNode(root)

		doneRepo := make(map[string]bool)
		var currTree *tview.TreeNode
		for _, v := range images {
			if !doneRepo[v.Repository] {
				currTree = tview.NewTreeNode(v.Repository).
					SetSelectable(false)
				root.AddChild(currTree)
			}
			leaf := tview.NewTreeNode(fmt.Sprintf("%s:%-20s %s %s %s", v.Repository, v.Tag, v.ID, v.CreatedAt, v.Size)).
				SetSelectable(true).
				SetColor(clrNormal).
				SetReference(fmt.Sprintf("%s:%s", v.Repository, v.Tag))
			currTree.AddChild(leaf)
			doneRepo[v.Repository] = true
		}
		tree.SetSelectedFunc(func(node *tview.TreeNode) {
			mainText := node.GetReference().(string)
			if selectedImages[mainText] {
				delete(selectedImages, mainText)
				node.SetColor(clrNormal)
			} else {
				selectedImages[mainText] = true
				node.SetColor(clrSelected)
			}
		})
		tree.SetDoneFunc(func(k tcell.Key) {
			if k == tcell.KeyEscape {
				app.Stop()
			}
		})
		if err := app.SetRoot(tree, true).SetFocus(tree).Run(); err != nil {
			panic(err)
		}
	}

	if len(selectedImages) > 0 {
		fmt.Print(cmdType, " rmi")
		for img := range selectedImages {
			fmt.Print(" ", img)
		}
		fmt.Println("")
	}
}
