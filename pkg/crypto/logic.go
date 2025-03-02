/*
 * Copyright (C) 2025 Anthrove
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package crypto

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher interface {
	HashPassword(password string, salt string) (string, error)
	ComparePassword(password string, hashedPassword string, salt string) (bool, error)
}

func GetPasswordHasher(passwordType string) (PasswordHasher, error) {
	switch passwordType {
	case "bcrypt":
		return NewBcryptHasher(bcrypt.DefaultCost), nil
	default:
		return nil, fmt.Errorf("unsupported password type: %s", passwordType)
	}
}
