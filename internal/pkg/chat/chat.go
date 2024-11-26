package chat

import (
	"fmt"

	"codeberg.org/Kaamkiya/sshout/internal/pkg/db"

	"github.com/charmbracelet/ssh"
)

func Run(sess ssh.Session, user db.User) {
	fmt.Fprintf(sess, "Someday, something fantastic will be here, %s.\n", user.Username)
}
