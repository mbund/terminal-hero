package main

import (
	stopwatch "github.com/charmbracelet/bubbles/v2/stopwatch"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/log"
)

type Game struct {
	width     int
	height    int
	stopwatch stopwatch.Model
}

func (m Game) Init() tea.Cmd {
	return m.stopwatch.Init()
}

func (m Game) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		log.Info("key press", "key", msg.Key())
	case tea.KeyReleaseMsg:
		log.Info("key release", "key", msg.Key())
		return m, m.stopwatch.Start()
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	var cmd tea.Cmd
	m.stopwatch, cmd = m.stopwatch.Update(msg)
	return m, cmd
}

func (m Game) View() tea.View {
	result := m.stopwatch.Elapsed().String()

	view := tea.NewView(result)
	view.KeyReleases = true
	return view
}
