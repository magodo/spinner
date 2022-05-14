package spinner

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func ExampleSpinner() {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	p := tea.NewProgram(New(ctx, "Loading", s))
	if err := p.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
