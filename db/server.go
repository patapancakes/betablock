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

func SetUserServerID(ctx context.Context, username string, sid []byte) error {
	_, err := conn.ExecContext(ctx, "REPLACE INTO players (username, server) VALUES (?, ?)", username, sid)
	if err != nil {
		return err
	}

	return nil
}

func GetUserServerID(ctx context.Context, username string) ([]byte, error) {
	var sid []byte
	err := conn.QueryRowContext(ctx, "SELECT server FROM players WHERE username = ? AND issued > DATE_SUB(UTC_TIMESTAMP(), INTERVAL 1 MINUTE)", username).Scan(&sid)
	if err != nil {
		return nil, err
	}

	return sid, nil
}

func DeleteUserServerID(ctx context.Context, username string) error {
	_, err := conn.ExecContext(ctx, "DELETE FROM players WHERE username = ?", username)
	if err != nil {
		return err
	}

	return nil
}
