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

import "context"

func InsertSession(ctx context.Context, username string, session []byte) error {
	_, err := conn.ExecContext(ctx, "REPLACE INTO sessions (username, session) VALUES (?, ?)", username, session)
	if err != nil {
		return err
	}

	return nil
}

func GetUsernameFromSession(ctx context.Context, session []byte) (string, error) {
	var username string
	err := conn.QueryRowContext(ctx, "SELECT username FROM sessions WHERE session = ? AND issued > DATE_SUB(UTC_TIMESTAMP(), INTERVAL 1 DAY)", session).Scan(&username)
	if err != nil {
		return "", err
	}

	return username, nil
}
