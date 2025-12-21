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
	"encoding/binary"
	"io"
	"net/http"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/icholy/replace"
)

const (
	host = "betablock.net"

	wwwHost   = "www." + host
	loginHost = "login." + host
	s3Host    = "s3." + host
	startHost = "start." + host
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
				// client
				replace.Bytes(strb("https://login.minecraft.net/session?name="), strb("https://"+loginHost+"/session?name=")),
				replace.Bytes(strb("http://www.minecraft.net/game/joinserver.jsp?user="), strb("http://"+wwwHost+"/game/joinserver.jsp?user=")),

				// legacy client
				replace.Bytes(strb("http://www.minecraft.net/resources/"), strb("http://"+wwwHost+"/resources/")),
				replace.Bytes(strb("http://www.minecraft.net/skin/"), strb("http://"+wwwHost+"/skin/")),
				replace.Bytes(strb("http://www.minecraft.net/cloak/get.jsp?user="), strb("http://"+wwwHost+"/cloak/get.jsp?user=")),

				// client resources
				replace.Bytes(strb("http://s3.amazonaws.com/MinecraftSkins/"), strb("http://"+s3Host+"/MinecraftSkins/")),
				replace.Bytes(strb("http://s3.amazonaws.com/MinecraftCloaks/"), strb("http://"+s3Host+"/MinecraftCloaks/")),
				replace.Bytes(strb("http://s3.amazonaws.com/MinecraftResources/"), strb("http://"+s3Host+"/MinecraftResources/")),

				// server
				replace.Bytes(strb("http://www.minecraft.net/game/checkserver.jsp?user="), strb("http://"+wwwHost+"/game/checkserver.jsp?user=")),

				// launcher
				replace.Bytes(strb("http://mcupdate.tumblr.com/"), strb("https://"+startHost+"/")),
				replace.Bytes(strb("https://login.minecraft.net/"), strb("https://"+loginHost+"/")),
				replace.Bytes(strb("http://s3.amazonaws.com/MinecraftDownload/"), strb("http://"+s3Host+"/MinecraftDownload/")),

				// legacy launcher
				replace.Bytes(strb("http://www.minecraft.net/game/getversion.jsp"), strb("http://"+wwwHost+"/game/getversion.jsp")),

				// launcher + client
				replace.Bytes(strb("minecraft"), strb("betablock")), // replace directory name
			)
		case filepath.Base(f.Name) == "minecraft.key":
			resp, err := http.Get("https://" + loginHost + "/session")
			if err != nil {
				continue
			}

			for _, cert := range resp.TLS.PeerCertificates {
				if cert.Subject.CommonName != host && !slices.Contains(cert.DNSNames, host) {
					continue
				}

				body = bytes.NewReader(cert.RawSubjectPublicKeyInfo)
				break
			}
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

func strb(s string) []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.BigEndian, uint8(1))
	binary.Write(buf, binary.BigEndian, uint16(len(s)))
	buf.WriteString(s)

	return buf.Bytes()
}
