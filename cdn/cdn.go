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
	"archive/zip"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strconv"
	"time"

	"github.com/patapancakes/betablock/patcher"

	"github.com/patapancakes/betablock/db"
)

type ListBucketResult struct {
	XMLName  xml.Name `xml:"ListBucketResult"`
	Name     string   `xml:"Name"`
	Contents []Object `xml:"Contents"`
}

type Object struct {
	Key  string `xml:"Key"`
	Size int    `xml:"Size"`

	Hash     string    `xml:"-"`
	Modified time.Time `xml:"-"`
}

func Handle(w http.ResponseWriter, r *http.Request) {
	file := filepath.Join("public", r.URL.Path)

	// object list
	if slices.Contains([]string{"/binaries/", "/resources/"}, r.URL.Path) {
		var err error

		res := ListBucketResult{Name: path.Base(r.URL.Path)}
		res.Contents, err = getFiles(file)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "text/xml")

		err = xml.NewEncoder(w).Encode(res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}

	var of io.Reader

	switch r.URL.Path {
	case "/binaries/minecraft.jar": // handle version selection and patching
		ticket, err := hex.DecodeString(r.URL.Query().Get("ticket"))
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to decode ticket: %s", err), http.StatusBadRequest)
			return
		}

		username, err := db.GetUsernameFromTicket(r.Context(), ticket)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to validate ticket: %s", err), http.StatusBadRequest)
			return
		}

		if r.URL.Query().Get("user") != username {
			http.Error(w, "username mismatch", http.StatusUnauthorized)
			return
		}

		version, err := db.GetUserClientVersion(r.Context(), username)
		if err != nil {
			version = "b1.7.3"
		}

		f, err := os.Open(filepath.Join("clients", version+".jar"))
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to open client for reading: %s", err), http.StatusInternalServerError)
			return
		}

		defer f.Close()

		s, err := f.Stat()
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to stat client: %s", err), http.StatusInternalServerError)
			return
		}

		zr, err := zip.NewReader(f, s.Size())
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to open client zip: %s", err), http.StatusInternalServerError)
			return
		}

		patched := new(bytes.Buffer)
		err = patcher.New(zr).Write(patched)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to patch client: %s", err), http.StatusInternalServerError)
			return
		}

		of = patched
	default: // normal file download
		f, err := os.Open(file)
		if err != nil {
			if os.IsNotExist(err) {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		defer f.Close()

		s, err := f.Stat()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if s.IsDir() {
			w.WriteHeader(http.StatusOK)
			return
		}

		of = f
	}

	b, err := io.ReadAll(of)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to read output data: %s", err), http.StatusInternalServerError)
		return
	}

	etag := fmt.Sprintf("\"%x\"", md5.Sum(b))

	if r.Header.Get("If-None-Match") == etag {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	w.Header().Set("ETag", etag)
	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	w.Write(b)
}

func getFiles(base string) ([]Object, error) {
	var files []Object

	err := filepath.WalkDir(base, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}

		defer f.Close()

		stat, err := f.Stat()
		if err != nil {
			return err
		}

		hash := md5.New()
		_, err = io.Copy(hash, f)
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(base, path)
		if err != nil {
			return err
		}

		files = append(files, Object{
			Key:      rel,
			Size:     int(stat.Size()),
			Modified: stat.ModTime(),
			Hash:     hex.EncodeToString(hash.Sum(nil)),
		})

		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}
