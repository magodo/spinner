package spinner

import (
	"context"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

func ExampleSpinner() {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	Run(func(ctx context.Context) {
		time.Sleep(1 * time.Second)
		if ctx.Err() != nil {
			return
		}
		time.Sleep(3 * time.Second)
	}, "sleeping", s)
}
