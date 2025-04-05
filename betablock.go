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
	"fmt"
	"log"
	"net/http"

	"github.com/patapancakes/betablock/api/login"
	"github.com/patapancakes/betablock/api/s3"
	"github.com/patapancakes/betablock/db"
	"github.com/patapancakes/betablock/frontend"

	"github.com/patapancakes/betablock/api/game"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// webserver related
	ip := flag.String("ip", "127.0.0.1", "ip to listen on")
	port := flag.Int("port", 80, "port to listen on")

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

	// game
	http.HandleFunc("/game/joinserver.jsp", game.JoinServer)
	http.HandleFunc("/game/checkserver.jsp", game.CheckServer)

	// action
	http.HandleFunc("/register.jsp", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
	})

	http.HandleFunc("/register", frontend.Register)
	http.HandleFunc("/setskin", frontend.SetCosmetic)
	http.HandleFunc("/setcape", frontend.SetCosmetic)
	http.HandleFunc("/setversion", frontend.SetVersion)

	// login
	http.HandleFunc("login.betablock.net/", login.Handle)

	// s3
	http.HandleFunc("s3.betablock.net/", s3.Handle)

	// http stuff
	http.ListenAndServe(fmt.Sprintf("%s:%d", *ip, *port), nil)
}
