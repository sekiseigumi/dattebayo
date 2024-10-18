package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sekiseigumi/dattebayo/internal/tui/screens"
	"github.com/sekiseigumi/dattebayo/shared"
)

func NewTUI(config shared.Config) *tea.Program {
	return tea.NewProgram(screens.Initialize(config), tea.WithAltScreen())
}
