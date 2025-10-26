package main

import (
	"fmt"
	"slices"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

type LeadeboardEntry struct {
	name   string
	score  int
	pubkey string
}

type Leaderboard struct {
	song    string
	entries []LeadeboardEntry
	width   int
	height  int
}

func (m Leaderboard) Init() tea.Cmd {
	return nil
}

func (m Leaderboard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyPressMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Leaderboard) View() tea.View {
	s := ""
	for i := range m.entries {
		entry := m.entries[i]
		style := lipgloss.NewStyle().SetString(fmt.Sprintf("%2v. %32v: %7v\n", i+1, entry.name, entry.score)).PaddingTop(1)

		switch i {
		case 0:
			style = style.Foreground(lipgloss.Color("220"))
		case 1:
			style = style.Foreground(lipgloss.Color("7"))
		case 2:
			style = style.Foreground(lipgloss.Color("94"))
		}

		s += style.Render()
	}

	result := lipgloss.Place(m.width, m.height, 0.5, 0.5, s)
	view := tea.NewView(result)
	view.KeyReleases = true
	return view
}

func (m *Leaderboard) AddEntry(entry LeadeboardEntry) {
	present := false
	for i := range m.entries {
		if m.entries[i].pubkey == entry.pubkey {
			m.entries[i].name = entry.name
			m.entries[i].score = max(entry.score, entry.score)
			present = true
			break
		}
	}
	if !present {
		m.entries = append(m.entries, entry)
	}
	slices.SortStableFunc(m.entries, func(a LeadeboardEntry, b LeadeboardEntry) int {
		if a.score < b.score {
			return 1
		}
		if a.score == b.score {
			return 0
		}
		return -1
	})
}
