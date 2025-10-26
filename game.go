package main

import (
	"fmt"

	stopwatch "github.com/charmbracelet/bubbles/v2/stopwatch"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	gotar_hero "github.com/mbund/terminal-hero/pkg/gotar-hero"
)

type Game struct {
	width     int
	height    int
	stopwatch stopwatch.Model
	mixer     *AudioMixer
	held      []bool
	cursor    gotar_hero.ChartCursor
	prevTime  float64
	positions [][]float64
}

func (m Game) Init() tea.Cmd {
	return m.stopwatch.Init()
}

func (m Game) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		log.Info("key press", "key", msg.Key())
		switch msg.Key().Text {
		case "1":
			m.held[0] = true
		case "2":
			m.held[1] = true
		case "3":
			m.held[2] = true
		case "4":
			m.held[3] = true
		case "5":
			m.held[4] = true
		case "q":
			return m, tea.Quit
		}
	case tea.KeyReleaseMsg:
		switch msg.Key().Text {
		case "1":
			m.held[0] = false
		case "2":
			m.held[1] = false
		case "3":
			m.held[2] = false
		case "4":
			m.held[3] = false
		case "5":
			m.held[4] = false
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	var cmd tea.Cmd
	m.stopwatch, cmd = m.stopwatch.Update(msg)
	m.update()

	return m, cmd
}

func floordiv(a, b int) int {
	return (a - mod(a, b)) / b
}

func mod(a, b int) int {
	return (a%b + b) % b
}

type rowColors struct {
	boxBorder lipgloss.Color
	boxFill   lipgloss.Color
	note      lipgloss.Color
	overlap   lipgloss.Color
}

// postitions is an array of half-character coordinates
func renderRow(charWidth int, positions []float64, held bool, colors rowColors) string {
	// boxLeft := 8
	// boxRight := 20
	a := []rune("\u2588\u2588\u2588\u2588 ")
	b := []rune("\u2590\u2588\u2588\u2588\u258c")
	result := ""
	result += lipgloss.NewStyle().Foreground(colors.boxBorder).Render("    ┌────┐") + "\n"
	line := make([]rune, charWidth)
	for i := range charWidth {
		line[i] = ' '
	}
	for _, pos := range positions {
		posChar := floordiv(int(pos), 2)
		posMod := mod(int(pos), 2)
		if posMod == 0 {
			for i := range 5 {
				if posChar+i >= 0 && posChar+i < charWidth {
					line[posChar+i] = a[i]
				}
			}
		} else {
			for i := range 5 {
				if posChar+i >= 0 && posChar+i < charWidth {
					line[posChar+i] = b[i]
				}
			}
		}
	}
	for _ = range 2 {
		if held {
			result += lipgloss.NewStyle().Foreground(colors.note).Render(string(line[0:5]))
			result += lipgloss.NewStyle().Foreground(colors.overlap).Background(colors.boxFill).Render(string(line[5:9]))
			result += lipgloss.NewStyle().Foreground(colors.note).Render(string(line[9:])) + "\n"
		} else {
			result += lipgloss.NewStyle().Foreground(colors.note).Render(string(line)) + "\n"
		}
	}
	result += lipgloss.NewStyle().Foreground(colors.boxBorder).Render("    └────┘") + "\n"
	return result
}

func lighten(c lipgloss.Color, percent float64) lipgloss.Color {
	r16, g16, b16, _ := c.RGBA()

	r := float64(r16) / 257
	g := float64(g16) / 257
	b := float64(b16) / 257

	f := percent / 100
	l := func(x float64) uint8 {
		val := x + (255-x)*f
		if val > 255 {
			val = 255
		}
		return uint8(val)
	}

	r2 := l(r)
	g2 := l(g)
	b2 := l(b)

	return lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", r2, g2, b2))
}

func (m *Game) update() {
	newTime := m.stopwatch.Elapsed().Seconds()
	deltaTime := newTime - m.prevTime

	for i := range 5 {
		if m.positions[i] == nil {
			m.positions[i] = make([]float64, 0.0)
		}
		for j := range m.positions[i] {
			m.positions[i][j] -= deltaTime * 64.0
		}
	}

	if m.prevTime == 0 {
		m.positions[1] = append(m.positions[1], 100)
	}

	m.prevTime = newTime
}

func (m Game) View() tea.View {

	green := lipgloss.Color("#19a11b")
	greens := rowColors{
		boxBorder: green,
		note:      green,
		boxFill:   lighten(green, 20),
		overlap:   green,
	}

	red := lipgloss.Color("#b72528")
	reds := rowColors{
		boxBorder: red,
		note:      red,
		boxFill:   lighten(red, 20),
		overlap:   red,
	}

	yellow := lipgloss.Color("#cab50c")
	yellows := rowColors{
		boxBorder: yellow,
		note:      yellow,
		boxFill:   lighten(yellow, 20),
		overlap:   yellow,
	}

	blue := lipgloss.Color("#138ed2")
	blues := rowColors{
		boxBorder: blue,
		note:      blue,
		boxFill:   lighten(blue, 20),
		overlap:   blue,
	}

	orange := lipgloss.Color("#a05206")
	oranges := rowColors{
		boxBorder: orange,
		note:      orange,
		boxFill:   lighten(orange, 20),
		overlap:   orange,
	}

	result := renderRow(m.width, m.positions[0], m.held[0], greens)
	result += renderRow(m.width, m.positions[1], m.held[1], reds)
	result += renderRow(m.width, m.positions[2], m.held[2], yellows)
	result += renderRow(m.width, m.positions[3], m.held[3], blues)
	result += renderRow(m.width, m.positions[4], m.held[4], oranges)

	view := tea.NewView(result)
	view.KeyReleases = true
	return view
}
