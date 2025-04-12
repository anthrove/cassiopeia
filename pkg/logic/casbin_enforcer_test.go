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

package logic

import (
	"strconv"
	"testing"
)

// Helper function to generate test data
func generateTestData(prefix string, n int) []string {
	data := make([]string, n)
	for i := 0; i < n; i++ {
		data[i] = prefix + strconv.Itoa(i)
	}
	return data
}

func TestCrossJoinFull(t *testing.T) {
	subIDs := generateTestData("sub", 2)
	v1s := generateTestData("v1", 2)
	v2s := generateTestData("v2", 2)
	v3s := generateTestData("v3", 2)
	v4s := generateTestData("v4", 2)
	v5s := generateTestData("v5", 2)

	result := crossJoin(subIDs, v1s, v2s, v3s, v4s, v5s)

	expectedLength := len(subIDs) * len(v1s) * len(v2s) * len(v3s) * len(v4s) * len(v5s)
	if len(result) != expectedLength {
		t.Errorf("Expected length %d, but got %d", expectedLength, len(result))
	}

	// Optional: Add more detailed checks for specific combinations if needed
}

func TestCrossJoinMissingData(t *testing.T) {
	subIDs := generateTestData("sub", 2)
	v1s := generateTestData("v1", 2)
	v2s := generateTestData("v2", 2)
	v3s := generateTestData("v3", 0)
	v4s := generateTestData("v4", 0)
	v5s := generateTestData("v5", 0)

	result := crossJoin(subIDs, v1s, v2s, v3s, v4s, v5s)

	expectedLength := len(subIDs) * len(v1s) * len(v2s)
	if len(result) != expectedLength {
		t.Errorf("Expected length %d, but got %d", expectedLength, len(result))
	}

	// Optional: Add more detailed checks for specific combinations if needed
}

func BenchmarkCrossJoin(b *testing.B) {
	subIDs := generateTestData("sub", 1000)
	v1s := generateTestData("v1", 5)
	v2s := generateTestData("v2", 10)
	v3s := generateTestData("v3", 2)
	v4s := generateTestData("v4", 2)
	v5s := generateTestData("v5", 2)

	for i := 0; i < b.N; i++ {
		_ = crossJoin(subIDs, v1s, v2s, v3s, v4s, v5s)
	}
}
