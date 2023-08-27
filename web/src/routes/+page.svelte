<script lang="ts">
    import { onMount } from 'svelte';
    import { writable } from 'svelte/store';
    import type InverterData from '../../types/InverterData';
    import { camelCaseToWords } from '$lib/stringUtils';

    const data = writable<InverterData | null>(null);
	onMount(() => {
        const ws = new WebSocket("ws://localhost:8080/last-ws");
        ws.addEventListener("message", (message: MessageEvent) => {
            const messageData: InverterData = JSON.parse(message.data);
            data.set(messageData);
        });
    })
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
