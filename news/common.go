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

package news

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/patapancakes/betablock/db"
)

//go:embed templates
var templatesFS embed.FS
var t = template.Must(template.New("main.html").Funcs(template.FuncMap{"raw": func(s string) template.HTML { return template.HTML(s) }}).ParseFS(templatesFS, "templates/*.html"))

//go:embed assets
var assetsFS embed.FS

var AssetsFS, _ = fs.Sub(assetsFS, "assets")

func Handle(w http.ResponseWriter, r *http.Request) {
	entries, err := db.GetNews(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get news entries: %s", err), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, entries)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to execute template: %s", err), http.StatusInternalServerError)
		return
	}
}
