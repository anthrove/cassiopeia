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
	"crypto/md5"
	"encoding/hex"
	"io"
)

// HashFileMD5 computes the MD5 hash of the content read from the provided io.Reader.
//
// Parameters:
//   - file: an io.Reader representing the file content to be hashed.
//
// Returns:
//   - A string representing the MD5 hash in hexadecimal format.
//   - An error if there is any issue during hashing.
func HashFileMD5(file io.Reader) (string, error) {
	hasher := md5.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}
	hash := hex.EncodeToString(hasher.Sum(nil))
	return hash, nil
}
