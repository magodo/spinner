package main

import (
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

	if err := spinner.Run(s, func(msg spinner.Messager) error {
		jobs := []struct {
			name     string
			detail   string
			duration time.Duration
			err      error
		}{
			{
				name: "Eatting",
				detail: `- Milk
- Egg
- Corn`,
				duration: time.Second,
			},
			{
				name: "Babysitting",
				detail: `- Read a book
- Blow bubbles
- Watch TV`,
				duration: 3 * time.Second,
			},
			{
				name:     "Coding",
				detail:   "Write some Go code...",
				duration: time.Second,
			},
			{
				name:     "Playing",
				detail:   "Play some Dota",
				duration: 3 * time.Second,
				err:      errors.New("overheating"),
			},
		}

		for _, job := range jobs {
			msg.SetDetail(job.detail)
			msg.SetStatus(job.name)
			time.Sleep(job.duration)
			if job.err != nil {
				return job.err
			}
		}
		return nil
	}); err != nil {
		log.Fatal(err)
	}
}
