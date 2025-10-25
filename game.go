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
	mixer     *AudioMixer
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
	x2 := int(m.stopwatch.Elapsed().Seconds() * 8)
	x := x2 / 2
	mod := x2 % 2

	result := ""

	a := []rune("\u2588\u2588\u2588\u2588 ")
	b := []rune("\u2590\u2588\u2588\u2588\u258c")

	width := 100
	result += "                                                         ┌────┐\n"
	line := make([]rune, width)
	for i := range width {
		line[i] = ' '
		if i >= x && i-x < 5 {
			if mod == 0 {
				line[i] = a[i-x]
			} else {
				line[i] = b[i-x]
			}
		}
	}
	result += string(line) + "\n"
	result += string(line) + "\n"
	result += "                                                         └────┘\n"

	view := tea.NewView(result)
	view.KeyReleases = true
	return view
}
