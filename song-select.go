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
			m.selected = m.songWrap(m.selected - 1)
		case "up", "k":
			m.selected = m.songWrap(m.selected + 1)
		case "space", "enter":
			fmt.Println("selected song", m.selected)
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

func (m SongSelect) songWrap(i int) int {
	if i >= len(m.songs) {
		return i - len(m.songs)
	}
	if i < 0 {
		return i + len(m.songs)
	}
	return i
}

func (m SongSelect) View() tea.View {
	previousSong := m.songs[m.songWrap(m.selected-1)]
	selectedSong := m.songs[m.selected]
	nextSong := m.songs[m.songWrap(m.selected+1)]

	lipgloss.NewStyle().Render("")

	s1 := lipgloss.NewStyle().SetString(fmt.Sprintf("%v: %v\nArtist: %v\nGenre: %v", m.songWrap(m.selected-1), previousSong.Title, previousSong.Artist, previousSong.Genre)).Border(lipgloss.NormalBorder()).BorderTop(false).BorderBottom(false).BorderRight(false).BorderLeft(true).PaddingLeft(1)
	s2 := lipgloss.NewStyle().SetString(fmt.Sprintf("%v: %v\nArtist: %v\nGenre: %v", m.selected, selectedSong.Title, selectedSong.Artist, selectedSong.Genre)).Border(lipgloss.NormalBorder()).BorderTop(false).BorderBottom(false).BorderRight(false).BorderLeft(true).BorderForeground(lipgloss.BrightMagenta).Bold(true).PaddingLeft(1)
	s3 := lipgloss.NewStyle().SetString(fmt.Sprintf("%v: %v\nArtist: %v\nGenre: %v", m.songWrap(m.selected+1), nextSong.Title, nextSong.Artist, nextSong.Genre)).Border(lipgloss.NormalBorder()).BorderTop(false).BorderBottom(false).BorderRight(false).BorderLeft(true).PaddingLeft(1)

	list := lipgloss.JoinVertical(0, s1.Render(), lipgloss.NewStyle().Render(""), s2.Render(), lipgloss.NewStyle().Render(""), s3.Render())
	result := lipgloss.Place(m.width, m.height, 0.5, 0.5, list)
	view := tea.NewView(result)
	view.KeyReleases = true
	return view
}
