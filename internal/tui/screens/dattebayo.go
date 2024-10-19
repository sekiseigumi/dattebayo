package screens

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/exp/rand"
)

func getRandomAsciiArt() string {
	asciiDir := "internal/assets/ascii"
	files, err := os.ReadDir(asciiDir)
	if err != nil {
		return "Error loading ASCII art"
	}

	if len(files) == 0 {
		return "No ASCII art files found"
	}

	rand.Seed(uint64(time.Now().UnixNano()))
	randomFile := files[rand.Intn(len(files))]

	content, err := os.ReadFile(filepath.Join(asciiDir, randomFile.Name()))
	if err != nil {
		return "Error reading ASCII art file"
	}

	return string(content)
}

type DattebayoScreen struct {
	asciiArt string
	timer    int
}

func dattebayo() tea.Model {
	return DattebayoScreen{
		asciiArt: getRandomAsciiArt(),
		timer: func() int {
			if globals.config.StartTimer != 0 {
				return globals.config.StartTimer
			} else {
				return 5
			}
		}(),
	}
}

func (d DattebayoScreen) Init() tea.Cmd {
	return tick
}

func tick() tea.Msg {
	time.Sleep(time.Second)
	return tickMsg{}
}

type tickMsg struct{}

func (d DattebayoScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return d, tea.Quit
		}
	case tickMsg:
		d.timer--
		if d.timer <= 0 {
			return screen().Switch(dashboard())
		}
		return d, tick
	}

	return d, nil
}

func (d DattebayoScreen) View() string {
	artLines := strings.Split(strings.TrimSpace(d.asciiArt), "\n")

	message := fmt.Sprintf("\n\nStarting Dattebayo Servers in %d seconds. Press Ctrl + C now to quit.", d.timer)
	content := strings.Join(artLines, "\n") + message

	centeredStyle := lipgloss.NewStyle().
		Width(globals.width).
		Height(globals.height).
		Align(lipgloss.Center).
		AlignVertical(lipgloss.Center)

	centeredContent := centeredStyle.Render(content)

	return centeredContent
}
