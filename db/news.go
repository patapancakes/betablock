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

type NewsEntry struct {
	Title  string
	Body   string
	Posted time.Time
}

func GetNews(ctx context.Context) ([]NewsEntry, error) {
	rows, err := conn.QueryContext(ctx, "SELECT title, body, posted FROM (SELECT title, body, posted, (DAYOFYEAR(posted) - 304 + 366) % 366 AS pos FROM news) t WHERE pos <= (DAYOFYEAR(NOW()) - 304 + 366) % 366 ORDER BY pos DESC")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var entries []NewsEntry
	for rows.Next() {
		var entry NewsEntry
		err = rows.Scan(&entry.Title, &entry.Body, &entry.Posted)
		if err != nil {
			return nil, err
		}

		entries = append(entries, entry)
	}

	return entries, nil
}
