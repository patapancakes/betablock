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
	"bytes"
	"fmt"
	"net/http"

	"github.com/patapancakes/betablock/db"
)

func CheckServer(w http.ResponseWriter, r *http.Request) {
	serverId, err := getPaddedServerID(r.URL.Query().Get("serverId"))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to decode serverId: %s", err), http.StatusBadRequest)
		return
	}

	sid, err := db.GetUserServerID(r.URL.Query().Get("user"))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get user serverId: %s", err), http.StatusBadRequest)
		return
	}
	if !bytes.Equal(serverId, sid) {
		http.Error(w, "serverId does not match", http.StatusUnauthorized)
		return
	}

	fmt.Fprint(w, "YES")
}
