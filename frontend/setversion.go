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
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/patapancakes/betablock/db"
)

var versions, _ = getVersions()

func getVersions() ([]Version, error) {
	f, err := AssetsFS.Open("assets/versions.csv")
	if err != nil {
		return nil, err
	}

	defer f.Close()

	cr := csv.NewReader(f)

	versionTimes := make(map[string]time.Time)
	for {
		records, err := cr.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		if len(records) < 2 {
			continue
		}

		unixTime, err := strconv.Atoi(records[1])
		if err != nil {
			return nil, err
		}

		versionTimes[records[0]] = time.Unix(int64(unixTime), 0)
	}

	entries, err := os.ReadDir("clients")
	if err != nil {
		return nil, err
	}

	var versions []Version
	for _, e := range entries {
		if e.Type().IsDir() {
			continue
		}

		if filepath.Ext(e.Name()) != ".jar" {
			continue
		}

		version := strings.TrimSuffix(e.Name(), filepath.Ext(e.Name()))
		versionTime, _ := versionTimes[version]

		versions = append(versions, Version{Name: version, Time: versionTime})
	}

	slices.SortFunc(versions, func(a, b Version) int {
		if !a.Time.IsZero() && !b.Time.IsZero() {
			return a.Time.Compare(b.Time)
		}

		return strings.Compare(a.Name, b.Name)
	})

	versions = append([]Version{{Name: "realtime"}}, versions...)

	return versions, nil
}

func SetVersion(w http.ResponseWriter, r *http.Request) {
	ad := ActionData{Header: "Set Version", Page: "setversion"}

	ad.Versions = versions

	var err error
	ad.Username, err = UsernameFromRequest(r)
	if err != nil {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/logout", http.StatusSeeOther)
		return
	}

	ad.Version, err = db.GetUserClientVersion(r.Context(), ad.Username)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/logout", http.StatusSeeOther)
		return
	}

	if r.Method == "GET" {
		err := t.Execute(w, ad)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to execute template: %s", err), http.StatusInternalServerError)
			return
		}

		return
	}

	version := r.PostFormValue("version")

	var found bool
	for _, e := range versions {
		if e.Name != version {
			continue
		}

		found = true
		break
	}
	if !found {
		Error(w, ad, "The specified version isn't available")
		return
	}

	err = db.SetUserClientVersion(r.Context(), ad.Username, version)
	if err != nil {
		Error(w, ad, "An unknown error occured while setting the version")
		return
	}

	ad.Success = true
	ad.Version = version

	err = t.Execute(w, ad)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to execute template: %s", err), http.StatusInternalServerError)
		return
	}
}
