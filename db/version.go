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
	"context"
	"time"
)

func GetUserClientVersion(ctx context.Context, username string) (string, error) {
	var version string
	err := conn.QueryRowContext(ctx, "SELECT version FROM versions WHERE username = ?", username).Scan(&version)
	if err != nil {
		return "", err
	}

	return version, nil
}

func GetUserClientVersionChanged(ctx context.Context, username string) (time.Time, error) {
	var changed time.Time
	err := conn.QueryRowContext(ctx, "SELECT changed FROM versions WHERE username = ?", username).Scan(&changed)
	if err != nil {
		return time.Now(), err
	}

	return changed, nil
}

func SetUserClientVersion(ctx context.Context, username string, version string) error {
	_, err := conn.ExecContext(ctx, "REPLACE INTO versions (username, version) VALUES (?, ?)", username, version)
	if err != nil {
		return err
	}

	return nil
}

func GetRealtimeVersion(ctx context.Context) (string, time.Time, error) {
	var version string
	var released time.Time
	err := conn.QueryRowContext(ctx, "SELECT id, released FROM (SELECT id, released, (DAYOFYEAR(released) - 304 + 366) % 366 AS pos FROM timeline) t WHERE pos <= (DAYOFYEAR(NOW()) - 304 + 366) % 366 ORDER BY pos DESC").Scan(&version, &released)
	if err != nil {
		return "", time.Time{}, err
	}

	return version, released, nil
}
