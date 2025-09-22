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

package game

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/patapancakes/betablock/db"
)

func JoinServer(w http.ResponseWriter, r *http.Request) {
	sessionId, err := hex.DecodeString(r.URL.Query().Get("sessionId"))
	if err != nil {
		http.Error(w, "Bad response", http.StatusBadRequest)
		return
	}

	serverId, err := getPaddedServerID(r.URL.Query().Get("serverId"))
	if err != nil {
		http.Error(w, "Bad response", http.StatusBadRequest)
		return
	}

	username, err := db.GetUsernameFromSession(r.Context(), sessionId)
	if err != nil {
		http.Error(w, "Bad login", http.StatusOK)
		return
	}
	if r.URL.Query().Get("user") != username {
		http.Error(w, "Bad login", http.StatusOK)
		return
	}

	err = db.SetUserServerID(r.Context(), username, serverId)
	if err != nil {
		http.Error(w, "Bad response", http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, "OK")
}
