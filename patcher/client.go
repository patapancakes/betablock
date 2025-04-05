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
	"bytes"
	"io"
	"path/filepath"
	"strings"
)

func PatchGameClient(r io.ReaderAt, s int64) (io.Reader, error) {
	zr, err := zip.NewReader(r, s)
	if err != nil {
		return nil, err
	}

	ozb := new(bytes.Buffer) // output zip buffer

	zw := zip.NewWriter(ozb)

	for _, f := range zr.File {
		// don't include META-INF
		if strings.HasPrefix(f.Name, "META-INF") {
			continue
		}

		// create directories, doesn't seem to get used
		if f.FileInfo().IsDir() {
			_, err := zw.Create(f.Name)
			if err != nil {
				return nil, err
			}

			continue
		}

		fr, err := f.Open()
		if err != nil {
			return nil, err
		}

		fb := new(bytes.Buffer)

		_, err = io.Copy(fb, fr)
		if err != nil {
			return nil, err
		}

		// only patch .class files
		if filepath.Ext(f.Name) == ".class" {
			// replace minecraft.net
			rep := bytes.ReplaceAll(fb.Bytes(), []byte("minecraft.net"), []byte("betablock.net"))

			// replace s3.amazonaws.com
			rep = bytes.ReplaceAll(rep, []byte("s3.amazonaws.com"), []byte("s3.betablock.net"))

			// look for getWorkingDirectory
			if bytes.Contains(rep, []byte("user.home")) {
				rep = bytes.ReplaceAll(rep, append([]byte{0x01, 0x00, 0x09}, []byte("minecraft")...), append([]byte{0x01, 0x00, 0x09}, []byte("betablock")...))
			}

			fb.Reset()

			_, err = fb.Write(rep)
			if err != nil {
				return nil, err
			}
		}

		fw, err := zw.Create(f.Name)
		if err != nil {
			return nil, err
		}

		_, err = io.Copy(fw, fb)
		if err != nil {
			return nil, err
		}
	}

	err = zw.Close()
	if err != nil {
		return nil, err
	}

	return ozb, nil
}
