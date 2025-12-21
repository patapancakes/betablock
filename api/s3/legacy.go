package s3

import (
	"encoding/csv"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

func HandleLegacyResources(w http.ResponseWriter, r *http.Request) {
	file := strings.TrimPrefix(r.URL.Path, "/resources/")

	// object list
	if file == "" {
		files, err := getFiles(filepath.Join("public/MinecraftResources", file))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "text/csv")

		cw := csv.NewWriter(w)

		defer cw.Flush()

		for _, f := range files {
			cw.Write([]string{f.Key, strconv.Itoa(f.Size), strconv.Itoa(int(f.Modified.UnixMilli()))})
		}

		return
	}

	http.Redirect(w, r, "//s3.betablock.net/MinecraftResources/"+file, http.StatusMovedPermanently)
}
