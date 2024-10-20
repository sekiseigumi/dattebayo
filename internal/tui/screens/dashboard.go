package screens

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type DashboardScreen struct {
	viewport viewport.Model
	ready    bool
}

type dashboardTickMsg time.Time

func dashboard() tea.Model {
	viewport := viewport.New(globals.width, globals.height-10)
	messages := []string{}

	for _, msg := range globals.logger.GetEntries() {
		messages = append(messages, globals.logger.FormatEntry(msg))
	}

	viewport.SetContent(strings.Join(messages, "\n"))

	return DashboardScreen{
		viewport: viewport,
		ready:    true,
	}
}

func (d DashboardScreen) Init() tea.Cmd {
	return tickCmd()
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return dashboardTickMsg(t)
	})
}

func (d DashboardScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		d.viewport.Width = globals.width
		d.viewport.Height = globals.width - 10
		d.updateContent()

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return d, tea.Quit
		}

	case dashboardTickMsg:
		d.updateContent()
		cmds = append(cmds, tickCmd())
	}

	d.viewport, cmd = d.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return d, tea.Batch(cmds...)
}

func (d *DashboardScreen) updateContent() {
	if globals.logger == nil {
		return
	}

	entries := globals.logger.GetEntries()
	var content string
	for _, entry := range entries {
		content += globals.logger.FormatEntry(entry) + "\n"
	}

	d.viewport.SetContent(content)
	d.viewport.GotoBottom()
}

func (d DashboardScreen) View() string {
	return d.viewport.View() + "\n\nPress Ctrl+C to quit"
}
