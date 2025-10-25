package main

import (
	"log/slog"

	stopwatch "github.com/charmbracelet/bubbles/v2/stopwatch"
	tea "github.com/charmbracelet/bubbletea/v2"
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
		slog.Info("key press", "key", msg.Key())
	case tea.KeyReleaseMsg:
		slog.Info("key release", "key", msg.Key())
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m Game) View() tea.View {
	result := "🬗"

	view := tea.NewView(result)
	view.KeyReleases = true
	return view
}
