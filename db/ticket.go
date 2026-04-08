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

func InsertTicket(ctx context.Context, username string, ticket []byte) error {
	_, err := conn.ExecContext(ctx, "REPLACE INTO tickets (username, ticket) VALUES (?, ?)", username, ticket)
	if err != nil {
		return err
	}

	return nil
}

func GetUsernameFromTicket(ctx context.Context, ticket []byte) (string, error) {
	var username string
	err := conn.QueryRowContext(ctx, "SELECT username FROM tickets WHERE ticket = ? AND issued > DATE_SUB(UTC_TIMESTAMP(), INTERVAL 1 DAY)", ticket).Scan(&username)
	if err != nil {
		return "", err
	}

	return username, nil
}

func DeleteTicket(ctx context.Context, ticket []byte) error {
	_, err := conn.ExecContext(ctx, "DELETE FROM tickets WHERE ticket = ?", ticket)
	if err != nil {
		return err
	}

	return nil
}
