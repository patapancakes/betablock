/*
	betablock - block server emulator
	Copyright (C) 2025  Pancakes <patapancakes@pagefault.games>

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU Affero General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU Affero General Public License for more details.

	You should have received a copy of the GNU Affero General Public License
	along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package frontend

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/patapancakes/betablock/db"

	"golang.org/x/crypto/bcrypt"
)

func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	ad := ActionData{Header: "Delete Account", Page: "delaccount"}

	if r.Method == "GET" {
		err := t.Execute(w, ad)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to execute template: %s", err), http.StatusInternalServerError)
			return
		}

		return
	}

	err := r.ParseForm()
	if err != nil {
		Error(w, ad, "An error occured while parsing your request")
		return
	}

	// validate username and password
	username := r.PostFormValue("username")

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

	err = db.DeleteAccount(r.Context(), username)
	if err != nil {
		Error(w, ad, "An error occured while deleting the account")
		return
	}

	ad.Success = true

	err = t.Execute(w, ad)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to execute template: %s", err), http.StatusInternalServerError)
		return
	}
}
