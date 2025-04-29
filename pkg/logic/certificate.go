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
	"context"
	"crypto/elliptic"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/anthrove/identity/pkg/object"
	"github.com/anthrove/identity/pkg/repository"
	"github.com/anthrove/identity/pkg/util"
	"github.com/go-playground/validator/v10"
)

func (is IdentityService) CreateCertificate(ctx context.Context, tenantID string, createCertificate object.CreateCertificate) (object.Certificate, error) {
	dbConn, _ := is.getDBConn(ctx)

	err := validate.Struct(createCertificate)

	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return object.Certificate{}, errors.Join(fmt.Errorf("problem while validating create certificate data"), util.ConvertValidationError(validateErrs))
		}
	}

	certificate := object.Certificate{
		TenantID:    tenantID,
		DisplayName: createCertificate.DisplayName,
		Algo:        createCertificate.Algorithm,
		BitSize:     createCertificate.BitSize,
		ExpiredAt:   createCertificate.ExpiredAt,
	}

	switch createCertificate.Algorithm {
	case "RS256":
		key, keyData, err := util.GenerateRSAKey(certificate.BitSize)
		if err != nil {
			return object.Certificate{}, err
		}

		cert, err := util.GenerateCertificate(key, &key.PublicKey, createCertificate.ExpiredAt, x509.SHA256WithRSA)
		if err != nil {
			return object.Certificate{}, err
		}

		certificate.Certificate = string(cert)
		certificate.PrivateKey = string(keyData)
	case "RS384":
		key, keyData, err := util.GenerateRSAKey(certificate.BitSize)
		if err != nil {
			return object.Certificate{}, err
		}

		cert, err := util.GenerateCertificate(key, &key.PublicKey, createCertificate.ExpiredAt, x509.SHA384WithRSA)
		if err != nil {
			return object.Certificate{}, err
		}

		certificate.Certificate = string(cert)
		certificate.PrivateKey = string(keyData)
	case "RS512":
		key, keyData, err := util.GenerateRSAKey(certificate.BitSize)
		if err != nil {
			return object.Certificate{}, err
		}

		cert, err := util.GenerateCertificate(key, &key.PublicKey, createCertificate.ExpiredAt, x509.SHA512WithRSA)
		if err != nil {
			return object.Certificate{}, err
		}

		certificate.Certificate = string(cert)
		certificate.PrivateKey = string(keyData)
	case "ES256":
		key, keyData, err := util.GenerateECDSAKey(elliptic.P256())
		if err != nil {
			return object.Certificate{}, err
		}

		cert, err := util.GenerateCertificate(key, &key.PublicKey, createCertificate.ExpiredAt, x509.SHA512WithRSA)
		if err != nil {
			return object.Certificate{}, err
		}

		certificate.Certificate = string(cert)
		certificate.PrivateKey = string(keyData)
	case "ES384":
		key, keyData, err := util.GenerateECDSAKey(elliptic.P384())
		if err != nil {
			return object.Certificate{}, err
		}

		cert, err := util.GenerateCertificate(key, &key.PublicKey, createCertificate.ExpiredAt, x509.SHA512WithRSA)
		if err != nil {
			return object.Certificate{}, err
		}

		certificate.Certificate = string(cert)
		certificate.PrivateKey = string(keyData)
	case "ES512":
		key, keyData, err := util.GenerateECDSAKey(elliptic.P521())
		if err != nil {
			return object.Certificate{}, err
		}

		cert, err := util.GenerateCertificate(key, &key.PublicKey, createCertificate.ExpiredAt, x509.SHA512WithRSA)
		if err != nil {
			return object.Certificate{}, err
		}

		certificate.Certificate = string(cert)
		certificate.PrivateKey = string(keyData)
	default:
		return object.Certificate{}, errors.New("given algorithm is not supported")
	}

	return repository.CreateCertificate(ctx, dbConn, tenantID, certificate)
}

func (is IdentityService) UpdateCertificate(ctx context.Context, tenantID string, certificateID string, updateCertificate object.UpdateCertificate) error {
	dbConn, _ := is.getDBConn(ctx)

	if len(tenantID) == 0 {
		return errors.New("tenantID is required")
	}

	err := validate.Struct(updateCertificate)

	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return errors.Join(fmt.Errorf("problem while validating create tenant data"), validateErrs)
		}
	}

	return repository.UpdateCertificate(ctx, dbConn, tenantID, certificateID, updateCertificate)
}

func (is IdentityService) KillCertificate(ctx context.Context, tenantID string, certificateID string) error {
	dbConn, _ := is.getDBConn(ctx)

	return repository.KillCertificate(ctx, dbConn, tenantID, certificateID)
}

func (is IdentityService) FindCertificate(ctx context.Context, tenantID string, certificateID string) (object.Certificate, error) {
	dbConn, _ := is.getDBConn(ctx)

	return repository.FindCertificate(ctx, dbConn, tenantID, certificateID)
}

func (is IdentityService) FindCertificates(ctx context.Context, tenantID string, pagination object.Pagination) ([]object.Certificate, error) {
	dbConn, _ := is.getDBConn(ctx)

	return repository.FindCertificates(ctx, dbConn, tenantID, pagination)
}

func (is IdentityService) FindAllCertificates(ctx context.Context) ([]object.Certificate, error) {
	dbConn, _ := is.getDBConn(ctx)

	return repository.FindAllCertificates(ctx, dbConn)
}
