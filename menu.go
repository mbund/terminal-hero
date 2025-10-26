package main

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/v2/spinner"
	"github.com/charmbracelet/bubbles/v2/stopwatch"
	gotar_hero "github.com/mbund/terminal-hero/pkg/gotar-hero"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

var (
	normal    = lipgloss.Color("#EEEEEE")
	subtle    = lipgloss.Color("#b0b0b0")
	highlight = lipgloss.Color("#a3e635")
)

type Menu struct {
	width       int
	height      int
	selected    int
	connected   bool
	mixer       *AudioMixer
	sessionData *sessionData
	spinner     spinner.Model
}

func (m Menu) Init() tea.Cmd {
	return tea.Batch(
		connectionStatus(m.sessionData.connected),
		m.spinner.Tick,
	)
}

const (
	BUTTON_PLAY = iota
	BUTTON_LEADERBOARD
	BUTTON_QUIT
	BUTTON_MAX = iota - 1
)

func (m Menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "down", "j":
			m.selected = min(m.selected+1, BUTTON_MAX)
		case "up", "k":
			m.selected = max(m.selected-1, 0)
		case "space", "enter":
			switch m.selected {
			case BUTTON_QUIT:
				return m, tea.Quit
			case BUTTON_PLAY:
				chart, err := gotar_hero.OpenChart("notes.chart")
				if err != nil {
					panic(err)
				}
				cursor, _ := gotar_hero.NewChartCursor(*chart, "ExpertSingle")
				game := Game{width: m.width, height: m.height, stopwatch: stopwatch.New(stopwatch.WithInterval(10 * time.Millisecond)), mixer: m.mixer, held: make([]bool, 5), positions: make([][]float64, 5), cursor: *cursor}
				return game, game.Init()
			}
		}
	case connectionMsg:
		m.connected = msg.connected
		return m, connectionStatus(m.sessionData.connected)
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

const text = `
 ▒▓████████▓▒ ▒▓████████▓▒ ▒▓███████▓▒  ▒▓██████████████▓▒  ▒▓█▓▒ ▒▓███████▓▒   ▒▓██████▓▒  ▒▓█▓▒               ▒▓█▓▒  ▒▓█▓▒ ▒▓████████▓▒ ▒▓███████▓▒   ▒▓██████▓▒   
    ▒▓█▓▒     ▒▓█▓▒        ▒▓█▓▒  ▒▓█▓▒ ▒▓█▓▒  ▒▓█▓▒  ▒▓█▓▒ ▒▓█▓▒ ▒▓█▓▒  ▒▓█▓▒ ▒▓█▓▒  ▒▓█▓▒ ▒▓█▓▒               ▒▓█▓▒  ▒▓█▓▒ ▒▓█▓▒        ▒▓█▓▒  ▒▓█▓▒ ▒▓█▓▒  ▒▓█▓▒  
    ▒▓█▓▒     ▒▓█▓▒        ▒▓█▓▒  ▒▓█▓▒ ▒▓█▓▒  ▒▓█▓▒  ▒▓█▓▒ ▒▓█▓▒ ▒▓█▓▒  ▒▓█▓▒ ▒▓█▓▒  ▒▓█▓▒ ▒▓█▓▒               ▒▓█▓▒  ▒▓█▓▒ ▒▓█▓▒        ▒▓█▓▒  ▒▓█▓▒ ▒▓█▓▒  ▒▓█▓▒  
    ▒▓█▓▒     ▒▓██████▓▒   ▒▓███████▓▒  ▒▓█▓▒  ▒▓█▓▒  ▒▓█▓▒ ▒▓█▓▒ ▒▓█▓▒  ▒▓█▓▒ ▒▓████████▓▒ ▒▓█▓▒               ▒▓████████▓▒ ▒▓██████▓▒   ▒▓███████▓▒  ▒▓█▓▒  ▒▓█▓▒  
    ▒▓█▓▒     ▒▓█▓▒        ▒▓█▓▒  ▒▓█▓▒ ▒▓█▓▒  ▒▓█▓▒  ▒▓█▓▒ ▒▓█▓▒ ▒▓█▓▒  ▒▓█▓▒ ▒▓█▓▒  ▒▓█▓▒ ▒▓█▓▒               ▒▓█▓▒  ▒▓█▓▒ ▒▓█▓▒        ▒▓█▓▒  ▒▓█▓▒ ▒▓█▓▒  ▒▓█▓▒  
    ▒▓█▓▒     ▒▓█▓▒        ▒▓█▓▒  ▒▓█▓▒ ▒▓█▓▒  ▒▓█▓▒  ▒▓█▓▒ ▒▓█▓▒ ▒▓█▓▒  ▒▓█▓▒ ▒▓█▓▒  ▒▓█▓▒ ▒▓█▓▒               ▒▓█▓▒  ▒▓█▓▒ ▒▓█▓▒        ▒▓█▓▒  ▒▓█▓▒ ▒▓█▓▒  ▒▓█▓▒  
    ▒▓█▓▒     ▒▓████████▓▒ ▒▓█▓▒  ▒▓█▓▒ ▒▓█▓▒  ▒▓█▓▒  ▒▓█▓▒ ▒▓█▓▒ ▒▓█▓▒  ▒▓█▓▒ ▒▓█▓▒  ▒▓█▓▒ ▒▓████████▓▒        ▒▓█▓▒  ▒▓█▓▒ ▒▓████████▓▒ ▒▓█▓▒  ▒▓█▓▒  ▒▓██████▓▒ 
`

const guitar = `                                                      
                                              ░░▒▒▓▓▓▓
                                              ▒▒▓▓▓▓▓▓
                                              ▓▓▓▓▒▒  
                                          ░░▓▓▓▓░░    
                                        ░░▒▒░░        
                                      ░░░░▒▒          
                                    ▓▓░░░░            
                                  ░░▒▒░░              
                                ░░▓▓                  
                            ░░▓▓▓▓                    
                  ░░░░    ░░▓▓▒▒                      
              ░░▓▓▓▓    ░░▒▒▒▒                        
              ▓▓▓▓▓▓  ░░▒▒▒▒                          
            ▒▒▓▓▓▓██▒▒▓▓▒▒                            
          ░░▓▓▓▓▓▓▓▓░░▒▒▒▒                            
    ▒▒████▓▓▓▓██▒▒▓▓▓▓▒▒▓▓░░▒▒                        
  ▓▓▓▓▓▓▓▓▓▓▓▓▒▒░░░░████████░░                        
  ▓▓▓▓▓▓▓▓▓▓▒▒▓▓▒▒▓▓██████▓▓                          
  ▓▓▓▓▓▓▓▓▓▓██░░██████▓▓░░                            
  ▒▒▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓██                                
    ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓                                
      ▒▒██▓▓▓▓██▓▓██▒▒                                
        ░░▓▓▓▓▓▓▓▓██                                  
`

func (m Menu) View() tea.View {
	guitarGradiant := lipgloss.NewStyle().Foreground(normal).Render(guitar)

	menu := ""
	for i := range BUTTON_MAX + 1 {
		var button string
		switch i {
		case BUTTON_PLAY:
			button = "Play"
		case BUTTON_LEADERBOARD:
			button = "Leaderboard"
		case BUTTON_QUIT:
			button = "Quit"
		}
		color := subtle
		if m.selected == i {
			color = highlight
		}
		if menu != "" {
			menu += "\n\n"
		}
		menu = lipgloss.JoinVertical(0.0, menu, lipgloss.NewStyle().Foreground(color).Bold(true).Border(lipgloss.NormalBorder()).BorderForeground(color).Padding(1).PaddingLeft(2).Width(54).Render(button))
	}

	connectionCommand := "ssh -T -p 23234 -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no 100.107.230.44 | aplay -f S16_LE -c 2 -r 44100 --buffer-size 1024"
	connectionBlock := AddTitle(lipgloss.NewStyle().Foreground(normal).Border(lipgloss.NormalBorder()).Padding(1, 2).Width(170).Render(connectionCommand), "Connect to audio:")

	var connectionStatus string
	if m.connected {
		connectionStatus = lipgloss.NewStyle().Foreground(highlight).Bold(true).Render("🎶Connected")
	} else {
		connectionStatus = lipgloss.NewStyle().Foreground(highlight).Bold(true).Render(m.spinner.View() + " Waiting for audio connection...")
	}

	result := lipgloss.JoinVertical(0.5,
		lipgloss.NewStyle().Foreground(highlight).Render(text),
		"\n\n\n",
		lipgloss.NewStyle().Foreground(subtle).Render(connectionBlock),
		"\n",
		connectionStatus,
		"\n\n\n",
		lipgloss.JoinHorizontal(0.5,
			guitarGradiant,
			"              ",
			menu,
		),
	)
	result = lipgloss.Place(m.width, m.height, 0.5, 0.5, result)
	view := tea.NewView(result)
	view.KeyReleases = true
	return view
}

type connectionMsg struct {
	connected bool
}

func connectionStatus(ch <-chan bool) tea.Cmd {
	return func() tea.Msg {
		value := <-ch
		return connectionMsg{connected: value}
	}
}

func AddTitle(rendered, title string) string {
	lines := strings.Split(rendered, "\n")
	if len(lines) == 0 {
		return rendered
	}

	// Replace characters 3 to 3+len(title)+2 in the first line
	firstLine := lines[0]
	runes := []rune(firstLine)

	titleWithSpaces := " " + title + " "
	titleRunes := []rune(titleWithSpaces)

	// Make sure we have enough characters to replace
	if len(runes) < 3+len(titleRunes) {
		return rendered
	}

	// Replace the characters
	for i, r := range titleRunes {
		runes[3+i] = r
	}

	lines[0] = string(runes)
	return strings.Join(lines, "\n")
}
