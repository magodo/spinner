package spinner

import (
	"context"
	"fmt"

	bspinner "github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type Spinner struct {
	ctx     context.Context
	spinner bspinner.Model
	msg     string
}

func New(ctx context.Context, msg string, spinner bspinner.Model) Spinner {
	return Spinner{ctx: ctx, spinner: spinner, msg: msg}
}

func (m Spinner) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m Spinner) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.ctx.Err() != nil {
		return m, tea.Quit
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		default:
			return m, nil
		}
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m Spinner) View() string {
	return fmt.Sprintf("%s %s\n", m.spinner.View(), m.msg)
}
