<script lang="ts">
    import { writable } from 'svelte/store';
    import type InverterData from '../../types/InverterData';
    import { camelCaseToWords } from '$lib/stringUtils';
    import Config from '$lib/components/Config.svelte';
    import { apiUrl } from '$lib/stores/configStore';
    import { webSocketStore } from '$lib/stores/webSocketStore';

    import { onMount, onDestroy } from 'svelte';

    let lastUpdateTime: number | null = null;
    let elapsedTime = 0;
    let lastChecksum1: string | null | undefined = null; // Store the last checksum for InverterNumber x
    let resetCount = 0; // Number of resets
    let totalResetTime = 0; // Total time elapsed during resets
    let firstResetOccurred = false; // Flag to track the first reset

    let start = writable<boolean>(false);

    const data1 = writable<InverterData | null>(null);
    const data2 = writable<InverterData | null>(null);

    let fields = writable<string[]>();

    function getValue(data: InverterData | null, field: string) {
        let value = data ? data[field] : null;
        if (typeof value === 'object' && value !== null) {
            return "?"; // Placeholder for nested object
        }
        return value !== undefined ? value : "?"; // Placeholder for missing data
    }

    $: {
        if ($data1){
            fields.set(Object.keys($data1));
        }else if($data2){
            fields.set(Object.keys($data2));
        }
    }

    $: {
        if ($apiUrl && typeof window !== 'undefined') {
            const existingWebSocket = $webSocketStore;
            if (existingWebSocket) {
                existingWebSocket.close(); // Close the previous WebSocket instance
            }
            if ($start) {
                const newWebSocket = new WebSocket($apiUrl);
                newWebSocket.addEventListener("message", (message: MessageEvent) => {
                    const messageData: InverterData = JSON.parse(message.data);
                    if (messageData.InverterNumber == 1){
                        data1.set(messageData);
                    }else if (messageData.InverterNumber == 2){
                        data2.set(messageData);
                    }else{
                        console.log("some other data arrived");
                    }
                });
                webSocketStore.set(newWebSocket);
            }
        }
    }

    let interval: NodeJS.Timeout | null = null;

    onMount(() => {
        startInterval(); // Start the interval initially
    });

    $: {
        if ($start && !interval) {
            startInterval(); // Start the interval if start becomes true
        } else if (!$start && interval) {
            clearInterval(interval); // Clear the interval if start becomes false
            interval = null;
        }
    }

    const startInterval = () => {
        lastUpdateTime = Date.now();
        interval = setInterval(updateTimer, 10); // Store the interval reference

        onDestroy(() => {
            if (interval) {
                clearInterval(interval);
            }
        });
    };

    const updateTimer = () => {
        if (lastUpdateTime) {
            const currentTime = Date.now();
            const timeSinceLastReset = currentTime - lastUpdateTime;

            // Check if the data changed and reset occurred
            if (lastChecksum1 !== $data1?.Checksum) {
                lastUpdateTime = currentTime;
                lastChecksum1 = $data1?.Checksum;

                // Exclude the first reset and those that occur in less than 3 seconds from the average calculation
                if (firstResetOccurred && $data1?.Checksum !== null && timeSinceLastReset >= 2000) {
                    resetCount += 1;
                    totalResetTime += timeSinceLastReset;
                } else {
                    firstResetOccurred = true;
                }
            }

            elapsedTime = Math.round((currentTime - lastUpdateTime));
        }
    };
</script>

<Config/>

<button on:click={()=>{start.set(!$start)}}>Play/Pause</button>

<div class="time-elapsed">
    <p>Time since reset: {elapsedTime}ms</p>
    <p>Average time since reset: {resetCount > 0 ? Math.floor(totalResetTime / resetCount) : '?'}ms</p>
    <p>Number of resets: {resetCount}</p>
</div>

<div>
    <table>
        <tr>
            <th>Field</th>
            <th>Inverter 1</th>
            <th>Inverter 2</th>
        </tr>
        {#if $data1 || $data2}
            {#each $fields as field}
                <tr>
                    <td>{camelCaseToWords(field)}</td>
                    <td>{getValue($data1, field)}</td>
                    <td>{getValue($data2, field)}</td>
                </tr>
            {/each}
        {:else}
            <tr>
                <td colspan="3">Loading...</td>
            </tr>
        {/if}
    </table>
</div>

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

    .time-elapsed {
        font-size: 14px;
        color: #999;
    }
</style>
