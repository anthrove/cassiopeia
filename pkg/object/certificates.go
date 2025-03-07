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

package object

import (
	"crypto/x509"
	"encoding/pem"
	"github.com/go-jose/go-jose/v4"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
	"time"
)

type Certificate struct {
	ID       string `json:"id" gorm:"primaryKey;type:char(25)" example:"BsOOg4igppKxYwhAQQrD3GCRZ"`
	TenantID string `json:"tenant_id" maxLength:"25" minLength:"25" example:"BsOOg4igppKxYwhAQQrD3GCRZ"`

	CreatedAt time.Time `json:"created_at" format:"date-time"`
	UpdatedAt time.Time `json:"updated_at" format:"date-time"`

	DisplayName string    `json:"display_name" gorm:"type:varchar(100)" maxLength:"100" example:"Certification Title"`
	Algorithm   string    `json:"algorithm" gorm:"type:varchar(100)" maxLength:"100" example:"RS512"`
	BitSize     int       `json:"bit_size" example:"2048"`
	ExpiredAt   time.Time `json:"expired_at" format:"date-time"`

	Certificate string `json:"certificate"`
	PrivateKey  string `json:"private_key"`

	Applications []Application `json:"-" swaggerignore:"true"`
}

// BeforeCreate is a GORM hook that is called before a new group record is inserted into the database.
// It generates a unique ID for the group if it is not already set.
//
// Parameters:
//   - db: a gorm.DB instance representing the database connection.
//
// Returns:
//   - An error if there is any issue generating the unique ID.
func (base *Certificate) BeforeCreate(db *gorm.DB) error {
	if base.ID == "" {
		id, err := gonanoid.New(25)
		if err != nil {
			return err
		}

		base.ID = id
	}

	return nil
}

type CreateCertificate struct {
	DisplayName string    `json:"display_name" validate:"required,max=100" maxLength:"100" example:"Certificate Title"`
	Algorithm   string    `json:"algorithm" validate:"required" example:"RS512"`
	BitSize     int       `json:"bit_size" example:"2048"`
	ExpiredAt   time.Time `json:"expired_at" format:"date-time"`
}

type UpdateCertificate struct {
	DisplayName string `json:"display_name" validate:"required,max=100" maxLength:"100" example:"Certificate Title"`
}

func (base *Certificate) ToJWK() (jose.JSONWebKey, error) {
	certPemBlock := []byte(base.Certificate)
	certDerBlock, _ := pem.Decode(certPemBlock)
	x509Cert, err := x509.ParseCertificate(certDerBlock.Bytes)
	if err != nil {
		return jose.JSONWebKey{}, err
	}

	return jose.JSONWebKey{
		Key:          x509Cert.PublicKey,
		KeyID:        base.ID,
		Algorithm:    base.Algorithm,
		Use:          "sig",
		Certificates: []*x509.Certificate{x509Cert},
	}, nil
}
