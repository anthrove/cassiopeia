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
	"github.com/anthrove/identity/pkg/object"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func (is IdentityService) GetCasbinEnforcer(ctx context.Context, tenantID string, modelID string, adapterID string) (*casbin.Enforcer, error) {
	dbModel, err := is.FindModel(ctx, tenantID, modelID)

	if err != nil {
		return nil, err
	}

	adapter, err := is.FindAdapter(ctx, tenantID, adapterID)
	if err != nil {
		return nil, err
	}

	var casAdapter *gormadapter.Adapter

	if adapter.ExternalDB {
		casAdapter, err = gormadapter.NewAdapter(adapter.Driver, "mysql_username:mysql_password@tcp(127.0.0.1:3306)/", adapter.TableName)
	} else {
		if is.db.Dialector.Name() == "sqlite" {
			db, err := gorm.Open(sqlite.Open(adapter.TableName + ".db"))

			if err != nil {
				return nil, err
			}

			casAdapter, err = gormadapter.NewAdapterByDBUseTableName(db, "", adapter.TableName)
		} else {
			casAdapter, err = gormadapter.NewAdapterByDBUseTableName(is.db, "", adapter.TableName)
		}
	}

	casModel, err := model.NewModelFromString(dbModel.Model)

	if err != nil {
		return nil, err
	}

	casEnforcer, err := casbin.NewEnforcer(casModel, casAdapter)

	if err != nil {
		return nil, err
	}

	return casEnforcer, nil
}

func (is IdentityService) SyncCasbinPermissions(ctx context.Context, tenantID string, enforcerID string) error {
	enforcer, err := is.FindEnforcer(ctx, tenantID, enforcerID)

	if err != nil {
		return err
	}

	permissions, err := is.FindPermissionsByEnforcer(ctx, tenantID, enforcerID)

	if err != nil {
		return err
	}

	policies := make([][]string, 0)
	for _, permission := range permissions {
		policy, err := is.propergatePermissionToPolicies(ctx, permission)

		if err != nil {
			return err
		}

		policies = append(policies, policy...)
	}

	casbinEnforcer, err := is.GetCasbinEnforcer(ctx, tenantID, enforcer.ModelID, enforcer.AdapterID)

	if err != nil {
		return err
	}

	casbinEnforcer.ClearPolicy()
	for _, policy := range policies {
		_, err = casbinEnforcer.AddPolicy(policy)

		if err != nil {
			return err
		}
	}

	if _, groupDefExists := casbinEnforcer.GetModel()["g"]; groupDefExists {
		if _, groupDefGroupExists := casbinEnforcer.GetModel()["g"]["g"]; groupDefGroupExists {
			users, err := is.FindUsers(ctx, tenantID, object.MaxPagination)

			if err != nil {
				return err
			}

			for _, user := range users {
				if len(user.Groups) == 0 {
					continue
				}

				var userGroups []string

				for _, group := range user.Groups {
					userGroups = append(userGroups, group.ID)
				}

				_, err = casbinEnforcer.AddRolesForUser(user.ID, userGroups)

				if err != nil {
					return err
				}
			}
		}
	}

	return casbinEnforcer.SavePolicy()
}

func (is IdentityService) propergatePermissionToPolicies(ctx context.Context, permission object.Permission) ([][]string, error) {
	rules := crossJoin(append(permission.Users, permission.Groups...), permission.V1, permission.V2, permission.V3, permission.V4, permission.V5)

	return rules, nil
}

func (is IdentityService) propergateGroupToUsers(ctx context.Context, tenantID string, groupID string) ([]string, []string, error) {
	_, err := is.FindGroup(ctx, tenantID, groupID)

	if err != nil {
		// TODO check if group not exists error and return empty slice without error
		return nil, nil, err
	}

	users, err := is.FindUsersInGroup(ctx, tenantID, groupID)

	if err != nil {
		return nil, nil, err
	}

	userIDs := make([]string, 0, len(users))
	groupIDs := []string{groupID}
	for _, user := range users {
		userIDs = append(userIDs, user.ID)
	}

	groups, err := is.FindGroupsByParentID(ctx, tenantID, groupID)

	if err != nil {
		return nil, nil, err
	}

	for _, group := range groups {
		uIDs, gIDs, err := is.propergateGroupToUsers(ctx, tenantID, group.ID)

		if err != nil {
			return nil, nil, err
		}

		userIDs = append(userIDs, uIDs...)
		groupIDs = append(groupIDs, gIDs...)
	}

	return userIDs, groupIDs, nil
}

func crossJoin(subIDs []string, v1s []string, v2s []string, v3s []string, v4s []string, v5s []string) [][]string {
	if len(v1s) == 0 {
		v1s = []string{""}
	}

	if len(v2s) == 0 {
		v2s = []string{""}
	}

	if len(v3s) == 0 {
		v3s = []string{""}
	}

	if len(v4s) == 0 {
		v4s = []string{""}
	}

	if len(v5s) == 0 {
		v5s = []string{""}
	}

	items := make([][]string, 0, len(subIDs)*len(v1s)*len(v2s)*len(v3s)*len(v4s)*len(v5s))
	// 28138796 ns / op
	for _, subID := range subIDs {
		for _, v1 := range v1s {
			for _, v2 := range v2s {
				for _, v3 := range v3s {
					for _, v4 := range v4s {
						for _, v5 := range v5s {
							items = append(items, []string{subID, v1, v2, v3, v4, v5})
						}
					}
				}
			}
		}
	}

	return items
}
