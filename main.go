package main

import (
	"fmt"
	"os"
	"time"

	i3 "go.i3wm.org/i3/v4"
)

//sizes of container for snapp
var sizes = map[string]string{
	"left":  "50ppt 100ppt",
	"right": "50ppt 100ppt",
	"top":   "100ppt 50ppt",
	"down":  "100ppt 50ppt",
}

//start positions of container for snapp direction
var positions = map[string]string{
	"left":  "0ppt 0ppt",
	"right": "50ppt 0ppt",
	"top":   "0ppt 0ppt",
	"down":  "0ppt 50ppt",
}

func last(index int, ids []int64) int {
	if index == 0 {
		return len(ids) - 1
	}
	return index - 1
}

func next(index int, ids []int64) int {
	if index == len(ids)-1 {
		return 0
	}
	return index + 1
}

func traverseNodes(parent *i3.Node) (ids []int64) {
	if parent.Window != 0 {
		ids = append(ids, parent.Window)
	}
	for _, node := range parent.Nodes {
		ids = append(ids, traverseNodes(node)...)
	}
	for _, node := range parent.FloatingNodes {
		ids = append(ids, traverseNodes(node)...)
	}
	return
}

func getWindowIds(tree i3.Tree) []int64 {
	ws := tree.Root.FindFocused(func(n *i3.Node) bool {
		return n.Type == i3.WorkspaceNode
	})
	return traverseNodes(ws)
}

func snapp(focused *i3.Node, dir string) {
	if !focused.IsFloating() {
		i3.RunCommand("floating toggle")
	}
	i3.RunCommand("resize set " + sizes[dir])
	i3.RunCommand("move position " + positions[dir])
}

func focus(tree i3.Tree, focused *i3.Node, direction string) {
	focusedIndex := 0
	windowIds := getWindowIds(tree)
	for i, id := range windowIds {
		if id == focused.Window {
			focusedIndex = i
		}
	}
	if focused.FullscreenMode == 1 {
		i3.RunCommand("fullscreen toggle")
	}
	if direction == "prev" {
		i3.RunCommand(fmt.Sprintf("[id=%d] focus", windowIds[last(focusedIndex, windowIds)]))
		return
	}
	i3.RunCommand(fmt.Sprintf("[id=%d] focus", windowIds[next(focusedIndex, windowIds)]))
}
func peek(tree i3.Tree, focused *i3.Node, timeout string) {
	focusedIndex := 0
	windowIds := getWindowIds(tree)
	for i, id := range windowIds {
		if id == focused.Window {
			focusedIndex = i
		}
	}
	t, err := time.ParseDuration(timeout)
	if err != nil {
		fmt.Println(err)
		return
	}
	for nextId := 0; int64(nextId) != focused.Window; time.Sleep(t) {
		focusedIndex = next(focusedIndex, windowIds)
		nextId = int(windowIds[focusedIndex])
		i3.RunCommand(fmt.Sprintf("[id=%d] focus", nextId))
	}
}

func main() {
	command := os.Args[1]
	arg := os.Args[2]
	tree, _ := i3.GetTree()
	focused := tree.Root.FindChild(func(m *i3.Node) bool {
		return m.Focused == true
	})
	switch command {
	case "snapp":
		snapp(focused, arg)
	case "focus":
		focus(tree, focused, arg)
	case "peek":
		peek(tree, focused, arg)
	}
}
