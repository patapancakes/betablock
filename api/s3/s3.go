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

package s3

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"text/template"
	"time"

	"github.com/patapancakes/betablock/patcher"

	"github.com/patapancakes/betablock/db"
)

var t = template.Must(template.New("index.xml").ParseFiles("templates/s3/index.xml"))

type Index struct {
	Name  string
	Files []IndexFile
}

type IndexFile struct {
	Path     string
	Modified time.Time
	Hash     string
	Size     int
}

func Handle(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join("public/", r.URL.Path)

	s, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// directory listing
	if s.IsDir() {
		if !slices.Contains([]string{"MinecraftDownload", "MinecraftResources"}, s.Name()) {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		var index Index
		index.Name = s.Name()
		index.Files, err = getFiles(path, path)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "text/xml")

		err = t.Execute(w, index)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}

	var of io.Reader

	switch r.URL.Path {
	case "/MinecraftDownload/minecraft.jar": // handle version selection and patching
		ticket, err := hex.DecodeString(r.URL.Query().Get("ticket"))
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to decode ticket: %s", err), http.StatusBadRequest)
			return
		}

		username, err := db.GetUsernameFromTicket(ticket)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to validate ticket: %s", err), http.StatusBadRequest)
			return
		}

		if r.URL.Query().Get("user") != username {
			http.Error(w, "username mismatch", http.StatusUnauthorized)
			return
		}

		version, err := db.GetUserClientVersion(username)
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
		f, err := os.Open(path)
		if err != nil {
			if os.IsNotExist(err) {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		defer f.Close()

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

func getFiles(base string, path string) ([]IndexFile, error) {
	var files []IndexFile

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, e := range entries {
		if e.IsDir() {
			fs, err := getFiles(base, filepath.Join(path, e.Name()))
			if err != nil {
				return nil, err
			}

			files = append(files, fs...)
			continue
		}

		fpath := filepath.Join(path, e.Name())

		f, err := os.Open(fpath)
		if err != nil {
			return nil, err
		}

		stat, err := f.Stat()
		if err != nil {
			return nil, err
		}

		hash := md5.New()
		_, err = io.Copy(hash, f)
		if err != nil {
			return nil, err
		}

		rel, err := filepath.Rel(base, fpath)
		if err != nil {
			return nil, err
		}

		files = append(files, IndexFile{
			Path:     rel,
			Modified: stat.ModTime(),
			Hash:     hex.EncodeToString(hash.Sum(nil)),
			Size:     int(stat.Size()),
		})
	}

	return files, nil
}
