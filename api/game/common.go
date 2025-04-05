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

package game

import (
	"encoding/hex"
	"errors"
)

var errServerIdTooLong = errors.New("server id is too long")

func getPaddedServerID(sid string) ([]byte, error) {
	if len(sid) > 16 {
		return nil, errServerIdTooLong
	}

	for range 16 - len(sid) {
		sid += "0"
	}

	serverId, err := hex.DecodeString(sid)
	if err != nil {
		return nil, err
	}

	return serverId, nil
}
