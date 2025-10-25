package main

import (
	"fmt"
	"log/slog"
	"strings"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

type Menu struct {
	term    string
	profile string
	width   int
	height  int
	bg      string
}

func (m Menu) Init() tea.Cmd {
	return nil
}

func (m Menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyPressMsg:
		slog.Info("Key press")
	case tea.KeyReleaseMsg:
		slog.Info("Release press")

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
	result := lipgloss.JoinVertical(0,
		lipgloss.NewStyle().Foreground(lipgloss.Color("#00cc00")).Render(text),
		"\n\n\n",
		lipgloss.JoinHorizontal(0.5,
			guitarGradiant,
		))
	result = lipgloss.Place(m.width, m.height, 0.5, 0.5, result)
	view := tea.NewView(result)
	view.KeyReleases = true
	return view
}
