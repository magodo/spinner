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
	if err := spinner.Run(s, func(ctx context.Context, logger spinner.Logger) error {
		logger.Print("Init")
		time.Sleep(1 * time.Second)

		// Press Ctrl-C within the first 1 second will be stopped here.
		if ctx.Err() != nil {
			return nil
		}

		logger.Print("Running")
		time.Sleep(1 * time.Second)
		logger.Print("Done")
		return nil
	}); err != nil {
		log.Fatal(err)
	}
}
