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

package main

import (
	"embed"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/patapancakes/betablock/api"
	"github.com/patapancakes/betablock/cdn"
	"github.com/patapancakes/betablock/db"
	"github.com/patapancakes/betablock/frontend"
	"github.com/patapancakes/betablock/news"

	_ "github.com/go-sql-driver/mysql"
)

//go:embed frontend/assets
var frontendAssetsFS embed.FS

func main() {
	// init database
	err := db.Init(os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_PROTO"), os.Getenv("DB_ADDR"), os.Getenv("DB_NAME"))
	if err != nil {
		log.Fatalf("error in database init: %s", err)
	}

	// frontend
	http.HandleFunc("/", frontend.About)
	http.HandleFunc("/download", frontend.Download)
	http.HandleFunc("/register", frontend.Register)
	http.HandleFunc("/login", frontend.Login)
	http.HandleFunc("/logout", frontend.Logout)
	http.HandleFunc("/setskin", frontend.SetCosmetic)
	http.HandleFunc("/setcape", frontend.SetCosmetic)
	http.HandleFunc("/setversion", frontend.SetVersion)
	http.HandleFunc("/changepw", frontend.ChangePW)

	assets, err := fs.Sub(frontendAssetsFS, "frontend")
	if err != nil {
		log.Fatalf("failed to create sub fs: %s", err)
	}

	http.Handle("GET /assets/", http.FileServerFS(assets))

	// launcher
	http.HandleFunc("api.betablock.net/launcher/login", api.Login)

	// server
	http.HandleFunc("GET api.betablock.net/server/checkserver", api.CheckServer)

	// client
	http.HandleFunc("GET api.betablock.net/client/joinserver", api.JoinServer)
	http.HandleFunc("GET api.betablock.net/client/session", api.Session)
	http.HandleFunc("GET api.betablock.net/client/resources/", cdn.HandleLegacyResources)
	http.HandleFunc("GET api.betablock.net/client/cloak", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "//cdn.betablock.net/MinecraftCloaks/"+r.URL.Query().Get("user")+".png", http.StatusMovedPermanently)
	})

	// cdn
	http.HandleFunc("cdn.betablock.net/", cdn.Handle)

	// news
	http.HandleFunc("GET news.betablock.net/", news.Handle)
	http.Handle("GET news.betablock.net/assets/", http.StripPrefix("/assets/", http.FileServerFS(news.AssetsFS)))

	httpProto := os.Getenv("HTTP_PROTO")
	httpAddr := os.Getenv("HTTP_ADDR")

	// http stuff
	if httpProto == "unix" {
		err = os.Remove(httpAddr)
		if err != nil && !os.IsNotExist(err) {
			log.Fatalf("failed to delete unix socket: %s", err)
		}
	}

	l, err := net.Listen(httpProto, httpAddr)
	if err != nil {
		log.Fatalf("failed to create web server listener: %s", err)
	}

	defer l.Close()

	if httpProto == "unix" {
		err = os.Chmod(httpAddr, 0777)
		if err != nil {
			log.Fatalf("failed to set unix socket permissions: %s", err)
		}
	}

	http.Serve(l, nil)
}
