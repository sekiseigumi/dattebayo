package screens

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sekiseigumi/dattebayo/internal/logger"
)

type DashboardScreen struct {
	viewport viewport.Model
	messages []string
	ready    bool
}

type logMessage logger.LogEntry

func waitForActivity(logger *logger.Logger) tea.Cmd {
	return func() tea.Msg {
		for {
			select {
			case <-logger.Notify():
				entries := logger.Entries()
				if len(entries) > 0 {
					lastEntry := entries[len(entries)-1]
					return logMessage(lastEntry)
				}
			}
		}
	}
}

func dashboard() tea.Model {
	return DashboardScreen{
		messages: []string{},
	}
}

func (d DashboardScreen) Init() tea.Cmd {
	return tea.Batch(
		waitForActivity(globals.logger),
		tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return tickMsg{}
		}),
	)
}

func (d DashboardScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	if !d.ready {
		d.viewport = viewport.New(globals.width, globals.height-10)
		d.ready = true
	} else {
		d.viewport.Width, d.viewport.Height = globals.width, globals.height-10
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return d, tea.Quit
		}

		return d, nil

	case logMessage:
		d.messages = append(d.messages, msg.Message)
		d.viewport.SetContent(strings.Join(d.messages, "\n"))
		d.viewport.GotoBottom()

		if len(d.messages) > 100 {
			d.messages = d.messages[len(d.messages)-100:]
		}

		return d, nil
	}

	d.viewport, cmd = d.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return d, tea.Batch(cmds...)
}

func (d DashboardScreen) View() string {
	if !d.ready {
		return "Starting up..."
	}

	return fmt.Sprintf("%s\n\n%s", d.viewport.View(), "Press Ctrl+C to quit")
}
