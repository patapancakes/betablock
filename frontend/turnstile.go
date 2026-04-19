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
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type TurnstileResponse struct {
	// this is all we care about
	Success bool `json:"success"`
}

func verifyTurnstile(r *http.Request) (bool, error) {
	if r.FormValue("cf-turnstile-response") == "" {
		return false, fmt.Errorf("missing cf-turnstile-response")
	}

	resp, err := http.PostForm("https://challenges.cloudflare.com/turnstile/v0/siteverify", url.Values{
		"secret":   {os.Getenv("TURNSTILE_KEY")},
		"response": {r.FormValue("cf-turnstile-response")},
		"remote":   {r.Header.Get("X-Forwarded-For")},
	})
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
