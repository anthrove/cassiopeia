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
	"errors"
	"fmt"
	"github.com/anthrove/identity/pkg/object"
	"github.com/anthrove/identity/pkg/provider/email"
	"github.com/go-playground/validator/v10"
)

func (is IdentityService) SendMail(ctx context.Context, tenantID string, providerID string, mailData object.SendMailData) error {
	providerObj, err := is.FindProvider(ctx, tenantID, providerID)
	if err != nil {
		return err
	}

	if providerObj.Category != "email" {
		return errors.New("provider category not email")
	}

	provider, err := email.GetEMailProvider(providerObj)
	if err != nil {
		return err
	}

	err = validate.Struct(mailData)
	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return errors.Join(fmt.Errorf("problem while validating send mail data"), validateErrs)
		}
	}

	err = provider.SendMail(mailData.To, mailData.Subject, mailData.Body)
	if err != nil {
		return err
	}

	return nil
}
