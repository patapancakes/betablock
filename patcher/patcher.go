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

package patcher

import (
	"archive/zip"
	"io"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/icholy/replace"
)

type Patcher struct {
	zip *zip.Reader
}

func New(zr *zip.Reader) *Patcher {
	return &Patcher{zip: zr}
}

func (p *Patcher) Write(out io.Writer) error {
	zw := zip.NewWriter(out)
	defer zw.Close()

	for _, f := range p.zip.File {
		// directories are automatically created
		if f.FileInfo().IsDir() {
			continue
		}

		fr, err := f.Open()
		if err != nil {
			return err
		}

		defer fr.Close()

		var body io.Reader = fr

		switch {
		case strings.HasPrefix(f.Name, "META-INF/"):
			// don't include signature files
			if slices.Contains([]string{".dsa", ".rsa", ".sf"}, strings.ToLower(filepath.Ext(f.Name))) {
				continue
			}

			body = replace.Chain(body, replace.Regexp(regexp.MustCompile("SHA1-Digest: (.*)"), nil))
		case filepath.Ext(f.Name) == ".class":
			body = replace.Chain(body,
				replace.String("minecraft.net", "betablock.net"),       // replace minecraft.net
				replace.String("s3.amazonaws.com", "s3.betablock.net"), // replace s3.amazonaws.com
				replace.Bytes(append([]byte{0x01, 0x00, 0x09}, []byte("minecraft")...), append([]byte{0x01, 0x00, 0x09}, []byte("betablock")...)), // replace directory name
			)
		}

		fw, err := zw.Create(f.Name)
		if err != nil {
			return err
		}

		_, err = io.Copy(fw, body)
		if err != nil {
			return err
		}
	}

	return nil
}
