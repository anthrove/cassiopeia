package crypto

import (
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

func Test_argon2IDHasher_HashPassword(t *testing.T) {
	type fields struct {
		memory      uint32
		iterations  uint32
		parallelism uint8
		saltLength  uint32
		keyLength   uint32
	}
	type args struct {
		password string
		salt     string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Valid password and salt",
			fields:  fields{memory: 65536, iterations: 3, parallelism: 2, saltLength: 16, keyLength: 32},
			args:    args{password: "password123", salt: "somesalt"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := argon2IDHasher{
				memory:      tt.fields.memory,
				iterations:  tt.fields.iterations,
				parallelism: tt.fields.parallelism,
				saltLength:  tt.fields.saltLength,
				keyLength:   tt.fields.keyLength,
			}
			got, err := a.HashPassword(tt.args.password, tt.args.salt)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.HasPrefix(got, "$argon2id$") {
				t.Errorf("HashPassword() got = %v, want prefix $argon2id$", got)
			}
		})
	}
}

func Test_argon2IDHasher_ComparePassword(t *testing.T) {
	type fields struct {
		memory      uint32
		iterations  uint32
		parallelism uint8
		saltLength  uint32
		keyLength   uint32
	}
	type args struct {
		password   string
		decodeHash string
		salt       string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "Correct password",
			fields:  fields{memory: 65536, iterations: 3, parallelism: 2, saltLength: 16, keyLength: 32},
			args:    args{password: "password123", decodeHash: "$argon2id$v=19$m=65536,t=3,p=2$c29tZXNhbHQ$eXdjYT7Y/ugNjS9HuJfjPSJ155z8+XPJMqK8smK12Z4", salt: "somesalt"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "Incorrect password",
			fields:  fields{memory: 65536, iterations: 3, parallelism: 2, saltLength: 16, keyLength: 32},
			args:    args{password: "wrongpassword", decodeHash: "$argon2id$v=19$m=65536,t=3,p=2$c29tZXNhbHQ$eXdjYT7Y/ugNjS9HuJfjPSJ155z8+XPJMqK8smK12Z4", salt: "somesalt"},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := argon2IDHasher{
				memory:      tt.fields.memory,
				iterations:  tt.fields.iterations,
				parallelism: tt.fields.parallelism,
				saltLength:  tt.fields.saltLength,
				keyLength:   tt.fields.keyLength,
			}
			got, err := a.ComparePassword(tt.args.password, tt.args.decodeHash, tt.args.salt)
			if (err != nil) != tt.wantErr {
				t.Errorf("ComparePassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ComparePassword() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_argon2IDHasher_decodeHash(t *testing.T) {
	type fields struct {
		memory      uint32
		iterations  uint32
		parallelism uint8
		saltLength  uint32
		keyLength   uint32
	}

	tests := []struct {
		name            string
		fields          fields
		encodedHash     string
		wantMemory      uint32
		wantIterations  uint32
		wantParallelism uint8
		wantSalt        []byte
		wantHash        []byte
		wantErr         bool
	}{
		{
			name:            "Valid hash",
			fields:          fields{memory: 65536, iterations: 3, parallelism: 2, saltLength: 16, keyLength: 32},
			encodedHash:     "$argon2id$v=19$m=65536,t=3,p=2$c29tZXNhbHQ$eXdjYT7Y/ugNjS9HuJfjPSJ155z8+XPJMqK8smK12Z4",
			wantMemory:      65536,
			wantIterations:  3,
			wantParallelism: 2,
			wantSalt:        []byte("somesalt"),
			wantHash:        []byte{121, 119, 99, 97, 62, 216, 254, 232, 13, 141, 47, 71, 184, 151, 227, 61, 34, 117, 231, 156, 252, 249, 115, 201, 50, 162, 188, 178, 98, 181, 217, 158},
			wantErr:         false,
		},
		{
			name:            "Invalid hash format",
			fields:          fields{memory: 65536, iterations: 3, parallelism: 2, saltLength: 16, keyLength: 32},
			encodedHash:     "$argon2id$v=19$m=65536,t=3,p=2$c29tZXNhbHQ",
			wantMemory:      0,
			wantIterations:  0,
			wantParallelism: 0,
			wantSalt:        nil,
			wantHash:        nil,
			wantErr:         true,
		},
		{
			name:            "Invalid version",
			fields:          fields{memory: 65536, iterations: 3, parallelism: 2, saltLength: 16, keyLength: 32},
			encodedHash:     "$argon2id$v=18$m=65536,t=3,p=2$c29tZXNhbHQ$eXdjYT7Y/ugNjS9HuJfjPSJ155z8+XPJMqK8smK12Z4",
			wantMemory:      0,
			wantIterations:  0,
			wantParallelism: 0,
			wantSalt:        nil,
			wantHash:        nil,
			wantErr:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := argon2IDHasher{
				memory:      tt.fields.memory,
				iterations:  tt.fields.iterations,
				parallelism: tt.fields.parallelism,
				saltLength:  tt.fields.saltLength,
				keyLength:   tt.fields.keyLength,
			}
			gotMemory, gotIterations, gotParallelism, gotSalt, gotHash, err := a.decodeHash(tt.encodedHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotMemory != tt.wantMemory {
				t.Errorf("decodeHash() gotMemory = %v, want %v", gotMemory, tt.wantMemory)
			}
			if gotIterations != tt.wantIterations {
				t.Errorf("decodeHash() gotIterations = %v, want %v", gotIterations, tt.wantIterations)
			}
			if gotParallelism != tt.wantParallelism {
				t.Errorf("decodeHash() gotParallelism = %v, want %v", gotParallelism, tt.wantParallelism)
			}
			if !reflect.DeepEqual(gotSalt, tt.wantSalt) {
				t.Errorf("decodeHash() gotSalt = %v, want %v", gotSalt, tt.wantSalt)
			}
			if !reflect.DeepEqual(gotHash, tt.wantHash) {
				t.Errorf("decodeHash() gotHash = %v, want %v", gotHash, tt.wantHash)
			}
		})
	}
}

func Test_argon2IDHasher_Speed(t *testing.T) {
	hasher := argon2IDHasher{
		memory:      16 * 1024,
		iterations:  3,
		parallelism: 2,
		saltLength:  16,
		keyLength:   32,
	}

	password := "password123"
	salt := "somesalt"
	testAmount := 1000

	// Measure hashing speed
	startHash := time.Now()
	for i := 0; i < testAmount; i++ {
		_, err := hasher.HashPassword(password, salt)
		if err != nil {
			t.Fatalf("HashPassword() error = %v", err)
		}
	}
	durationHash := time.Since(startHash)

	// Generate a single hash for comparison
	hash, err := hasher.HashPassword(password, salt)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	// Measure comparison speed
	startCompare := time.Now()
	for i := 0; i < testAmount; i++ {
		_, err := hasher.ComparePassword(password, hash, salt)
		if err != nil {
			t.Fatalf("ComparePassword() error = %v", err)
		}
	}
	durationCompare := time.Since(startCompare)

	// Calculate average times
	averageHashTime := durationHash / time.Duration(testAmount)
	averageCompareTime := durationCompare / time.Duration(testAmount)

	t.Logf("Time taken for %d hashings: %v", testAmount, durationHash)
	t.Logf("Time taken for %d comparisons: %v", testAmount, durationCompare)
	t.Logf("Average time per hash: %v", averageHashTime)
	t.Logf("Average time per comparison: %v", averageCompareTime)
}

func Test_argon2IDHasher_SpeedParallel(t *testing.T) {
	hasher := argon2IDHasher{
		memory:      16 * 1024,
		iterations:  3,
		parallelism: 2,
		saltLength:  16,
		keyLength:   32,
	}

	password := "password123"
	salt := "somesalt"
	testAmount := 1000

	var wg sync.WaitGroup

	// Measure hashing speed in parallel
	startHash := time.Now()
	for i := 0; i < testAmount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := hasher.HashPassword(password, salt)
			if err != nil {
				t.Errorf("HashPassword() error = %v", err)
			}
		}()
	}
	wg.Wait()
	durationHash := time.Since(startHash)

	// Generate a single hash for comparison
	hash, err := hasher.HashPassword(password, salt)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	// Measure comparison speed in parallel
	startCompare := time.Now()
	for i := 0; i < testAmount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := hasher.ComparePassword(password, hash, salt)
			if err != nil {
				t.Errorf("ComparePassword() error = %v", err)
			}
		}()
	}
	wg.Wait()
	durationCompare := time.Since(startCompare)

	// Calculate average times
	averageHashTime := durationHash / time.Duration(testAmount)
	averageCompareTime := durationCompare / time.Duration(testAmount)

	t.Logf("Time taken for %d parallel hashings: %v", testAmount, durationHash)
	t.Logf("Time taken for %d parallel comparisons: %v", testAmount, durationCompare)
	t.Logf("Average time per parallel hash: %v", averageHashTime)
	t.Logf("Average time per parallel comparison: %v", averageCompareTime)
}
