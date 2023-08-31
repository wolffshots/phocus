// $lib/stores/configStore.ts

import { writable } from "svelte/store";

export const apiUrl = writable('ws://localhost:8080/last-ws');
