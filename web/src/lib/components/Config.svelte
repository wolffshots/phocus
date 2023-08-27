<!-- Config.svelte -->

<script lang="ts">
    import { onMount } from 'svelte';
    import { setCookie, getCookie } from '$lib/cookieUtils'; // Create cookieUtils file
	import { apiUrl } from '$lib/stores/configStore';
    
    let editedApiUrl = '';

    // Load the saved configuration from the cookie on component mount
    onMount(() => {
        apiUrl.set(getCookie('apiUrl') || 'ws://localhost:8080/last-ws');
    });

    // Handle the save action when the save button is clicked
    function handleSave() {
        if (editedApiUrl) {
            setCookie('apiUrl', editedApiUrl);
            apiUrl.set(editedApiUrl);
        }
    }
</script>

<input type="text" bind:value={editedApiUrl} placeholder="API URL/Hostname">
<button on:click={handleSave}>Save</button>