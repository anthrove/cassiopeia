import { writable, type Writable } from 'svelte/store';
import { browser } from '$app/environment';

// Helper functions for cookies
function getCookie(name: string): string | null {
    const match = document.cookie.match(new RegExp('(^| )' + name + '=([^;]+)'));
    return match ? decodeURIComponent(match[2]) : null;
}

function setCookie(name: string, value: string, days = 365) {
    const expires = new Date(Date.now() + days * 864e5).toUTCString();
    document.cookie = `${name}=${encodeURIComponent(value)}; expires=${expires}; path=/`;
}

// This will hold all the writable stores per key
let contextReference: { [key: string]: Writable<any> } = {};

// Load initial cookie state
const COOKIE_NAME = 'contextStore'; // Cookie where the JSON state is saved
let initialState: Record<string, any> = {};

if (browser) {
    const cookie = getCookie(COOKIE_NAME);
    if (cookie) {
        try {
            initialState = JSON.parse(cookie);
        } catch {
            console.error('Failed to parse contextStore cookie');
        }
    }
}

function create<T = any>(parameter: string): Writable<T> {
    if (contextReference[parameter]) {
        return contextReference[parameter] as Writable<T>;
    }

    const initial = initialState[parameter] ?? null;
    const store = writable<T>(initial);

    // On change, update cookie
    store.subscribe((value) => {
        initialState[parameter] = value;
        try {
            setCookie(COOKIE_NAME, JSON.stringify(initialState));
        } catch (e) {
            console.error('Failed to stringify contextStore for cookie', e);
        }
    });

    contextReference[parameter] = store;
    return store;
}

export { create };
