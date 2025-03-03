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

package util

import (
	"crypto/rand"
	"encoding/base64"
	"math/big"
	"strconv"
)

func RandomSaltString(size int) (string, error) {
	salt := make([]byte, size)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(salt), nil
}

func RandomNumber(amount int) int {
	var numberStr string
	for i := 0; i < amount; i++ {
		digit, _ := rand.Int(rand.Reader, big.NewInt(10)) // generates a random digit between 0 and 9
		numberStr += digit.String()
	}
	number, _ := strconv.Atoi(numberStr)
	return number
}
