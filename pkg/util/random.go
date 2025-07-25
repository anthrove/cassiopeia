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
	"github.com/sethvargo/go-diceware/diceware"
	"math/big"
	"strconv"
	"strings"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

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

func RandomString(length int) (string, error) {
	result := make([]byte, length)
	for i := range result {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[randomIndex.Int64()]
	}
	return string(result), nil
}

func RandomPassPhrase(length int, concatSymbol string) (string, error) {
	list, err := diceware.Generate(length)
	if err != nil {
		return "", err
	}
	return strings.Join(list, concatSymbol), nil
}
