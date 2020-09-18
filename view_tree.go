package main

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// How to color the tree: normal or selected
var treeColorNormal = tcell.ColorWhite
var treeColorSelected = tcell.ColorGreen

func displayTree(app *tview.Application, images []imageJSONLine) {
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
		leaf := tview.NewTreeNode(fmt.Sprintf("%-40s %s %-30s %s", v.Tag, v.ID, v.CreatedAt, v.Size)).
			SetSelectable(true).
			SetColor(treeColorNormal).
			SetReference(fmt.Sprintf("%s:%s", v.Repository, v.Tag))
		currTree.AddChild(leaf)
		doneRepo[v.Repository] = true
	}
	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		mainText := node.GetReference().(string)
		if selectedImages[mainText] {
			delete(selectedImages, mainText)
			node.SetColor(treeColorNormal)
		} else {
			selectedImages[mainText] = true
			node.SetColor(treeColorSelected)
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
