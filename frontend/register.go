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
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/patapancakes/betablock/db"
)

var isValidUsername = regexp.MustCompile("^[A-Za-z0-9_]{3,16}$").MatchString

func Register(w http.ResponseWriter, r *http.Request) {
	ad := ActionData{Header: "Register", Page: "register"}

	username, err := UsernameFromRequest(r)
	if err != nil && err != http.ErrNoCookie {
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
		return
	}

	ad.Username = username

	// show register page
	if r.Method == "GET" {
		err := t.Execute(w, ad)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to execute template: %s", err), http.StatusInternalServerError)
			return
		}

		return
	}

	// try to register, show success page if ok
	err = r.ParseForm()
	if err != nil {
		Error(w, ad, "An error occured while parsing your request")
		return
	}

	username = strings.TrimSpace(r.PostForm.Get("username"))
	if !isValidUsername(username) {
		Error(w, ad, "The username specified is invalid")
		return
	}

	password := strings.TrimSpace(r.PostForm.Get("password"))
	if len(password) > 72 {
		Error(w, ad, "The password specified is too long")
		return
	}

	err = db.InsertAccount(r.Context(), username, password)
	if err != nil {
		Error(w, ad, "An error occured while creating the account (username taken?)")
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
