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
	"sync"
)

var sessions map[string]map[string]any
var sessionMutex sync.RWMutex

func init() {
	sessions = make(map[string]map[string]any)
}

func (is IdentityService) UpdateSession(ctx context.Context, sessionID string, session map[string]any) {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	sessions[sessionID] = session
}

func (is IdentityService) KillSession(ctx context.Context, sessionID string) {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	delete(sessions, sessionID)
}

func (is IdentityService) FindSession(ctx context.Context, sessionID string) map[string]any {
	sessionMutex.RLock()
	defer sessionMutex.RUnlock()

	if session, ok := sessions[sessionID]; ok {
		return session
	}

	return make(map[string]any)
}
