package main

import (
	"context"
	"errors"
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
		if err := returnError(); err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
		logger.Print("Shouldn't see this")
		return nil
	}); err != nil {
		log.Fatal(err)
	}
}

func returnError() error {
	return errors.New("Some error")
}
