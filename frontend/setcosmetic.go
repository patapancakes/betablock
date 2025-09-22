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
	"fmt"
	"image/png"
	"net/http"
	"os"
	"path/filepath"
)

func SetCosmetic(w http.ResponseWriter, r *http.Request) {
	cape := r.URL.Path == "/setcape"

	ad := ActionData{Header: "Set Skin", Page: "setskin"}
	if cape {
		ad.Header = "Set Cape"
		ad.Page = "setcape"
	}

	username, err := UsernameFromRequest(r)
	if err != nil {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

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

	// parse form data
	err = r.ParseMultipartForm(maxUploadSize)
	if err != nil {
		Error(w, ad, "An error occured while parsing your request")
		return
	}

	// decode and validate image
	f, fh, err := r.FormFile("image")
	if err != nil {
		Error(w, ad, "An error occured while reading the image")
		return
	}

	defer f.Close()

	if fh.Size > maxUploadSize {
		Error(w, ad, "The image is too large")
		return
	}
	if fh.Header.Get("Content-Type") != "image/png" {
		Error(w, ad, "The image is the wrong type")
		return
	}

	image, err := png.Decode(f)
	if err != nil {
		Error(w, ad, "The image couldn't be decoded")
		return
	}

	dim := image.Bounds()
	if dim.Dx() > 64 || dim.Dy() > 64 {
		Error(w, ad, "The image dimensions are too large")
		return
	}

	dir := "MinecraftSkins"
	if cape {
		dir = "MinecraftCloaks"
	}

	// open dest file for writing
	sf, err := os.OpenFile(filepath.Join("public", dir, username+".png"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		Error(w, ad, "An error occured while creating the image")
		return
	}

	defer sf.Close()

	err = png.Encode(sf, image)
	if err != nil {
		Error(w, ad, "An error occured while encoding the image")
		return
	}

	ad.Success = true

	// write page
	err = t.Execute(w, ad)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to execute template: %s", err), http.StatusInternalServerError)
		return
	}
}
