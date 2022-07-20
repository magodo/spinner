package spinner

import (
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

// UserFunc executes with spinner running.
type UserFunc func(msg Messager) error

func Run(sp spinner.Model, f UserFunc, opts ...tea.ProgramOption) error {
	statusCh := make(chan string, 1)
	detailCh := make(chan string, 1)
	msg := Messager{
		statusCh: statusCh,
		detailCh: detailCh,
	}

	doneCh := make(chan any)
	var result error
	go func() {
		result = f(msg)
		close(statusCh)
		close(detailCh)
		close(doneCh)
	}()
	s := model{
		m:        sp,
		statusCh: statusCh,
		detailCh: detailCh,
		doneCh:   doneCh,
	}
	if err := tea.NewProgram(s, opts...).Start(); err != nil {
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

	interrupt bool

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
			m.interrupt = true
			return m, tea.Quit
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

	if m.interrupt {
		return fmt.Sprintf("%s %s\n", m.m.View(), "Stoppping on interrupt")
	}

	output := fmt.Sprintf("%s %s\n", m.m.View(), m.status)
	if m.detail != "" {
		output += "\n" + m.detail
	}
	return output
}
