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

package cdn

import (
	"encoding/csv"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

func HandleLegacyResources(w http.ResponseWriter, r *http.Request) {
	file := strings.TrimPrefix(r.URL.Path, "/client/resources/")

	// object list
	if file == "" {
		files, err := getFiles(filepath.Join("public/resources", file))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "text/csv")

		cw := csv.NewWriter(w)

		for _, f := range files {
			cw.Write([]string{f.Key, strconv.Itoa(f.Size), strconv.Itoa(int(f.Modified.UnixMilli()))})
		}

		cw.Flush()

		return
	}

	http.Redirect(w, r, "//cdn.betablock.net/resources/"+file, http.StatusMovedPermanently)
}
