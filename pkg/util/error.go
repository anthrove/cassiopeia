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
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

// ConvertValidationError converts a validator.ValidationErrors object into a single error message.
// It extracts the invalid field names from the validation errors and concatenates them into a single string.
//
// Parameters:
//   - err: a validator.ValidationErrors object containing the validation errors.
//
// Returns:
//   - An error object with a formatted message listing the invalid fields.
func ConvertValidationError(err validator.ValidationErrors) error {
	fields := make([]string, 0, len(err))

	for _, e := range err {
		if e == nil {
			continue
		}

		fields = append(fields, e.StructNamespace())
	}

	return fmt.Errorf("the following fields are not valid: %s", strings.Join(fields, ", "))
}
