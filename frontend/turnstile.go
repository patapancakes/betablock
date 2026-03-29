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

package frontend

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

type TurnstileRequest struct {
	Secret   string `json:"secret"`
	Response string `json:"response"`
	Remote   string `json:"remote"`
}

type TurnstileResponse struct {
	// this is all we care about
	Success bool `json:"success"`
}

func verifyTurnstile(r *http.Request) (bool, error) {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(TurnstileRequest{
		Secret:   os.Getenv("TS_SECRET_KEY"),
		Response: r.FormValue("cf-turnstile-response"),
		Remote:   r.Header.Get("X-Forwarded-For"),
	})
	if err != nil {
		return false, err
	}

	resp, err := http.Post("https://challenges.cloudflare.com/turnstile/v0/siteverify", "application/json", buf)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	var tr TurnstileResponse
	err = json.NewDecoder(resp.Body).Decode(&tr)
	if err != nil {
		return false, err
	}

	return tr.Success, nil
}
