package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"log"
	"strconv"
	"time"
)

func main() {
	kbtui, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer kbtui.Close()

	kbtui.SetManagerFunc(layout)
	if err := initKeybindings(kbtui); err != nil {
		log.Fatalln(err)
	}
	go testAsync(kbtui)
	if err := kbtui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := kbtui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func testAsync(kbtui *gocui.Gui) {
	for i := 0; i < 50; i++ {
		printToView(kbtui, "Chat", "Message #" + strconv.Itoa(i) + "\n")
		time.Sleep(1 * time.Second)
	}
	clearView(kbtui, "Chat")
}

func clearView(kbtui *gocui.Gui, viewName string) {
	kbtui.Update(func(g *gocui.Gui) error {
		inputView, err := kbtui.View(viewName)
		if err != nil {
			return err
		} else {
			inputView.Clear()
			inputView.SetCursor(0, 0)
			inputView.SetOrigin(0, 0)
		}
		return nil
	})
}

func printToView(kbtui *gocui.Gui, viewName string, message string) {
	kbtui.Update(func(g *gocui.Gui) error {
		updatingView, err := kbtui.View(viewName)
		if err != nil {
			return err
		} else {
			fmt.Fprintf(updatingView, message)
		}
		return nil
	})
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if feedView, err := g.SetView("Feed", 12, 0, maxX-1, maxY/5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		feedView.Autoscroll = true
		fmt.Fprintln(feedView, "Feed Window")
	}
	if chatView, err2 := g.SetView("Chat", 12, maxY/5+1, maxX-1, maxY-5); err2 != nil {
		if err2 != gocui.ErrUnknownView {
			return err2
		}
		chatView.Autoscroll = true
		fmt.Fprintln(chatView, "Chat Window")
	}
	if inputView, err3 := g.SetView("Input", 12, maxY-4, maxX-1, maxY-1); err3 != nil {
		if err3 != gocui.ErrUnknownView {
			return err3
		}
		 if _, err := g.SetCurrentView("Input"); err != nil {
		 	return err
		 }
		inputView.Editable = true
	}
	if listView, err4 := g.SetView("List", 0, 0, 10, maxY-1); err4 != nil {
		if err4 != gocui.ErrUnknownView {
			return err4
		}
		fmt.Fprintf(listView, "Lists\nWindow")
	}
	return nil
}
func getInputString(g *gocui.Gui) (string, error) {
	inputView, _ := g.View("Input")
	return inputView.Line(0)
}
func initKeybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return gocui.ErrQuit
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("Input", gocui.KeyEnter, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			handleInput(g)
			return nil
		}); err != nil {
		return err
	}
	return nil
}

func handleInput(g *gocui.Gui) {
	printToView(g, "Chat", "Enter was hit!\n")
	inputString, _ := getInputString(g)
	printToView(g, "Chat", inputString+"\n")
	clearView(g, "Input")
}
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
