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

	"github.com/patapancakes/betablock/api/login"
	"github.com/patapancakes/betablock/api/s3"
	"github.com/patapancakes/betablock/db"
	"github.com/patapancakes/betablock/frontend"

	"github.com/patapancakes/betablock/api/game"

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

	http.HandleFunc("/register", frontend.Register)
	http.HandleFunc("/login", frontend.Login)
	http.HandleFunc("/logout", frontend.Logout)
	http.HandleFunc("/setskin", frontend.SetCosmetic)
	http.HandleFunc("/setcape", frontend.SetCosmetic)
	http.HandleFunc("/setversion", frontend.SetVersion)
	http.HandleFunc("/changepw", frontend.ChangePW)

	// action
	http.Handle("/", http.RedirectHandler("/register", http.StatusSeeOther))
	http.Handle("/register.jsp", http.RedirectHandler("/register", http.StatusSeeOther))

	// game
	http.HandleFunc("GET /game/joinserver.jsp", game.JoinServer)
	http.HandleFunc("GET /game/checkserver.jsp", game.CheckServer)
	http.HandleFunc("POST /game/getversion.jsp", login.Login) // legacy login

	// login
	http.Handle("GET login.betablock.net/", http.RedirectHandler("https://betablock.net", http.StatusSeeOther))
	http.HandleFunc("POST login.betablock.net/", login.Login)
	http.HandleFunc("GET login.betablock.net/session", login.Session)

	// s3
	http.HandleFunc("s3.betablock.net/", s3.Handle)

	// legacy assets
	http.HandleFunc("GET /resources/", s3.HandleLegacyResources)
	http.HandleFunc("GET /cloak/get.jsp", s3.HandleLegacyCloak)
	http.Handle("GET /skin/", http.StripPrefix("/skin/", http.FileServer(http.Dir("public/MinecraftSkins"))))

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
