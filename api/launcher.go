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

package api

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/patapancakes/betablock/db"
)

func Login(w http.ResponseWriter, r *http.Request) {
	username, err := db.GetCanonicalUsername(r.Context(), r.PostFormValue("user"))
	if err != nil {
		http.Error(w, "Bad login", http.StatusOK)
		return
	}

	// password
	err = db.ValidatePassword(r.Context(), username, r.PostFormValue("password"))
	if err != nil {
		http.Error(w, "Bad login", http.StatusOK)
		return
	}

	// ticket
	ticket := make([]byte, 16)
	_, err = rand.Read(ticket)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	err = db.InsertTicket(r.Context(), username, ticket)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// session
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

	// latest version
	latestVersion, err := db.GetUserClientVersionChanged(r.Context(), username)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		latestVersion = time.UnixMilli(0)
	}

	version, _ := db.GetUserClientVersion(r.Context(), username)
	if version == "realtime" {
		version, latestVersion, err = db.GetRealtimeVersion(r.Context())
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
	}

	fmt.Fprintf(w, "%d:%x:%s:%x", latestVersion.Unix(), ticket, username, session)
}
