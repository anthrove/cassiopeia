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

import "time"

type Group struct {
	ID       string `json:"id" gorm:"primaryKey;type:char(25)"`
	TenantID string `json:"tenant_id"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	ParentGroupID *string `json:"parent_group_id"`
	ParentGroup   *Group  `json:"-"`

	DisplayName string `json:"displayName;type:varchar(100)"`

	Enabled bool `json:"enabled"`
}
