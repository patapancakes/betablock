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
)

func About(w http.ResponseWriter, r *http.Request) {
	ad := ActionData{Header: "Betablock", Page: "about"}

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
}
