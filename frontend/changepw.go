package frontend

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/patapancakes/betablock/db"
	"golang.org/x/crypto/bcrypt"
)

func ChangePW(w http.ResponseWriter, r *http.Request) {
	ad := ActionData{Header: "Change Password", Page: "changepw"}

	username, err := UsernameFromRequest(r)
	if err != nil {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

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

	// validate username and password
	err = db.ValidatePassword(r.Context(), ad.Username, r.PostFormValue("password"))
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

	err = db.UpdatePassword(r.Context(), username, r.PostFormValue("newpassword"))
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	ad.Success = true
	ad.Username = username

	err = t.Execute(w, ad)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to execute template: %s", err), http.StatusInternalServerError)
		return
	}
}
