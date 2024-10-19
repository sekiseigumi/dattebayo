package screens

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sekiseigumi/dattebayo/internal/dns"
	"github.com/sekiseigumi/dattebayo/internal/logger"
	"github.com/sekiseigumi/dattebayo/shared"
	"golang.org/x/term"
)

type ScreenSwitcher struct {
	currentScreen tea.Model
}

type Globals struct {
	width     int
	height    int
	config    shared.Config
	dnsServer *dns.DNSServer
	logger    *logger.Logger
}

var globals Globals

func (s ScreenSwitcher) Init() tea.Cmd {
	return s.currentScreen.Init()
}

func (s ScreenSwitcher) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var model tea.Model

	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		globals.width, globals.height = m.Width, m.Height
	}

	model, cmd = s.currentScreen.Update(msg)

	return ScreenSwitcher{currentScreen: model}, cmd
}

func (s ScreenSwitcher) View() string {
	return s.currentScreen.View()
}

func (s ScreenSwitcher) Switch(screen tea.Model) (tea.Model, tea.Cmd) {
	s.currentScreen = screen
	return s.currentScreen, s.currentScreen.Init()
}

func screen() ScreenSwitcher {
	screen := dattebayo()

	return ScreenSwitcher{
		currentScreen: screen,
	}
}

func Initialize(config shared.Config, dnsServer *dns.DNSServer, log *logger.Logger) tea.Model {

	width, height, err := term.GetSize(int(os.Stdout.Fd()))

	if err != nil {
		width = 80
		height = 30
	}

	globals.width = width
	globals.height = height
	globals.config = config
	globals.dnsServer = dnsServer
	globals.logger = log

	return screen()
}
