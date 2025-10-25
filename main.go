package main

// An example Bubble Tea server. This will put an ssh session into alt screen
// and continually print up to date terminal information.

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
)

const (
	host = "0.0.0.0"
	port = "23234"
)

type associationTuple struct {
	pubkey        string
	user          string
	ip            string
	clientVersion string
}

func PublicKeyToAuthString(pubKey ssh.PublicKey) string {
	keyBytes := base64.StdEncoding.EncodeToString(pubKey.Marshal())

	return fmt.Sprintf("%s %s", pubKey.Type(), keyBytes)
}

func newAssociationTuple(sess ssh.Session) associationTuple {
	return associationTuple{
		pubkey:        PublicKeyToAuthString(sess.PublicKey()),
		user:          sess.User(),
		ip:            sess.RemoteAddr().(*net.TCPAddr).IP.String(),
		clientVersion: sess.Context().ClientVersion(),
	}
}

type sessionData struct {
	mixer *AudioMixer
}

var (
	associations   map[associationTuple]*sessionData = make(map[associationTuple]*sessionData)
	associationsMu sync.RWMutex
)

func getOrCreateSessionData(tuple associationTuple) *sessionData {
	associationsMu.Lock()
	defer associationsMu.Unlock()

	if data, exists := associations[tuple]; exists {
		return data
	}

	const channels = 2
	const mixAmp = 1.0
	const framesPerWrite = 128
	const sampleRate = 44100
	const bytesPerSample = 2

	data := &sessionData{
		mixer: NewAudioMixer(channels, mixAmp, framesPerWrite, sampleRate, bytesPerSample),
	}
	associations[tuple] = data
	return data
}

func AudioMiddleware() wish.Middleware {
	return func(next ssh.Handler) ssh.Handler {
		return func(sess ssh.Session) {
			_, _, active := sess.Pty()
			tuple := newAssociationTuple(sess)
			if active {
				log.Info("Inserting association", "tuple", tuple)
				getOrCreateSessionData(tuple)
				next(sess)
				return
			}

			sessionData := getOrCreateSessionData(tuple)
			if sessionData == nil || sessionData.mixer == nil {
				log.Error("failed to get session data")
				_ = sess.Exit(1)
				return
			}
			log.Info("associate", "tuple", tuple, "session", sessionData)
			sendAudio(sess, sessionData.mixer)
			_ = sess.Exit(0)
		}
	}
}

func CleanupMiddleware() wish.Middleware {
	return func(next ssh.Handler) ssh.Handler {
		return func(sess ssh.Session) {
			log.Info("cleanup")
			next(sess)
		}
	}
}

func main() {
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithPublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
			return true
		}),
		wish.WithMiddleware(
			CleanupMiddleware(),
			bubbletea.Middleware(teaHandler),
			AudioMiddleware(),
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

// You can wire any Bubble Tea model up to the middleware with a function that
// handles the incoming ssh.Session. Here we just grab the terminal info and
// pass it to the new model. You can also return tea.ProgramOptions (such as
// tea.WithAltScreen) on a session by session basis.
func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	// This should never fail, as we are using the activeterm middleware.
	pty, _, _ := s.Pty()

	// When running a Bubble Tea app over SSH, you shouldn't use the default
	// lipgloss.NewStyle function.
	// That function will use the color profile from the os.Stdin, which is the
	// server, not the client.
	// We provide a MakeRenderer function in the bubbletea middleware package,
	// so you can easily get the correct renderer for the current session, and
	// use it to create the styles.
	// The recommended way to use these styles is to then pass them down to
	// your Bubble Tea model.
	renderer := bubbletea.MakeRenderer(s)
	txtStyle := renderer.NewStyle().Foreground(lipgloss.Color("10"))
	quitStyle := renderer.NewStyle().Foreground(lipgloss.Color("8"))

	bg := "light"
	if renderer.HasDarkBackground() {
		bg = "dark"
	}

	tuple := newAssociationTuple(s)
	sessionData := getOrCreateSessionData(tuple)

	m := model{
		term:         pty.Term,
		profile:      renderer.ColorProfile().Name(),
		width:        pty.Window.Width,
		height:       pty.Window.Height,
		bg:           bg,
		txtStyle:     txtStyle,
		quitStyle:    quitStyle,
		mixer:        sessionData.mixer,
		activeHandle: nil,
	}
	return m, []tea.ProgramOption{tea.WithAltScreen()}
}

// Just a generic tea.Model to demo terminal information of ssh.
type model struct {
	term         string
	profile      string
	width        int
	height       int
	bg           string
	txtStyle     lipgloss.Style
	quitStyle    lipgloss.Style
	mixer        *AudioMixer
	activeHandle *PlaybackHandle
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "1":
			log.Info("got 1 - playing notification")
			_, err := m.mixer.Play("strum.raw", 1)
			if err != nil {
				log.Errorf("failed to play audio: %v", err)
				return m, tea.Quit
			}
			return m, nil
		case "2":
			log.Info("got 2 - playing notification")
			_, err := m.mixer.Play("strum2.raw", 2)
			if err != nil {
				log.Errorf("failed to play audio: %v", err)
				return m, tea.Quit
			}
			return m, nil
		case "r":
			log.Info("got r - playing audio")
			handle, err := m.mixer.Play("audio.raw", 1)
			if err != nil {
				log.Errorf("failed to play audio: %v", err)
				return m, tea.Quit
			}
			m.activeHandle = handle

			go func() {
				ticker := time.NewTicker(500 * time.Millisecond)
				defer ticker.Stop()

				for range ticker.C {
					if !handle.IsPlaying() {
						m.activeHandle = nil
						break
					}
					log.Infof("audio.raw Progress: %d%%", int(handle.Progress()*100))
				}
			}()
			return m, nil
		case "s":
			log.Info("got s - stopping active audio")
			if m.activeHandle != nil && m.activeHandle.IsPlaying() {
				go m.activeHandle.Stop()
				m.activeHandle = nil
			}
			return m, nil
		case "p", " ":
			paused := m.mixer.TogglePause()
			if paused {
				log.Info("audio paused")
			} else {
				log.Info("audio resumed")
			}
			return m, nil
		}
	}
	return m, nil
}

func (m model) View() string {
	pauseStatus := ""
	if m.mixer.IsPaused() {
		pauseStatus = " [PAUSED]"
	}

	s := fmt.Sprintf("Your term is %s\nYour window size is %dx%d\nBackground: %s\nColor Profile: %s%s",
		m.term, m.width, m.height, m.bg, m.profile, pauseStatus)

	instructions := "Press '1' or '2' for short sounds, 'r' for long audio\n"
	instructions += "Press 's' to stop the active audio\n"
	instructions += "Press 'p' or space to pause/resume audio\n"

	return m.txtStyle.Render(s) + "\n\n" + m.quitStyle.Render(instructions+"Press 'q' to quit\n")
}
