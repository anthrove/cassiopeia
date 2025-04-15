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
<script>

    import {goto} from "$app/navigation";

    const urlParams = new URLSearchParams(window.location.search)

    // static variables right now... we need to get it dynamically later
    let tenantID = 'TK8ZPdka-YEGQjJepmPlTgOEm';
    let applicationID = 'OHqVwht2rT_WDt-jYugzUUyWB';

    let requestID = urlParams.get("request_id")
    let email = '';
    let password = '';

    async function handleSubmit() {
        // Handle form submission
        console.log('Email:', email);
        console.log('Password:', password);

        const response = await fetch('/api/v1/tenant/' + tenantID + "/application/" + applicationID + "/login", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                "username": email,
                "password": password,
                "request_id": requestID
            })
        });
        const body = await response.json();
        console.log(body);

        goto(body.data.redirect_uri)
    };
</script>

<div class="min-h-screen flex items-center justify-center bg-gray-900">
    <div class="p-8 rounded-lg shadow-lg w-full max-w-md bg-gray-800">
        <h1 class="text-2xl font-bold mb-6 text-center text-white">Login</h1>
        <form on:submit|preventDefault={() => handleSubmit()}>
            <div class="mb-4">
                <label class="block text-sm font-bold mb-2 text-gray-300" for="email">
                    Email
                </label>
                <input
                        bind:value={email}
                        class="shadow appearance-none border rounded w-full py-2 px-3 leading-tight focus:outline-none focus:shadow-outline bg-gray-700 text-gray-300 border-gray-600"
                        id="email"
                        placeholder="Enter your Username"
                        required
                />
            </div>
            <div class="mb-6">
                <label class="block text-sm font-bold mb-2 text-gray-300" for="password">
                    Password
                </label>
                <input
                        bind:value={password}
                        class="shadow appearance-none border rounded w-full py-2 px-3 mb-3 leading-tight focus:outline-none focus:shadow-outline bg-gray-700 text-gray-300 border-gray-600"
                        id="password"
                        placeholder="Enter your password"
                        required
                        type="password"
                />
            </div>
            <div class="flex items-center justify-between">
                <button
                        class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
                        type="submit"
                >
                    Sign In
                </button>
            </div>
        </form>
    </div>
</div>

<style>
    /* Add any additional custom styles here */
</style>
