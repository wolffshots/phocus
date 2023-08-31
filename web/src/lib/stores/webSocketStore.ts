// $lib/stores/webSocketStore.ts

import { writable } from 'svelte/store';

export const webSocketStore = writable<WebSocket | null>(null);