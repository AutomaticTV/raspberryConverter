package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var encryptionKey = "something-very-secret"
var loggedUserSession = sessions.NewCookieStore([]byte(encryptionKey))

func PasswordIsCorrect(username string, password string) bool {
	hashedPassword, err := getHashedPasswword(username)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if bcrypt.CompareHashAndPassword(hashedPassword, []byte(password)) != nil {
		return false
	}
	return true
}

func IsLoggedIn(r *http.Request) bool {
	session, err := loggedUserSession.Get(r, "authenticated-user-session")
	if err != nil || session.Values["username"] != "admin" {
		return false
	}
	return true
}

func Login(w http.ResponseWriter, r *http.Request, username string) error {
	session, err := loggedUserSession.New(r, "authenticated-user-session")
	session.Values["username"] = username
	if err == nil {
		err = session.Save(r, w)
	}
	if err != nil {
		fmt.Println(err)
		return errors.New("Error creating session")
	}
	return err
}

func Logout(w http.ResponseWriter, r *http.Request) error {
	session, _ := loggedUserSession.Get(r, "authenticated-user-session")
	session.Values["username"] = ""
	return session.Save(r, w)
}

func UpdatePassword(username string, oldPass string, newPass string, rePass string) error {
	if newPass != rePass {
		return errors.New("The second entry of the new password doesn't match.")
	}
	err := updatePassword(username, oldPass, newPass)
	if err != nil {
		fmt.Println(err)
		return errors.New("Incorrect username or old password.")
	}
	return nil
}
