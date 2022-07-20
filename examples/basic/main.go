package main

import (
	"context"
	"log"
	"time"

	bspinner "github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/magodo/spinner"
)

func main() {
	s := bspinner.New()
	s.Spinner = bspinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	if err := spinner.Run(s, func(ctx context.Context, msg spinner.Messager) error {
		msg.SetStatus("Init")
		time.Sleep(1 * time.Second)
		msg.SetStatus("Running")
		time.Sleep(1 * time.Second)
		return nil
	}); err != nil {
		log.Fatal(err)
	}
}
