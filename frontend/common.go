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
	"html/template"
	"net/http"
)

type ActionData struct {
	Page    string
	Header  string
	Error   string
	Success bool

	Versions []string
}

const maxUploadSize = 1024 * 16

var t = template.Must(template.New("main.html").ParseGlob("templates/frontend/*.html"))

func Error(w http.ResponseWriter, ad ActionData, reason string) error {
	ad.Error = reason
	err := t.Execute(w, ad)
	if err != nil {
		return err
	}

	return nil
}
