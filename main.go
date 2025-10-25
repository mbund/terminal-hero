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

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish/v2"
	"github.com/charmbracelet/wish/v2/bubbletea"
	"github.com/charmbracelet/wish/v2/logging"
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

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	// This should never fail, as we are using the activeterm middleware.
	pty, _, _ := s.Pty()

	tuple := newAssociationTuple(s)
	sessionData := getOrCreateSessionData(tuple)
	m := Menu{
		width:  pty.Window.Width,
		height: pty.Window.Height,
		mixer:  sessionData.mixer,
	}

	return m, []tea.ProgramOption{}

}
