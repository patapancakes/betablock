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

package login

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"time"

	"github.com/patapancakes/betablock/db"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	username, err := db.GetCanonicalUsername(r.PostForm.Get("user"))
	if err != nil {
		fmt.Fprint(w, "Bad login")
		return
	}

	// password
	err = db.ValidatePassword(username, r.PostForm.Get("password"))
	if err != nil {
		fmt.Fprint(w, "Bad login")
		return
	}

	// ticket
	ticket := make([]byte, 16)
	_, err = rand.Read(ticket)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = db.InsertTicket(username, ticket)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// session
	session := make([]byte, 16)
	_, err = rand.Read(session)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = db.InsertSession(username, session)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// latest version
	latestVersion, err := db.GetUserClientVersionChanged(username)
	if err != nil {
		latestVersion = time.UnixMilli(0)
	}

	fmt.Fprintf(w, "%d:%x:%s:%x", latestVersion.Unix(), ticket, username, session)
}
