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

import "testing"

func Test_bcryptHasher_ComparePassword(t *testing.T) {
	type fields struct {
		costs int
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
			name: "Normal without Salt",
			fields: fields{
				costs: 10,
			},
			args: args{
				password:       "password",
				hashedPassword: "$2a$10$kTpM/XVZ3p6YTHCCnzL7IOKZKlYqqEsB71VSj3fJZQnQ1.sF33ah6",
				salt:           "",
			},
			want: true,
		},
		{
			name: "Normal with user salt",
			fields: fields{
				costs: 10,
			},
			args: args{
				password:       "password",
				hashedPassword: "$2a$10$Zj2Y8LPM0zXf/VL86/5B2u554W/PO9vYXTiiWAme0E.voBH9iho9q",
				salt:           "123",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := bcryptHasher{
				costs: tt.fields.costs,
			}
			got, err := h.ComparePassword(tt.args.password, tt.args.hashedPassword, tt.args.salt)
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
