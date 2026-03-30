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
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/patapancakes/betablock/patcher"
)

func Download(w http.ResponseWriter, r *http.Request) {
	ad := ActionData{Header: "Get Betablock", Page: "download"}

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

	if os.Getenv("TS_SITE_KEY") != "" {
		ok, err := verifyTurnstile(r)
		if err != nil {
			Error(w, ad, "Server error")
			return
		}
		if !ok {
			Error(w, ad, "Verification failed")
			return
		}
	}

	const maxSize = 1024 * 1024 * 4 // 4MB

	f, fh, err := r.FormFile("launcher")
	if err != nil {
		Error(w, ad, "No file attached")
		return
	}
	if fh.Size > int64(maxSize) {
		Error(w, ad, "File too large")
	}

	b, err := io.ReadAll(http.MaxBytesReader(w, f, int64(maxSize)))
	if err != nil {
		Error(w, ad, "Malformed file")
		return
	}

	br := bytes.NewReader(b)

	zr, err := zip.NewReader(br, br.Size())
	if err != nil {
		Error(w, ad, "File is not a JAR")
		return
	}

	w.Header().Set("Content-Type", "application/java-archive")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s-patched.jar\"", strings.TrimSuffix(fh.Filename, path.Ext(fh.Filename))))
	err = patcher.New(zr).Write(w)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to write patched zip: %s", err), http.StatusInternalServerError)
		return
	}
}
