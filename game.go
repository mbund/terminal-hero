package main

import (
	"fmt"
	"math"
	"strconv"

	stopwatch "github.com/charmbracelet/bubbles/v2/stopwatch"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	gotar_hero "github.com/mbund/terminal-hero/pkg/gotar-hero"
)

type NotePos struct {
	position float64
	length   float64
}

type Game struct {
	width        int
	height       int
	stopwatch    stopwatch.Model
	mixer        *AudioMixer
	held         []bool
	cursor       gotar_hero.ChartCursor
	prevTime     float64
	notes        [][]NotePos
	accTime      float64
	startedAudio bool
	strumming    bool
	strumInfo    string
	score        float64
}

var (
	NoteSpawn = 450
	NoteSpeed = 200
)

func (m Game) Init() tea.Cmd {
	return m.stopwatch.Init()
}

func (m Game) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.strumming = false
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		log.Info("pressed", "key", msg.Key().Text)
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
		case "j", "k", "space":
			m.strumming = true
		case "q":
			return m, tea.Quit
		}
	case tea.KeyReleaseMsg:
		log.Info("released", "key", msg.Key().Text)
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

	done := m.update()
	if done {
		return m, tea.Quit
	}

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
func renderRow(charWidth int, positions []NotePos, held bool, colors rowColors, ticks_per_char float64) string {
	// a := []rune("\u2588\u2588\u2588\u2588 ")
	// b := []rune("\u2590\u2588\u2588\u2588\u258c")
	result := ""
	result += lipgloss.NewStyle().Foreground(colors.boxBorder).Render("   ┌──────┐") + "\n"
	line := make([]rune, charWidth)
	for i := range charWidth {
		line[i] = ' '
	}
	for _, pos := range positions {
		posChar := floordiv(int(pos.position), 2)
		posMod := mod(int(pos.position), 2)
		char_len := max(5, int(pos.length/ticks_per_char))
		if posMod == 0 {
			for i := range char_len {
				if posChar+i >= 0 && posChar+i < charWidth {
					line[posChar+i] = '\u2588'
				}
			}
		} else {
			for i := range char_len {
				if posChar+i >= 0 && posChar+i < charWidth {
					if i == 0 {
						line[posChar+i] = '\u2590'
					}
					if i == charWidth {
						line[posChar+i] = '\u258c'
					}
					line[posChar+i] = '\u2588'
				}
			}
		}
	}
	for range 2 {
		if held {
			result += lipgloss.NewStyle().Foreground(colors.note).Render(string(line[0:5]))
			result += lipgloss.NewStyle().Foreground(colors.overlap).Background(colors.boxFill).Render(string(line[5:9]))
			result += lipgloss.NewStyle().Foreground(colors.note).Render(string(line[9:])) + "\n"
		} else {
			result += lipgloss.NewStyle().Foreground(colors.note).Render(string(line)) + "\n"
		}
	}
	result += lipgloss.NewStyle().Foreground(colors.boxBorder).Render("   └──────┘") + "\n"
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

func (m *Game) handleEvents(events []any) {
	for i := range events {
		switch u := events[i].(type) {
		case []gotar_hero.Note:
			for j := range u {
				note := u[j]
				if note.Typ > 4 || note.Typ < 0 {
					// silently discared bad notes
					continue
				}
				m.notes[note.Typ] = append(m.notes[note.Typ], NotePos{float64(NoteSpawn), float64(note.Len)})
			}
			// notes
		case *gotar_hero.TempoChange:
			// tempo changes are automatically handled by the cursor
		case *gotar_hero.TSChange:
			// time signature change are automatically handled by the cursor
		}
	}
}

func (m *Game) update() bool {
	m.mixer.mu.Lock()
	newTime := m.mixer.elapsedTime
	deltaTime := newTime - m.prevTime
	m.mixer.mu.Unlock()

	events, adv := m.cursor.NextEvent()
	advTime := float64(adv) / m.cursor.CurrentTicksPerSecond()

	// log.Info("update", "accTime", m.accTime, "advTime", advTime)
	m.accTime += deltaTime

	// if we have accumalated more time than needs to be advanced
	// we need to consume these events
	// log.Info("tick", "adv", adv)
	for m.accTime >= advTime && adv > 0 {
		// consume the events
		m.handleEvents(events)
		m.accTime -= advTime
		m.cursor.AdvanceTick(adv)
		// setup next events
		events, adv = m.cursor.NextEvent()
		advTime = float64(adv) * m.cursor.CurrentTicksPerSecond()
	}

	noteDist := make([]float64, 5)
	for i := range 5 {
		oldPositions := m.notes[i]
		noteDist[i] = math.NaN()
		m.notes[i] = make([]NotePos, 0.0)
		if oldPositions == nil {
			continue
		}

		for j := range oldPositions {
			targetPosition := 10.0
			dist := math.Abs(oldPositions[j].position - targetPosition)
			if dist <= 32.0 && oldPositions[j].length == 0 {
				noteDist[i] = dist
				// delete hit notes that are 0 length
				if m.strumming && m.held[i] {
					log.Info("hit note", "note", i, "dist", dist)
					m.score += 40 * (32 - dist)
				}
			}

			secondsPerChar := 2.0 / float64(NoteSpeed)

			ticksPerChar := m.cursor.CurrentTicksPerSecond() * secondsPerChar

			after := targetPosition - oldPositions[j].position
			if after < oldPositions[j].length/ticksPerChar && after > 0 {
				// we are in the note
				if m.strumming && m.held[i] {
					m.score += deltaTime * 100.0
					log.Info("strumming in held note", "dt", deltaTime, "score", m.score)
					// make this not NaN so this is not considered a false positive
					noteDist[i] = 0
				}
			}

			if oldPositions[j].position+oldPositions[j].length >= -32.0 {
				m.notes[i] = append(m.notes[i], NotePos{oldPositions[j].position - deltaTime*float64(NoteSpeed), oldPositions[j].length})
			} else if oldPositions[j].length != 0 {
				// for now just ignore missed long notes
				log.Info("missed", "note", i)
				m.strumInfo = fmt.Sprintf("miss %d", i)
				m.score -= 50
			}
		}
	}

	total_positions := len(m.notes[0]) + len(m.notes[2]) + len(m.notes[3]) + len(m.notes[4])

	if adv == 0 {
		log.Info("no more events", "positions left", total_positions)
	}

	if adv == 0 && total_positions == 0 {
		// all the notes have passed and there are no more events coming so we are done
		return true
	}

	if m.strumming {
		m.strumInfo = ""
		for i := range 5 {
			if m.held[i] {
				if math.IsNaN(noteDist[i]) {
					m.strumInfo += fmt.Sprintf("false positive %d; ", i)
				} else {
					m.strumInfo += fmt.Sprintf("distance %d %f; ", i, noteDist[i])
				}
			} else {
				if math.IsNaN(noteDist[i]) {
					// true negative
				} else {
					m.strumInfo += fmt.Sprintf("false negative %d; ", i)
				}
			}

		}
	}

	if newTime > (float64(NoteSpawn)-20.0)/float64(NoteSpeed) && !m.startedAudio {
		_, _ = m.mixer.Play("audio.raw", 1.0)
		m.startedAudio = true
	}

	m.prevTime = newTime

	return false
}

func (m Game) View() tea.View {

	green := lipgloss.Color("#19a11b")
	greens := rowColors{
		boxBorder: green,
		note:      green,
		boxFill:   lighten(green, 30),
		overlap:   green,
	}

	red := lipgloss.Color("#b72528")
	reds := rowColors{
		boxBorder: red,
		note:      red,
		boxFill:   lighten(red, 30),
		overlap:   red,
	}

	yellow := lipgloss.Color("#cab50c")
	yellows := rowColors{
		boxBorder: yellow,
		note:      yellow,
		boxFill:   lighten(yellow, 30),
		overlap:   yellow,
	}

	blue := lipgloss.Color("#138ed2")
	blues := rowColors{
		boxBorder: blue,
		note:      blue,
		boxFill:   lighten(blue, 30),
		overlap:   blue,
	}

	orange := lipgloss.Color("#a05206")
	oranges := rowColors{
		boxBorder: orange,
		note:      orange,
		boxFill:   lighten(orange, 20),
		overlap:   orange,
	}

	result := m.strumInfo + " score: " + strconv.Itoa(int(m.score)) + "\n"

	secondsPerChar := 2.0 / float64(NoteSpeed)

	ticksPerChar := m.cursor.CurrentTicksPerSecond() * secondsPerChar

	result += renderRow(m.width, m.notes[0], m.held[0], greens, ticksPerChar)
	result += renderRow(m.width, m.notes[1], m.held[1], reds, ticksPerChar)
	result += renderRow(m.width, m.notes[2], m.held[2], yellows, ticksPerChar)
	result += renderRow(m.width, m.notes[3], m.held[3], blues, ticksPerChar)
	result += renderRow(m.width, m.notes[4], m.held[4], oranges, ticksPerChar)

	view := tea.NewView(result)
	view.KeyReleases = true
	return view
}
