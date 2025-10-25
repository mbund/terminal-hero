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
	y2 := int(m.stopwatch.Elapsed().Seconds() * 2)
	y := y2 / 2
	ymod := y2 % 2

	result := ""

	a1 := "\u2588\u2588\u2588\u2588"
	a2 := "\u2588\u2588\u2588\u2588"
	b1 := "\u2584\u2584\u2584\u2584"
	b2 := "\u2588\u2588\u2588\u2588"
	b3 := "\u2580\u2580\u2580\u2580"

	for i := range 20 {
		diff := (i - y)
		switch ymod {
		case 0:
			switch diff {
			case 0:
				result += a1
			case 1:
				result += a2
			default:
				result += "    "
			}
		case 1:
			switch diff {
			case 0:
				result += b1
			case 1:
				result += b2
			case 2:
				result += b3
			default:
				result += "    "
			}

		}
		result += "\n"
	}

	view := tea.NewView(result)
	view.KeyReleases = true
	return view
}
