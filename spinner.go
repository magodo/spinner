package spinner

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func Run(f func(ctx context.Context), msg string, sp spinner.Model) error {
	doneCh := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		f(ctx)
		doneCh <- struct{}{}
	}()
	s := model{m: sp, msg: msg, doneCh: doneCh, cancel: cancel}
	return tea.NewProgram(s).Start()
}

type model struct {
	m   spinner.Model
	msg string
	// doneCh is used to receive the done msg from the invoked function.
	doneCh <-chan struct{}
	// cancel is used to cancel the ctx of the invoked function, to notify the invoked function that it is interrupted
	cancel context.CancelFunc
	// canceled indicates that the app is canceled
	canceled bool
}

func (m model) Init() tea.Cmd {
	return m.m.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	select {
	case <-m.doneCh:
		return m, tea.Quit
	default:
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.cancel()
			m.canceled = true
			return m, nil
		default:
			return m, nil
		}
	default:
		var cmd tea.Cmd
		m.m, cmd = m.m.Update(msg)
		return m, cmd
	}
}

func (m model) View() string {
	if m.canceled {
		return fmt.Sprintf("%s Stopping on interrupt\n", m.m.View())
	}
	return fmt.Sprintf("%s %s\n", m.m.View(), m.msg)
}
