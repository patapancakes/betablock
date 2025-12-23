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
	"encoding/base64"
	"html/template"
	"net/http"

	"github.com/patapancakes/betablock/db"
)

type ActionData struct {
	Page    string
	Header  string
	Error   string
	Success bool

	Username string
	Version  string

	Versions []string
}

const maxUploadSize = 1024 * 16

var t = template.Must(template.New("main.html").ParseGlob("templates/*.html"))

func Error(w http.ResponseWriter, ad ActionData, reason string) error {
	ad.Error = reason
	err := t.Execute(w, ad)
	if err != nil {
		return err
	}

	return nil
}

func UsernameFromRequest(r *http.Request) (string, error) {
	sessionCookie, err := r.Cookie("session")
	if err != nil {
		return "", err
	}
	session, err := base64.StdEncoding.DecodeString(sessionCookie.Value)
	if err != nil {
		return "", err
	}
	username, err := db.GetUsernameFromSession(r.Context(), session)
	if err != nil {
		return "", err
	}

	return username, nil
}
