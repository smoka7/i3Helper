package main

import (
	"fmt"
	"os"
	"time"

	i3 "go.i3wm.org/i3/v4"
)

//sizes of container for snap
var sizes = map[string]string{
	"left":   "50ppt 100ppt",
	"right":  "50ppt 100ppt",
	"top":    "100ppt 50ppt",
	"bottom": "100ppt 50ppt",
}

//start positions of container for snap direction
var positions = map[string]string{
	"left":   "0ppt 0ppt",
	"right":  "50ppt 0ppt",
	"top":    "0ppt 0ppt",
	"bottom": "0ppt 50ppt",
}

var helpMessage = `valid commands are:
focus next
focus prev
snap left
snap top
snap bottom
snap right
peek 500ms
#Valid time units are ns,us,ms,s,m,h.`

func last(index int, len int) int {
	if index == 0 {
		return len - 1
	}
	return index - 1
}

func next(index int, len int) int {
	if index == len-1 {
		return 0
	}
	return index + 1
}

func traverseNodes(parent *i3.Node) (nodes []*i3.Node) {
	if parent.Window != 0 {
		nodes = append(nodes, parent)
	}
	for _, node := range parent.Nodes {
		nodes = append(nodes, traverseNodes(node)...)
	}
	for _, node := range parent.FloatingNodes {
		nodes = append(nodes, traverseNodes(node)...)
	}
	return
}

func getWindowNodes(tree i3.Tree) []*i3.Node {
	ws := tree.Root.FindFocused(func(n *i3.Node) bool {
		return n.Type == i3.WorkspaceNode
	})
	return traverseNodes(ws)
}

func snap(tree i3.Tree, dir string) {
	focused := tree.Root.FindChild(func(m *i3.Node) bool {
		return m.Focused == true
	})
	if _, ok := sizes[dir]; !ok {
		fmt.Println("valid directions are top,bottom,left and right")
		os.Exit(1)
	}
	if !focused.IsFloating() {
		i3.RunCommand("floating toggle")
	}
	i3.RunCommand("resize set " + sizes[dir])
	i3.RunCommand("move position " + positions[dir])
}

func focus(tree i3.Tree, direction string) {
	focusedIndex, fullScreen := 0, 0
	windowNodes := getWindowNodes(tree)
	windowCount := len(windowNodes)
	for i, node := range windowNodes {
		if node.Focused {
			focusedIndex = i
			fullScreen = int(node.FullscreenMode)
		}
	}
	if fullScreen >= 1 {
		i3.RunCommand("fullscreen toggle")
	}
	if direction == "prev" {
		i3.RunCommand(fmt.Sprintf("[id=%d] focus", windowNodes[last(focusedIndex, windowCount)].Window))
	} else if direction == "next" {
		i3.RunCommand(fmt.Sprintf("[id=%d] focus", windowNodes[next(focusedIndex, windowCount)].Window))
	}
	fmt.Println("valid directions are next and prev")
	os.Exit(1)
}

func peek(tree i3.Tree, timeout string) {
	var focusedId int64
	focusedIndex := 0
	windowNodes := getWindowNodes(tree)
	windowCount := len(windowNodes)
	for i, node := range windowNodes {
		if node.Focused {
			focusedIndex = i
			focusedId = node.Window
		}
	}
	t, err := time.ParseDuration(timeout)
	if err != nil {
		fmt.Println(err)
		return
	}
	for nextId := 0; int64(nextId) != focusedId; time.Sleep(t) {
		focusedIndex = next(focusedIndex, windowCount)
		nextId = int(windowNodes[focusedIndex].Window)
		i3.RunCommand(fmt.Sprintf("[id=%d] focus", nextId))
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println(helpMessage)
		os.Exit(1)
	}
	command := os.Args[1]
	arg := os.Args[2]
	tree, _ := i3.GetTree()
	switch command {
	case "snap":
		snap(tree, arg)
	case "focus":
		focus(tree, arg)
	case "peek":
		peek(tree, arg)
	default:
		fmt.Println(helpMessage)
		os.Exit(1)
	}
}
