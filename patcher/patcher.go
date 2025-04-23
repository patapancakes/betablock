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
	"bufio"
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"slices"
	"strings"
)

func NewPatcher(zr *zip.Reader) (io.Reader, error) {
	ozb := new(bytes.Buffer) // output zip buffer

	zw := zip.NewWriter(ozb)

	for _, f := range zr.File {
		// create directories, doesn't seem to get used
		if f.FileInfo().IsDir() {
			_, err := zw.Create(f.Name)
			if err != nil {
				return nil, err
			}
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

		switch {
		case strings.HasPrefix(f.Name, "META-INF/"):
			// don't include signature files
			if slices.Contains([]string{".dsa", ".rsa", ".sf"}, strings.ToLower(filepath.Ext(f.Name))) {
				continue
			}

			scanner := bufio.NewScanner(bytes.NewReader(fb.Bytes()))

			fb.Reset()

			for scanner.Scan() {
				if strings.HasPrefix(scanner.Text(), "SHA1-Digest") {
					continue
				}

				fmt.Fprintln(fb, scanner.Text())
			}
		case filepath.Ext(f.Name) == ".class":
			// replace minecraft.net
			rep := bytes.ReplaceAll(fb.Bytes(), []byte("minecraft.net"), []byte("betablock.net"))

			// replace s3.amazonaws.com
			rep = bytes.ReplaceAll(rep, []byte("s3.amazonaws.com"), []byte("s3.betablock.net"))

			// replace directory name in getWorkingDirectory call
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

	err := zw.Close()
	if err != nil {
		return nil, err
	}

	return ozb, nil
}
