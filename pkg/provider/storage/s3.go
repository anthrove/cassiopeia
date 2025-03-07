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

package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/anthrove/identity/pkg/object"
	"github.com/aws/aws-sdk-go/service/s3control"
	"github.com/go-playground/validator/v10"
	"github.com/qor/oss/s3"
)

type s3Configuration struct {
	AccessID  string `json:"access_id" validate:"required"`
	AccessKey string `json:"access_key" validate:"required"`
	Bucket    string `json:"bucket" validate:"required"`
	Endpoint  string `json:"endpoint" validate:"required"`
	Region    string `json:"region" validate:"required"`
}

type s3Provider struct {
	s3.Client
	provider object.Provider
}

func newS3Provider(provider object.Provider) (Provider, error) {
	var parameters map[string]string
	err := json.Unmarshal(provider.Parameter, &parameters)
	if err != nil {
		return nil, err
	}

	// Specific for Minio
	s3config := s3.Config{
		AccessID:         parameters["access_id"],
		AccessKey:        parameters["access_key"],
		Region:           parameters["region"],
		Bucket:           parameters["bucket"],
		Endpoint:         parameters["endpoint"],
		S3Endpoint:       parameters["endpoint"],
		ACL:              s3control.BucketCannedACLPublicRead,
		S3ForcePathStyle: true,
	}

	return s3Provider{
		Client:   *s3.New(&s3config),
		provider: provider,
	}, nil
}

func (s s3Provider) GetConfigurationFields() []object.ProviderConfigurationField {
	return []object.ProviderConfigurationField{
		{
			FieldKey:  "access_id",
			FieldType: "text",
		},
		{
			FieldKey:  "access_key",
			FieldType: "text",
		},
		{
			FieldKey:  "region",
			FieldType: "text",
		},
		{
			FieldKey:  "bucket",
			FieldType: "text",
		},
		{
			FieldKey:  "endpoint",
			FieldType: "text",
		},
	}
}

func (s s3Provider) ValidateConfigurationFields(provider object.Provider) error {
	s3Configuration := s3Configuration{}

	// TODO: more validation?

	err := json.Unmarshal(s.provider.Parameter, &s3Configuration)
	if err != nil {
		return err
	}

	// use a single instance of Validate, it caches struct info
	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(s3Configuration)
	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return errors.Join(fmt.Errorf("problem while validating create tenant data"), validateErrs)
		}
	}

	return nil
}

func (s s3Provider) GetBucketName() (string, error) {
	var parameters map[string]string
	err := json.Unmarshal(s.provider.Parameter, &parameters)
	if err != nil {
		return "", err
	}

	return parameters["bucket"], nil

}
