package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sekiseigumi/dattebayo/internal/dns"
	"github.com/sekiseigumi/dattebayo/internal/logger"
	"github.com/sekiseigumi/dattebayo/internal/tui/screens"
	"github.com/sekiseigumi/dattebayo/shared"
)

func NewTUI(config shared.Config, dnsServer *dns.DNSServer, log *logger.Logger) *tea.Program {
	return tea.NewProgram(screens.Initialize(config, dnsServer, log), tea.WithAltScreen())
}
