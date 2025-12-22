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
	"flag"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/patapancakes/betablock/api"
	"github.com/patapancakes/betablock/cdn"
	"github.com/patapancakes/betablock/db"
	"github.com/patapancakes/betablock/frontend"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// webserver related
	proto := flag.String("proto", "tcp", "proto for web server")
	addr := flag.String("addr", "127.0.0.1:80", "address for web server")

	// database related
	dbuser := flag.String("dbuser", "betablock", "database user's name")
	dbpass := flag.String("dbpass", "", "database user's password")
	dbproto := flag.String("dbproto", "tcp", "database connection protocol")
	dbaddr := flag.String("dbaddr", "127.0.0.1:3306", "database address")
	dbname := flag.String("dbname", "betablock", "database name")
	flag.Parse()

	// init database
	err := db.Init(*dbuser, *dbpass, *dbproto, *dbaddr, *dbname)
	if err != nil {
		log.Fatalf("error in database init: %s", err)
	}

	// frontend
	http.Handle("/", http.RedirectHandler("//www.betablock.net/register", http.StatusSeeOther))

	http.HandleFunc("/register", frontend.Register)
	http.HandleFunc("/login", frontend.Login)
	http.HandleFunc("/logout", frontend.Logout)
	http.HandleFunc("/setskin", frontend.SetCosmetic)
	http.HandleFunc("/setcape", frontend.SetCosmetic)
	http.HandleFunc("/setversion", frontend.SetVersion)
	http.HandleFunc("/changepw", frontend.ChangePW)

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

	// http stuff
	if *proto == "unix" {
		err = os.Remove(*addr)
		if err != nil && !os.IsNotExist(err) {
			log.Fatalf("failed to delete unix socket: %s", err)
		}
	}

	l, err := net.Listen(*proto, *addr)
	if err != nil {
		log.Fatalf("failed to create web server listener: %s", err)
	}

	defer l.Close()

	if *proto == "unix" {
		err = os.Chmod(*addr, 0777)
		if err != nil {
			log.Fatalf("failed to set unix socket permissions: %s", err)
		}
	}

	http.Serve(l, nil)
}
