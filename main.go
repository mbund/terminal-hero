package main

// An example Bubble Tea server. This will put an ssh session into alt screen
// and continually print up to date terminal information.

import (
	"context"
	"errors"
	"fmt"

	// "fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	_ "github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
)

const (
	host = "localhost"
	port = "23234"
)

func main() {
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server", "host", host, "port", port)
	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	// This should never fail, as we are using the activeterm middleware.
	pty, _, _ := s.Pty()

	renderer := bubbletea.MakeRenderer(s)

	bg := "light"
	if renderer.HasDarkBackground() {
		bg = "dark"
	}

	m := Menu{
		term:     pty.Term,
		profile:  renderer.ColorProfile().Name(),
		width:    pty.Window.Width,
		height:   pty.Window.Height,
		bg:       bg,
		renderer: renderer,
	}
	return m, []tea.ProgramOption{tea.WithAltScreen()}
}

// Just a generic tea.Model to demo terminal information of ssh.
type Menu struct {
	term     string
	profile  string
	width    int
	height   int
	bg       string
	renderer *lipgloss.Renderer
}

func (m Menu) Init() tea.Cmd {
	return nil
}

func (m Menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
                                                                                                                                                                                                                                                                                                   
                        ░░░░
                  ░░▓▓▓▓░░
                  ▓▓▓▓▓▓
                ██▓▓▓▓
                ██▓▓▓▓▒▒
              ▒▒▓▓▓▓▓▓▓▓
              ▒▒▓▓▓▓▓▓▓▓
              ▓▓▓▓▓▓▓▓░░
              ░░▓▓▓▓▓▓
                ▓▓██▓▓
                ▓▓▓▓▓▓
                ████▓▓
                ██▓▓▓▓
                ██▓▓▓▓
                ██▓▓▓▓
                ██▓▓▓▓
                ██▓▓▓▓
                ▓▓▓▓▓▓
                ██▓▓▓▓
                ██▓▓▓▓
                ██▓▓▓▓
                ██▓▓▓▓
                ██▓▓▓▓
      ░░        ██▓▓▒▒
  ░░▓▓▓▓        ██▓▓▓▓
  ▓▓▓▓▓▓        ██▓▓▓▓
  ▓▓▓▓▓▓▓▓    ░░▓▓▓▓▓▓
  ▒▒▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓        ▓▓▓▓░░
  ░░▓▓▓▓▓▓▓▓░░░░▓▓▓▓▓▓      ▒▒░░░░
    ▓▓▓▓▓▓▓▓▒▒        ▓▓▓▓▒▒░░  ░░
    ▓▓▓▓▓▓▓▓            ░░      ▒▒
    ▒▒▓▓▓▓▓▓▒▒▒▒██████          ▒▒
      ▓▓▓▓▓▓                  ▒▒
    ░░▓▓▓▓▓▓  ▒▒██████        ▓▓
    ▓▓▓▓▓▓░░
  ░░▓▓▓▓▓▓    ████              ▓▓
  ▓▓▓▓▓▓▓▓        ████  ██      ░░▓▓
  ▓▓▓▓▓▓▓▓░░  ░░▒▒▒▒▒▒      ██▓▓  ▒▒
▓▓▓▓▓▓▓▓▓▓▓▓▓▓▒▒▓▓▓▓▓▓▓▓▓▓▒▒░░  ██░░▒▒
▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓░░  ▓▓
▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓
▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓
░░▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓░░
  ░░▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▒▒
          ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓░░
`

func (m Menu) View() string {
	var guitarGradiant = ""
	var white = 180
	for line := range strings.SplitSeq(guitar, "\n") {
		white -= 2
		guitarGradiant += m.renderer.NewStyle().Foreground(lipgloss.Color(fmt.Sprintf("#ff%02x%02x", white, white))).Render(line) + "\n"
	}
	var textGradiant = ""
	for line := range strings.SplitSeq(text, "\n") {
		for _, c := range line {
			var color string
			switch c {
			case '█':
				color = "#ffffff"
			case '▓':
				c = '█'
				color = "#ff4444"
			case '▒':
				c = '█'
				color = "#ff0000"
			}
			textGradiant += m.renderer.NewStyle().Foreground(lipgloss.Color(color)).Render(fmt.Sprintf("%c", c))
		}
		textGradiant += "\n"
	}
	result := lipgloss.JoinVertical(0,
		m.renderer.NewStyle().Foreground(lipgloss.Color("#ff0000")).Render(textGradiant),
		lipgloss.JoinHorizontal(0.5,
			guitarGradiant,
		))
	return m.renderer.Place(m.width, m.height, 1, 0.5, result)
}
