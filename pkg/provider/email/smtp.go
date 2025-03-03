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

package email

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/anthrove/identity/pkg/object"
	"github.com/go-playground/validator/v10"
	"net"
	"net/mail"
	"net/smtp"
)

type smtpConfiguration struct {
	Host     string `json:"host" validate:"required,max=100"`
	Port     int    `json:"port" validate:"required"`
	SSL      bool   `json:"ssl"`
	Username string `json:"username" validate:"required,max=100"`
	Password string `json:"password" validate:"required,max=100"`
	From     string `json:"from" validate:"required,max=100,email"`
	FromName string `json:"from_name" validate:"required,max=100"`
}

type smtpProvider struct {
	provider object.Provider
}

func newSMTPProvider(provider object.Provider) (Provider, error) {
	return smtpProvider{provider: provider}, nil
}

func (s smtpProvider) GetConfigurationFields() []object.ProviderConfigurationField {
	return []object.ProviderConfigurationField{
		{
			FieldKey:  "host",
			FieldType: "text",
		},
		{
			FieldKey:  "port",
			FieldType: "int",
		},
		{
			FieldKey:  "ssl",
			FieldType: "bool",
		},
		{
			FieldKey:  "username",
			FieldType: "text",
		},
		{
			FieldKey:  "password",
			FieldType: "secret",
		},
		{
			FieldKey:  "from",
			FieldType: "text",
		},
		{
			FieldKey:  "from_name",
			FieldType: "text",
		},
	}
}

func (s smtpProvider) ValidateConfigurationFields() error {
	smtpConfig := smtpConfiguration{}

	err := json.Unmarshal(s.provider.Parameter, &smtpConfig)
	if err != nil {
		return err
	}

	// use a single instance of Validate, it caches struct info
	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(smtpConfig)
	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return errors.Join(fmt.Errorf("problem while validating create tenant data"), validateErrs)
		}
	}

	return nil
}

func (s smtpProvider) SendMail(toAddress string, subject string, body string) error {
	smtpConfig := smtpConfiguration{}

	err := json.Unmarshal(s.provider.Parameter, &smtpConfig)
	if err != nil {
		return err
	}

	from := mail.Address{Name: smtpConfig.FromName, Address: smtpConfig.From}
	to := mail.Address{Address: toAddress}

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subject

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	servername := fmt.Sprintf("%s:%d", smtpConfig.Host, smtpConfig.Port)
	host, _, _ := net.SplitHostPort(servername)

	auth := smtp.PlainAuth("", smtpConfig.Username, smtpConfig.Password, host)

	var tlsConfig *tls.Config
	if smtpConfig.SSL {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         host,
		}
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", servername, tlsConfig)
	if err != nil {
		return err
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	defer c.Quit()

	// Auth
	if err = c.Auth(auth); err != nil {
		return err
	}

	// To && From
	if err = c.Mail(from.Address); err != nil {
		return err
	}

	if err = c.Rcpt(toAddress); err != nil {
		return err
	}

	// Data
	w, err := c.Data()
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = w.Write([]byte(message))
	if err != nil {
		return err
	}

	return nil
}
