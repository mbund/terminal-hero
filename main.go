package main

// An example Bubble Tea server. This will put an ssh session into alt screen
// and continually print up to date terminal information.

import (
	"context"
	"errors"

	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish/v2"
	"github.com/charmbracelet/wish/v2/activeterm"
	"github.com/charmbracelet/wish/v2/bubbletea"
	"github.com/charmbracelet/wish/v2/logging"
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

	m := Leaderboard{
		width:   pty.Window.Width,
		height:  pty.Window.Height,
		song:    "Song Name",
		entries: []LeadeboardEntry{{"Rock God", 100000, "key1"}, {"The Guitar Hero", 90000, "key2"}, {"EPIC Gamer", 80000, "key3"}, {"s t e v e", 70000, "key3"}, {"anonymous", 65000, "key4"}, {"John Wick", 60000, "key5"}, {"Gandalf", 55000, "key6"}},
	}
	m.AddEntry(LeadeboardEntry{"GOAT", 200000, "key4"})
	m.AddEntry(LeadeboardEntry{"Mark", 75000, "key7"})
	return m, []tea.ProgramOption{}
}
