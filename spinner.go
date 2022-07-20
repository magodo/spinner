package spinner

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type Logger chan string

func (m Logger) Print(v ...any) {
	m <- fmt.Sprint(v...)
}
func (m Logger) Printf(format string, v ...any) {
	m <- fmt.Sprintf(format, v...)
}
func (m Logger) Println(v ...any) {
	m <- fmt.Sprintln(v...)
}

// UserFunc executes with spinner running. Each call on Print functions on the logger updates the message that the spinner shows.
// Users are expected to periodically check the context to see whether it is terminated, which might due to one of the following reasons:
// - Context deadline execeeded (if any)
// - User pressed "Ctrl-C" to interrupt
// - Calling Fatal-like functions on the logger
type UserFunc func(ctx context.Context, logger Logger) error

func Run(sp spinner.Model, f UserFunc) error {
	ctx, cancel := context.WithCancel(context.Background())
	logger := make(chan string)
	doneCh := make(chan any)
	var result error
	go func() {
		result = f(ctx, logger)
		cancel()
		close(logger)
		close(doneCh)
	}()
	s := model{m: sp, mch: logger, ctx: ctx, cancel: cancel, doneCh: doneCh}
	if err := tea.NewProgram(s).Start(); err != nil {
		return err
	}
	return result
}

type model struct {
	m   spinner.Model
	msg string
	mch <-chan string

	ctx context.Context

	// cancel is used to cancel the ctx of the invoked function, to notify the invoked function that it is interrupted
	cancel context.CancelFunc

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
	case message := <-m.mch:
		m.msg = message
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
		return fmt.Sprintf("%s %s\n", m.m.View(), "Stoppping on interrupt")
	}
	return fmt.Sprintf("%s %s\n", m.m.View(), m.msg)
}
