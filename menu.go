package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

type Menu struct {
	width    int
	height   int
	selected int
}

func (m Menu) Init() tea.Cmd {
	return nil
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
				return Game{width: m.width, height: m.height}, nil
			}
		}
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
	var guitarGradiant = ""
	var white = 160
	for line := range strings.SplitSeq(guitar, "\n") {
		white -= 4
		guitarGradiant += lipgloss.NewStyle().Foreground(lipgloss.Color(fmt.Sprintf("#%02xff%02x", white, white))).Render(line) + "\n"
	}

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
		color := lipgloss.Color("#004400")
		if m.selected == i {
			color = lipgloss.Color("#009900")
		}
		if menu != "" {
			menu += "\n\n"
		}
		menu = lipgloss.JoinVertical(0.0, menu, lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Background(color).Padding(1).PaddingLeft(2).Width(54).Render(button))
	}

	connectionCommand := "ssh -T -p 23234 -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no 100.107.230.44 | aplay -f S16_LE -c 2 -r 44100 --buffer-size 1024"
	connectionBlock := "┌─ Connect to audio: ─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐\n"
	connectionBlock += "│  " + fmt.Sprintf("%-158s", "") + " │\n"
	connectionBlock += "│  " + fmt.Sprintf("%-158s", connectionCommand) + " │\n"
	connectionBlock += "│  " + fmt.Sprintf("%-158s", "") + " │\n"
	connectionBlock += "└─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘\n"

	result := lipgloss.JoinVertical(0.5,
		lipgloss.NewStyle().Foreground(lipgloss.Color("#00cc00")).Render(text),
		"\n\n\n",
		lipgloss.NewStyle().Foreground(lipgloss.Color("#00cc00")).Render(connectionBlock),
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
