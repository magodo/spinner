package spinner

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type Messager struct {
	detailCh chan string
	statusCh chan string
}

func (m *Messager) SetStatus(v string) {
	m.statusCh <- v
}
func (m *Messager) SetDetail(v string) {
	m.detailCh <- v
}

// UserFunc executes with spinner running. Users are expected to periodically check the context to see whether it is terminated, which might due to one of the following reasons:
// - Context deadline execeeded (if any)
// - User pressed "Ctrl-C" to interrupt
type UserFunc func(ctx context.Context, msg Messager) error

func Run(sp spinner.Model, f UserFunc) error {
	ctx, cancel := context.WithCancel(context.Background())

	statusCh := make(chan string, 1)
	detailCh := make(chan string, 1)
	msg := Messager{
		statusCh: statusCh,
		detailCh: detailCh,
	}

	doneCh := make(chan any)
	var result error
	go func() {
		result = f(ctx, msg)
		cancel()
		close(statusCh)
		close(detailCh)
		close(doneCh)
	}()
	s := model{
		m:        sp,
		statusCh: statusCh,
		detailCh: detailCh,
		ctx:      ctx,
		cancel:   cancel,
		doneCh:   doneCh,
	}
	if err := tea.NewProgram(s).Start(); err != nil {
		return err
	}
	return result
}

type model struct {
	m spinner.Model

	detail string
	status string

	statusCh <-chan string
	detailCh <-chan string

	ctx context.Context

	cancel   context.CancelFunc
	canceled bool

	doneCh <-chan any
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

	select {
	case v := <-m.statusCh:
		m.status = v
	default:
	}

	select {
	case v := <-m.detailCh:
		m.detail = v
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
	select {
	case <-m.doneCh:
		return ""
	default:
	}

	if m.canceled {
		return fmt.Sprintf("%s %s\n", m.m.View(), "Stoppping on interrupt")
	}

	output := fmt.Sprintf("%s %s\n", m.m.View(), m.status)
	if m.detail != "" {
		output += "\n" + m.detail
	}
	return output
}
