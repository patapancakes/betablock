package s3

import (
	"encoding/csv"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func HandleLegacyResources(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join("public/MinecraftResources", strings.TrimPrefix(r.URL.Path, "/resources/"))

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
		var index Index
		index.Files, err = getFiles(path, path)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "text/csv")

		cw := csv.NewWriter(w)

		defer cw.Flush()

		for _, f := range index.Files {
			cw.Write([]string{f.Path, strconv.Itoa(f.Size), strconv.Itoa(int(f.Modified.UnixMilli()))})
		}

		return
	}

	// normal file download
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

	io.Copy(w, f)
}

func HandleLegacyCloak(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open(filepath.Join("public/MinecraftCloaks", r.URL.Query().Get("user")+".png"))
	if err != nil {
		if os.IsNotExist(err) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer f.Close()

	io.Copy(w, f)
}
