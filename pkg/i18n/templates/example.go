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

package templates

const EMailLVerificationTemplateWithLink = `<!DOCTYPE html>
<html>
<head>
    <style>
        body {
            background-color: #f7fafc;
            color: #2d3748;
            font-family: Arial, sans-serif;
        }
        .container {
            max-width: 600px;
            margin: 40px auto;
            padding: 20px;
            background-color: #ffffff;
            border-radius: 8px;
            box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
        }
        .header {
            font-size: 24px;
            font-weight: bold;
            margin-bottom: 16px;
        }
        .content {
            margin-bottom: 16px;
        }
        .button {
            display: inline-block;
            padding: 10px 20px;
            color: #ffffff;
            background-color: #4299e1;
            border-radius: 4px;
            text-decoration: none;
        }
        .button:hover {
            background-color: #3182ce;
        }
        .link {
            margin-top: 16px;
            font-size: 14px;
            color: #4299e1;
            word-break: break-all;
        }
        .footer {
            margin-top: 16px;
            font-size: 12px;
            color: #718096;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1 class="header">Welcome, {{.DisplayName}}!</h1>
        <p class="content">Thank you for registering. Please verify your email address by clicking the link below:</p>
        <a href="{{.VerificationLink}}" class="button">Verify Email</a>
        <p class="content">If the button above does not work, please copy and paste the following link into your browser:</p>
        <p class="link">{{.VerificationLink}}</p>
        <p class="footer">If you did not register for this account, please ignore this email.</p>
    </div>
</body>
</html>`

const EMailLVerificationTemplateWithCode = `<!DOCTYPE html>
<html>
<head>
    <style>
        body {
            background-color: #f7fafc;
            color: #2d3748;
            font-family: Arial, sans-serif;
        }
        .container {
            max-width: 600px;
            margin: 40px auto;
            padding: 20px;
            background-color: #ffffff;
            border-radius: 8px;
            box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
        }
        .header {
            font-size: 24px;
            font-weight: bold;
            margin-bottom: 16px;
        }
        .content {
            margin-bottom: 16px;
        }
        .code {
            display: inline-block;
            padding: 10px 20px;
            color: #ffffff;
            background-color: #4299e1;
            border-radius: 4px;
            font-family: monospace;
            font-size: 18px;
        }
        .footer {
            margin-top: 16px;
            font-size: 12px;
            color: #718096;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1 class="header">Welcome, {{.DisplayName}}!</h1>
        <p class="content">Thank you for registering. Please verify your email address by using the code below:</p>
        <div class="code">{{.VerificationCode}}</div>
        <p class="footer">If you did not register for this account, please ignore this email.</p>
    </div>
</body>
</html>`
