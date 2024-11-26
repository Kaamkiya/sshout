package sshout

import (
	"fmt"
	"errors"
	"net"
	"time"
//	"path/filepath"
	_ "embed"

	"codeberg.org/Kaamkiya/sshout/internal/pkg/auth"
	"codeberg.org/Kaamkiya/sshout/internal/pkg/db"
	"codeberg.org/Kaamkiya/sshout/internal/pkg/chat"

	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/logging"
)

//go:embed banner.txt
var banner string

func Run() {
	hostAddr := net.JoinHostPort("localhost", "2222")

	var err error
	var user db.User
	
	err = db.Open()
	if err != nil {
		log.Error("Failed to open database", "error", err)
	}
	defer db.Close()

	srv, err := wish.NewServer(
		wish.WithBannerHandler(func(ctx ssh.Context) string {
			return fmt.Sprintf(banner, ctx.User())
		}),
		wish.WithIdleTimeout(15 * time.Minute),
		wish.WithAddress(hostAddr),
		wish.WithMiddleware(
			func(next ssh.Handler) ssh.Handler {
				return func(sess ssh.Session) {
					user, err = auth.Run(sess)
					log.Info("User logged in", "user", user, "error", err)
					if err != nil {
						wish.Println(sess, "Authentication failed.")
					} else {
						chat.Run(sess, user)
						next(sess)
					}
				}
			},
			logging.Middleware(),
		),
	)

	if err != nil {
		log.Error("Failed to create server", "error", err)
	}

	log.Info("Starting server", "addr", hostAddr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Failed to run server", "error", err)
	}
	log.Info("Stopping server.")
}
