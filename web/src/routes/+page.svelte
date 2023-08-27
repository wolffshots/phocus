<script lang="ts">
    import { writable } from 'svelte/store';
    import type InverterData from '../../types/InverterData';
    import { camelCaseToWords } from '$lib/stringUtils';
    import Config from '$lib/components/Config.svelte';
    import { apiUrl } from '$lib/stores/configStore';
    import { webSocketStore } from '$lib/stores/webSocketStore';

    const data = writable<InverterData | null>(null);
    $: {
        if ($apiUrl && typeof window !== 'undefined') {
            const existingWebSocket = $webSocketStore;
            if (existingWebSocket) {
                existingWebSocket.close(); // Close the previous WebSocket instance
            }
            const newWebSocket = new WebSocket($apiUrl);
            newWebSocket.addEventListener("message", (message: MessageEvent) => {
                const messageData: InverterData = JSON.parse(message.data);
                data.set(messageData);
            });
            webSocketStore.set(newWebSocket);
        }
    }
</script>

<style>
    table {
        border-collapse: collapse;
        width: 100%;
    }
    
    th, td {
        padding: 8px;
        text-align: left;
        border-bottom: 1px solid #ddd;
    }
    
    tr:nth-child(even) {
        background-color: #f2f2f2;
    }
</style>

<Config/>

<table>
    <tr>
        <th>Name</th>
        <th>Value</th>
    </tr>
    {#if $data}
        {#each Object.entries($data) as [key, value]}
            <tr>
                <td>{camelCaseToWords(key)}</td>
                <td>
                    {#if typeof value === 'object' && value !== null}
                        <table>
                            {#each Object.entries(value) as [subKey, subValue]}
                                <tr>
                                    <td>{camelCaseToWords(subKey)}</td>
                                    <td>{subValue}</td>
                                </tr>
                            {/each}
                        </table>
                    {:else}
                        {value}
                    {/if}
                </td>
            </tr>
        {/each}
    {:else}
        <tr>
            <td colspan="2">Loading...</td>
        </tr>
    {/if}
</table>
