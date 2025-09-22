package frontend

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/patapancakes/betablock/db"
	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	ad := ActionData{Header: "Login", Page: "login"}

	username, err := UsernameFromRequest(r)
	if err != nil && err != http.ErrNoCookie {
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
		return
	}

	ad.Username = username

	if r.Method == "GET" {
		err := t.Execute(w, ad)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to execute template: %s", err), http.StatusInternalServerError)
			return
		}

		return
	}

	err = r.ParseForm()
	if err != nil {
		Error(w, ad, "An error occured while parsing your request")
		return
	}

	// validate username and password
	username = r.PostFormValue("username")

	err = db.ValidatePassword(r.Context(), username, r.PostFormValue("password"))
	if err != nil {
		var reason string
		switch err {
		case sql.ErrNoRows:
			reason = "The specified user doesn't exist"
		case bcrypt.ErrMismatchedHashAndPassword:
			reason = "The password is incorrect"
		default:
			reason = "An unknown error occured during account validation"
		}

		Error(w, ad, reason)
		return
	}

	session := make([]byte, 16)
	_, err = rand.Read(session)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	err = db.InsertSession(r.Context(), username, session)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    base64.StdEncoding.EncodeToString(session),
		MaxAge:   60 * 60 * 24 * 7,
		HttpOnly: true,
	})

	ad.Success = true
	ad.Username = username

	err = t.Execute(w, ad)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to execute template: %s", err), http.StatusInternalServerError)
		return
	}
}
