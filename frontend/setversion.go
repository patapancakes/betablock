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
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/patapancakes/betablock/db"

	"golang.org/x/crypto/bcrypt"
)

func SetVersion(w http.ResponseWriter, r *http.Request) {
	ad := ActionData{Header: "Set Version", Page: "setversion"}
	entries, err := os.ReadDir("clients")
	if err != nil {
		Error(w, ad, "An error occured while getting available client versions")
		return
	}

	var versions []string
	for _, e := range entries {
		if e.Type().IsDir() {
			continue
		}

		if filepath.Ext(e.Name()) != ".jar" {
			continue
		}

		versions = append(versions, strings.TrimSuffix(e.Name(), filepath.Ext(e.Name())))
	}

	ad.Versions = versions

	if r.Method == "GET" {
		err := t.Execute(w, ad)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to execute template: %s", err), http.StatusInternalServerError)
			return
		}

		return
	}

	err = r.ParseForm()
	if err != nil {
		Error(w, ad, "An error occured while parsing your request")
		return
	}

	// validate username and password
	username := r.PostFormValue("username")

	err = db.ValidatePassword(r.Context(), username, r.PostFormValue("password"))
	if err != nil {
		var reason string
		switch err {
		case sql.ErrNoRows:
			reason = "The specified user doesn't exist"
		case bcrypt.ErrMismatchedHashAndPassword:
			reason = "The password is incorrect"
		default:
			reason = "An unknown error occured during account validation"
		}

		Error(w, ad, reason)
		return
	}

	version := r.PostForm.Get("version")
	if !slices.Contains(versions, version) {
		Error(w, ad, "The specified version isn't available")
		return
	}

	err = db.SetUserClientVersion(r.Context(), username, version)
	if err != nil {
		Error(w, ad, "An unknown error occured while setting the version")
		return
	}

	ad.Success = true

	err = t.Execute(w, ad)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to execute template: %s", err), http.StatusInternalServerError)
		return
	}
}
