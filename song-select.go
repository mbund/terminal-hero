package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/log"
	gotar_hero "github.com/mbund/terminal-hero/pkg/gotar-hero"
)

type Song struct {
	Title      string
	Artist     string
	Album      string
	Genre      string
	Year       string
	Charter    string
	Difficulty int
	Length     float64
}

func NewSong(chart gotar_hero.Chart) Song {
	return Song{
		Title:      chart.Title,
		Artist:     chart.Artist,
		Album:      chart.Album,
		Genre:      chart.Genre,
		Year:       chart.Year,
		Charter:    chart.Charter,
		Difficulty: chart.Difficulty,
		Length:     chart.Length,
	}
}

type SongSelect struct {
	songs    []Song
	selected int
	width    int
	height   int
}

func (m SongSelect) Init() tea.Cmd {
	return nil
}

func (m SongSelect) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "down", "j":
			m.selected = m.selected + 1
			if m.selected == len(m.songs) {
				m.selected = 0
			}
		case "up", "k":
			m.selected = m.selected - 1
			if m.selected == -1 {
				m.selected = len(m.songs) - 1
			}
		case "space", "enter":
			switch m.selected {
			case BUTTON_QUIT:
				return m, tea.Quit
			case BUTTON_PLAY:
				return Game{width: m.width, height: m.height}, nil
			}
		case "q":
			return m, tea.Quit
		}
		log.Info("select", "i", m.selected)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m SongSelect) View() tea.View {
	selected_song := m.songs[m.selected]

	s := lipgloss.NewStyle().SetString(fmt.Sprintf("%v: %v\nArtist: %v\nGenre: %v", m.selected, selected_song.Title, selected_song.Artist, selected_song.Genre))
	result := lipgloss.Place(m.width, m.height, 0.5, 0.5, s.Render())
	view := tea.NewView(result)
	view.KeyReleases = true
	return view
}
