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

package db

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// account
func InsertAccount(username string, password string) error {
	digest, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = conn.Exec("INSERT INTO accounts (username, password) VALUES (?, ?)", username, digest)
	if err != nil {
		return err
	}

	return nil
}

func DeleteAccount(username string) error {
	_, err := conn.Exec("DELETE FROM accounts WHERE username = ?", username)
	if err != nil {
		return err
	}

	return nil
}

func ValidatePassword(username string, password string) error {
	var stored []byte
	err := conn.QueryRow("SELECT password FROM accounts WHERE username = ?", username).Scan(&stored)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(stored, []byte(password))
	if err != nil {
		return err
	}

	return nil
}

func GetCanonicalUsername(username string) (string, error) {
	var canonical string
	err := conn.QueryRow("SELECT username FROM accounts WHERE username = ?", username).Scan(&canonical)
	if err != nil {
		return "", err
	}

	return canonical, nil
}

// session
func InsertSession(username string, session []byte) error {
	_, err := conn.Exec("REPLACE INTO sessions (username, session) VALUES (?, ?)", username, session)
	if err != nil {
		return err
	}

	return nil
}

func GetUsernameFromSession(session []byte) (string, error) {
	var username string
	err := conn.QueryRow("SELECT username FROM sessions WHERE session = ?", session).Scan(&username)
	if err != nil {
		return "", err
	}

	return username, nil
}

// ticket
func InsertTicket(username string, ticket []byte) error {
	_, err := conn.Exec("REPLACE INTO tickets (username, ticket) VALUES (?, ?)", username, ticket)
	if err != nil {
		return err
	}

	return nil
}

func GetUsernameFromTicket(ticket []byte) (string, error) {
	var username string
	err := conn.QueryRow("SELECT username FROM tickets WHERE ticket = ?", ticket).Scan(&username)
	if err != nil {
		return "", err
	}

	return username, nil
}

// server id
func SetUserServerID(username string, sid []byte) error {
	_, err := conn.Exec("REPLACE INTO players (username, server) VALUES (?, ?)", username, sid)
	if err != nil {
		return err
	}

	return nil
}

func GetUserServerID(username string) ([]byte, error) {
	var sid []byte
	err := conn.QueryRow("SELECT server FROM players WHERE username = ?", username).Scan(&sid)
	if err != nil {
		return nil, err
	}

	return sid, nil
}

// version
func GetUserClientVersion(username string) (string, error) {
	var version string
	err := conn.QueryRow("SELECT version FROM versions WHERE username = ?", username).Scan(&version)
	if err != nil {
		return "", err
	}

	return version, nil
}

func GetUserClientVersionChanged(username string) (time.Time, error) {
	var changed time.Time
	err := conn.QueryRow("SELECT changed FROM versions WHERE username = ?", username).Scan(&changed)
	if err != nil {
		return time.Now(), err
	}

	return changed, nil
}

func SetUserClientVersion(username string, version string) error {
	_, err := conn.Exec("REPLACE INTO versions (username, version) VALUES (?, ?)", username, version)
	if err != nil {
		return err
	}

	return nil
}
