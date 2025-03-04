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
	"bytes"
	"testing"
)

func Test_scryptHasher_HashPassword(t *testing.T) {
	type fields struct {
		CostFactor      int
		BlockSize       int
		Parallelization int
		KeyLength       int
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
			name: "valid hash",
			fields: fields{
				CostFactor:      16384,
				BlockSize:       8,
				Parallelization: 1,
				KeyLength:       32,
			},
			args: args{
				password: "password123",
				salt:     "SomeSalt",
			},
			wantErr: false,
		},
		{
			name: "empty password",
			fields: fields{
				CostFactor:      16384,
				BlockSize:       8,
				Parallelization: 1,
				KeyLength:       32,
			},
			args: args{
				password: "",
				salt:     "SomeSalt",
			},
			wantErr: true,
		},
		{
			name: "empty salt",
			fields: fields{
				CostFactor:      16384,
				BlockSize:       8,
				Parallelization: 1,
				KeyLength:       32,
			},
			args: args{
				password: "password123",
				salt:     "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := scryptHasher{
				CostFactor:      tt.fields.CostFactor,
				BlockSize:       tt.fields.BlockSize,
				Parallelization: tt.fields.Parallelization,
				KeyLength:       tt.fields.KeyLength,
			}
			_, err := s.HashPassword(tt.args.password, tt.args.salt)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_scryptHasher_ComparePassword(t *testing.T) {
	type fields struct {
		CostFactor      int
		BlockSize       int
		Parallelization int
		KeyLength       int
	}
	type args struct {
		password       string
		hashedPassword string
		salt           string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "valid password",
			fields: fields{
				CostFactor:      16384,
				BlockSize:       8,
				Parallelization: 1,
				KeyLength:       32,
			},
			args: args{
				password:       "password123",
				hashedPassword: "$scrypt$costFactor=16384$blockSize=8$parallelization=1$keyLength=32$U29tZVNhbHQ$udrWtp+zcLL4CGV9wqC2aMQm3BkMFs5F1OIPkUow6qE",
				salt:           "SomeSalt",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid password",
			fields: fields{
				CostFactor:      16384,
				BlockSize:       8,
				Parallelization: 1,
				KeyLength:       32,
			},
			args: args{
				password:       "wrongpassword",
				hashedPassword: "$scrypt$costFactor=16384$blockSize=8$parallelization=1$keyLength=32$U29tZVNhbHQ$9pWMSDVXvi0VVYRbk/+ZCIzyQpAEI+i9yNQPNYZB9Gc",
				salt:           "SomeSalt",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "invalid hash format",
			fields: fields{
				CostFactor:      16384,
				BlockSize:       8,
				Parallelization: 1,
				KeyLength:       32,
			},
			args: args{
				password:       "password123",
				hashedPassword: "invalidhashformat",
				salt:           "SomeSalt",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := scryptHasher{
				CostFactor:      tt.fields.CostFactor,
				BlockSize:       tt.fields.BlockSize,
				Parallelization: tt.fields.Parallelization,
				KeyLength:       tt.fields.KeyLength,
			}
			got, err := s.ComparePassword(tt.args.password, tt.args.hashedPassword, tt.args.salt)
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

func Test_scryptHasher_decodeHash(t *testing.T) {
	type fields struct {
		CostFactor      int
		BlockSize       int
		Parallelization int
		KeyLength       int
	}
	type args struct {
		encodedHash string
	}
	tests := []struct {
		name                string
		fields              fields
		args                args
		wantCostFactor      int
		wantBlockSize       int
		wantParallelization int
		wantKeyLength       int
		wantSalt            []byte
		wantHash            []byte
		wantErr             bool
	}{
		{
			name: "valid hash",
			fields: fields{
				CostFactor:      16384,
				BlockSize:       8,
				Parallelization: 1,
				KeyLength:       32,
			},
			args: args{
				encodedHash: "$scrypt$costFactor=16384$blockSize=8$parallelization=1$keyLength=32$U29tZVNhbHQ$udrWtp+zcLL4CGV9wqC2aMQm3BkMFs5F1OIPkUow6qE",
			},
			wantCostFactor:      16384,
			wantBlockSize:       8,
			wantParallelization: 1,
			wantKeyLength:       32,
			wantSalt:            []byte("SomeSalt"),
			wantHash:            []byte{185, 218, 214, 182, 159, 179, 112, 178, 248, 8, 101, 125, 194, 160, 182, 104, 196, 38, 220, 25, 12, 22, 206, 69, 212, 226, 15, 145, 74, 48, 234, 161},
			wantErr:             false,
		},
		{
			name: "invalid hash format",
			fields: fields{
				CostFactor:      16384,
				BlockSize:       8,
				Parallelization: 1,
				KeyLength:       32,
			},
			args: args{
				encodedHash: "invalidhashformat",
			},
			wantCostFactor:      0,
			wantBlockSize:       0,
			wantParallelization: 0,
			wantKeyLength:       0,
			wantSalt:            nil,
			wantHash:            nil,
			wantErr:             true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := scryptHasher{
				CostFactor:      tt.fields.CostFactor,
				BlockSize:       tt.fields.BlockSize,
				Parallelization: tt.fields.Parallelization,
				KeyLength:       tt.fields.KeyLength,
			}
			gotCostFactor, gotBlockSize, gotParallelization, gotKeyLength, gotSalt, gotHash, err := s.decodeHash(tt.args.encodedHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCostFactor != tt.wantCostFactor {
				t.Errorf("decodeHash() gotCostFactor = %v, want %v", gotCostFactor, tt.wantCostFactor)
			}
			if gotBlockSize != tt.wantBlockSize {
				t.Errorf("decodeHash() gotBlockSize = %v, want %v", gotBlockSize, tt.wantBlockSize)
			}
			if gotParallelization != tt.wantParallelization {
				t.Errorf("decodeHash() gotParallelization = %v, want %v", gotParallelization, tt.wantParallelization)
			}
			if gotKeyLength != tt.wantKeyLength {
				t.Errorf("decodeHash() gotKeyLength = %v, want %v", gotKeyLength, tt.wantKeyLength)
			}
			if !bytes.Equal(gotSalt, tt.wantSalt) {
				t.Errorf("decodeHash() gotSalt = %v, want %v", gotSalt, tt.wantSalt)
			}
			if !bytes.Equal(gotHash, tt.wantHash) {
				t.Errorf("decodeHash() gotHash = %v, want %v", gotHash, tt.wantHash)
			}
		})
	}
}

func Benchmark_scryptHasher_HashPassword(b *testing.B) {
	hasher := scryptHasher{
		CostFactor:      16384,
		BlockSize:       8,
		Parallelization: 1,
		KeyLength:       32,
	}
	password := "password123"
	salt := "SomeSalt"

	for i := 0; i < b.N; i++ {
		_, err := hasher.HashPassword(password, salt)
		if err != nil {
			b.Fatalf("HashPassword() error = %v", err)
		}
	}
}

func Benchmark_scryptHasher_HashPasswordParallel(b *testing.B) {
	hasher := scryptHasher{
		CostFactor:      16384,
		BlockSize:       8,
		Parallelization: 1,
		KeyLength:       32,
	}
	password := "password123"
	salt := "SomeSalt"

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := hasher.HashPassword(password, salt)
			if err != nil {
				b.Fatalf("HashPassword() error = %v", err)
			}
		}
	})
}
