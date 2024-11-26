package auth

import (
	"fmt"
	"errors"

	"codeberg.org/Kaamkiya/sshout/internal/pkg/db"

	"golang.org/x/crypto/bcrypt"

	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/huh"
)

func Run(sess ssh.Session) (db.User, error) {
	var hasAccount bool
	var password string
	var username string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Do you already have an account?").
				Value(&hasAccount),
			huh.NewInput().
				Title("Username").
				Description("Enter your preferred username if you don't have an account").
				Value(&username),
			huh.NewInput().
				Title("Password").
				Value(&password).
				EchoMode(huh.EchoModePassword).
				CharLimit(70).
				Validate(func(pw string) error {
					if len(pw) < 6 {
						return errors.New("Password too short")
					}
					return nil
				}),
		),
	).WithOutput(sess).
		WithInput(sess)
	
	if err := form.Run(); err != nil {
		fmt.Fprintln(sess, "Failed to read input, sorry.")
		return db.User{}, err
	}

	user, err := db.GetUser(username)
	if err != nil && !hasAccount {
		fmt.Fprintln(sess, "Welcome to sshout! Creating your account...")

		passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 16)

		if err != nil {
			return db.User{}, err
		}

		fmt.Fprintln(sess, "Success!")
		return db.AddUser(username, string(passwordHash))
	} else if err != nil {
		fmt.Fprintln(sess, "Error: couldn't find this account.")
		return user, err
	} else {
		fmt.Fprintln(sess, "Authenticating...")
		if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
			fmt.Fprintln(sess, "Welcome back!")
		}

		return user, nil
	}

	return user, errors.New("Error: this should be unreachable")
}
