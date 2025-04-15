<!--
  - Copyright (C) 2025 Anthrove
  -
  - Licensed under the Apache License, Version 2.0 (the "License");
  - you may not use this file except in compliance with the License.
  - You may obtain a copy of the License at
  -
  -      http://www.apache.org/licenses/LICENSE-2.0
  -
  - Unless required by applicable law or agreed to in writing, software
  - distributed under the License is distributed on an "AS IS" BASIS,
  - WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  - See the License for the specific language governing permissions and
  - limitations under the License.
  -->

<script lang="ts">
    import type {PageProps} from './$types';

    let { data }: PageProps = $props();
    const user = data.data.user;
</script>

<div class="min-h-screen flex items-center justify-center bg-gray-900">
    <div class="bg-gray-800 p-8 rounded-lg shadow-lg w-full max-w-md">
        <h1 class="text-3xl font-bold mb-6 text-center text-white">Profile</h1>
        <div class="text-white">
            <div class="mb-4">
                <label class="block text-sm font-bold mb-2">Username:</label>
                <p class="bg-gray-700 p-2 rounded">{user.username}</p>
            </div>
            <div class="mb-4">
                <label class="block text-sm font-bold mb-2">Display Name:</label>
                <p class="bg-gray-700 p-2 rounded">{user.display_name}</p>
            </div>
            <div class="mb-4">
                <label class="block text-sm font-bold mb-2">Email:</label>
                <p class="bg-gray-700 p-2 rounded">{user.email}</p>
            </div>
            <div class="mb-4">
                <label class="block text-sm font-bold mb-2">Email Verified:</label>
                <p class="bg-gray-700 p-2 rounded">{user.email_verified ? 'Yes' : 'No'}</p>
            </div>
            <div class="mb-4">
                <label class="block text-sm font-bold mb-2">Created At:</label>
                <p class="bg-gray-700 p-2 rounded">{new Date(user.created_at).toLocaleString()}</p>
            </div>
            <div class="mb-4">
                <label class="block text-sm font-bold mb-2">Updated At:</label>
                <p class="bg-gray-700 p-2 rounded">{new Date(user.updated_at).toLocaleString()}</p>
            </div>
        </div>
    </div>
</div>

<style>
    /* Add any additional custom styles here */
</style>
