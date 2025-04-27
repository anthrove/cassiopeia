// Advanced API aware fetch

import { browser } from '$app/environment';
import { NetworkError } from '$lib/errors/NetworkError';
import { get, writable } from 'svelte/store';

const APISERVER = 'http://localhost:8080/api/'

interface RequestOptions extends RequestInit {
    body?: any; // Allow any type, e.g., object, string, FormData, etc.
}

export function useApiFetch(event?: import('@sveltejs/kit').RequestEvent) {
    return async <T>(endpoint: string, options: RequestOptions = {}): Promise<APIResponse<T>> => {
        let token: string | undefined;

        if (browser) {
            // Client: get token from cookie
            const match = document.cookie.match(/(^| )token=([^;]+)/);
            token = match?.[2];
        } else if (event) {
            // Server: get token from event cookies
            token = event.cookies.get('token');
        }

        const headers = new Headers(options.headers || {});
        if (token) headers.set('Authorization', `Bearer ${token}`);
        headers.set('Content-Type', 'application/json');

        const http_target = endpoint.startsWith('http') ? endpoint : `${APISERVER}${endpoint}`

        const isObject = typeof options.body === 'object' && !(options.body instanceof FormData);
        if (isObject) {
            options.body = JSON.stringify(options.body)
        }

        const res = await fetch(http_target, {
            ...options,
            headers,
        });

        if (!res.ok) {
            console.error(`API Error ${res.status}: ${res.statusText}`);
            // Optionally throw or return a custom error here
            let body = ""
            try {
                body = await res.json()
            } catch { 
                return {
                    data: null,
                    error: `${res.status}: ${res.statusText}. No body`
                }
            }
            throw new NetworkError(res.statusText, res.status, body)
        }
        try{
            return await res.json()
        }catch{
            return {
                data: null,
                error: 'No body'
            }
        }
        //return responseBody.data
    };
}
