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
	"crypto/sha256"
	"hash"
	"reflect"
	"sync"
	"testing"
)

func Test_pbkdf2Hasher_HashPassword(t *testing.T) {
	type fields struct {
		Digest     func() hash.Hash
		Iterations int
		KeyLen     int
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
				Digest:     sha256.New,
				Iterations: 10000,
				KeyLen:     32,
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
				Digest:     sha256.New,
				Iterations: 10000,
				KeyLen:     32,
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
				Digest:     sha256.New,
				Iterations: 10000,
				KeyLen:     32,
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
			p := pbkdf2Hasher{
				Digest:     tt.fields.Digest,
				Iterations: tt.fields.Iterations,
				KeyLen:     tt.fields.KeyLen,
			}
			_, err := p.HashPassword(tt.args.password, tt.args.salt)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_pbkdf2Hasher_ComparePassword(t *testing.T) {
	type fields struct {
		Digest     func() hash.Hash
		Iterations int
		KeyLen     int
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
				Digest:     sha256.New,
				Iterations: 10000,
				KeyLen:     32,
			},
			args: args{
				password:       "password123",
				hashedPassword: "$pbkdf2$iter=10000$keylen=32$U29tZVNhbHQ$9pWMSDVXvi0VVYRbk/+ZCIzyQpAEI+i9yNQPNYZB9Gc",
				salt:           "SomeSalt",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid password",
			fields: fields{
				Digest:     sha256.New,
				Iterations: 10000,
				KeyLen:     32,
			},
			args: args{
				password:       "wrongpassword",
				hashedPassword: "$pbkdf2$iter=10000$keylen=32$U29tZVNhbHRz$HehMDXY3FiZPjarf/q9TsN0a8IJIahtm5jc+wgMaBXI",
				salt:           "SomeSalt",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "invalid hash format",
			fields: fields{
				Digest:     sha256.New,
				Iterations: 10000,
				KeyLen:     32,
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
			p := pbkdf2Hasher{
				Digest:     tt.fields.Digest,
				Iterations: tt.fields.Iterations,
				KeyLen:     tt.fields.KeyLen,
			}
			got, err := p.ComparePassword(tt.args.password, tt.args.hashedPassword, tt.args.salt)
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

func Test_pbkdf2Hasher_decodeHash(t *testing.T) {
	type fields struct {
		Digest     func() hash.Hash
		Iterations int
		KeyLen     int
	}
	type args struct {
		encodedHash string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantIterations int
		wantKeyLen     int
		wantDigest     func() hash.Hash
		wantSalt       []byte
		wantHash       []byte
		wantErr        bool
	}{
		{
			name: "valid hash",
			fields: fields{
				Digest:     sha256.New,
				Iterations: 10000,
				KeyLen:     32,
			},
			args: args{
				encodedHash: "$pbkdf2$iter=10000$keylen=32$U29tZVNhbHQ$9pWMSDVXvi0VVYRbk/+ZCIzyQpAEI+i9yNQPNYZB9Gc",
			},
			wantIterations: 10000,
			wantKeyLen:     32,
			wantDigest:     sha256.New,
			wantSalt:       []byte("SomeSalt"),
			wantHash:       []byte{246, 149, 140, 72, 53, 87, 190, 45, 21, 85, 132, 91, 147, 255, 153, 8, 140, 242, 66, 144, 4, 35, 232, 189, 200, 212, 15, 53, 134, 65, 244, 103},
			wantErr:        false,
		},
		{
			name: "invalid hash format",
			fields: fields{
				Digest:     sha256.New,
				Iterations: 10000,
				KeyLen:     32,
			},
			args: args{
				encodedHash: "invalidhashformat",
			},
			wantIterations: 0,
			wantKeyLen:     0,
			wantDigest:     nil,
			wantSalt:       nil,
			wantHash:       nil,
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := pbkdf2Hasher{
				Iterations: tt.fields.Iterations,
				KeyLen:     tt.fields.KeyLen,
			}
			gotIterations, gotKeyLen, gotSalt, gotHash, err := p.decodeHash(tt.args.encodedHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotIterations != tt.wantIterations {
				t.Errorf("decodeHash() gotIterations = %v, want %v", gotIterations, tt.wantIterations)
			}
			if gotKeyLen != tt.wantKeyLen {
				t.Errorf("decodeHash() gotKeyLen = %v, want %v", gotKeyLen, tt.wantKeyLen)
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

func Benchmark_pbkdf2Hasher_Speed(b *testing.B) {
	hasher := pbkdf2Hasher{
		Digest:     sha256.New,
		Iterations: 10000,
		KeyLen:     32,
	}

	password := "password123"
	salt := "SomeSalt"

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := hasher.HashPassword(password, salt)
		if err != nil {
			b.Fatalf("HashPassword() error = %v", err)
		}
	}
}

func Benchmark_pbkdf2Hasher_SpeedParallel(b *testing.B) {
	hasher := pbkdf2Hasher{
		Digest:     sha256.New,
		Iterations: 10000,
		KeyLen:     32,
	}

	password := "password123"
	salt := "SomeSalt"

	var wg sync.WaitGroup

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := hasher.HashPassword(password, salt)
			if err != nil {
				b.Errorf("HashPassword() error = %v", err)
			}
		}()
	}
	wg.Wait()
}
