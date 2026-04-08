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

	"golang.org/x/crypto/bcrypt"
)

func InsertAccount(ctx context.Context, username string, password string) error {
	digest, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = conn.ExecContext(ctx, "INSERT INTO accounts (username, password) VALUES (?, ?)", username, digest)
	if err != nil {
		return err
	}

	return nil
}

func DeleteAccount(ctx context.Context, username string) error {
	_, err := conn.ExecContext(ctx, "DELETE FROM accounts WHERE username = ?", username)
	if err != nil {
		return err
	}

	return nil
}

func ValidatePassword(ctx context.Context, username string, password string) error {
	var stored []byte
	err := conn.QueryRowContext(ctx, "SELECT password FROM accounts WHERE username = ?", username).Scan(&stored)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(stored, []byte(password))
	if err != nil {
		return err
	}

	return nil
}

func UpdatePassword(ctx context.Context, username string, password string) error {
	digest, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = conn.ExecContext(ctx, "UPDATE accounts SET password = ? WHERE username = ?", digest, username)
	if err != nil {
		return err
	}

	return nil
}

func GetCanonicalUsername(ctx context.Context, username string) (string, error) {
	var canonical string
	err := conn.QueryRowContext(ctx, "SELECT username FROM accounts WHERE username = ?", username).Scan(&canonical)
	if err != nil {
		return "", err
	}

	return canonical, nil
}
