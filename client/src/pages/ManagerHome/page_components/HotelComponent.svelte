<script>
    import hotel_door from '../../../resources/single-door.svg';
    import used_room_logo from '../../../resources/woman.svg'
    import room_logo from '../../../resources/window.svg';
    import { push } from 'svelte-spa-router';
    import { onMount } from 'svelte';
    import { server_name } from '../../../server_info';
    
    export let session_key = "";
    export let used_rooms = 0;
    export let rooms_selected = () => {};
    
    let rooms = [];

    onMount(() => {
        const headers = new Headers();
        headers.set("X-sk", session_key);

        fetch(`${server_name}/reservations?guest=*`, {method: 'GET', headers: headers})
            .then(promise => {
                if (promise.ok) {
                    promise.json().then(response => {
                        rooms = response;
                    });
                }
            })
    });

    $: if (rooms.length > 0) {
        let used_room_counter = 0;
        rooms.forEach(r => used_room_counter += r.used_by >= 0 ? 1 : 0);
        used_rooms = used_room_counter;
    }
</script>

<style>

    #hotel-component {
        width: 70%;
        height: 84vh;
        background: var(--color-gradiant);
        border-radius: 15px 15px 0 0;
        box-shadow: -5px 10px 26px 14px rgba(0, 0, 0, 0.2);
    }

    #rooms-container {
        box-sizing: border-box;
        overflow-y: auto;
        display: grid;
        height: 60vh;
        max-height: 60vh;
        padding-top: 5%;
        grid-template-columns: repeat(3, 1fr);
        column-gap: 0vw;
        row-gap: 5vw;
        -ms-overflow-style: none;  /* IE and Edge */
        scrollbar-width: none;  /* Firefox */
    }

    #rooms-container::-webkit-scrollbar {
        display: none;
    }

    .room-container{
        display: flex;
        justify-content: center;
    }

    :global(.room-container svg) {
        cursor: pointer;
        width: 8.5vw;
        fill: white;
    }

    #door-container {
        display: flex;
        height: 24vh;
        justify-content: center;
        align-items: flex-end;
    }

    :global(#door-container svg) {
        width: 8.5vw;
        fill: white;
    } 
</style>

<div id="hotel-component">
    <div id="rooms-container">
        {#each rooms as room, h }
            <div id="room-{room.number}" on:click={() => rooms_selected(rooms[h])} class="room-container">
                {#if room.used_by >= 0}
                    {@html used_room_logo}
                {:else}
                    {@html room_logo}
                {/if}
            </div>
        {/each}
    </div>
    <div id="door-container">
        {@html hotel_door}
    </div>
</div>