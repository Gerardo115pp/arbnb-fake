<script>
    import { ROOMS_CLASSES } from '../../../enums/rooms';
    import { server_name } from '../../../server_info';

    export let session_key;
    export let show_modal = false;
    export let room_data = {};

    let user_data = {}

    const handleBakcgroundClick = e => {
        if (e.target === e.currentTarget) {
            show_modal = !show_modal;
        }
    }

    const requestUserData = () => {
        if (room_data.used_by >= 0) {
            const headers = new Headers();
            headers.set("X-sk", session_key);

            const request = new Request(`${server_name}/user?guest=${room_data.used_by}`, {method: 'GET', headers: headers});
            fetch(request)
                .then(promise => {
                    if(promise.ok) {
                        promise.json().then(response => user_data = response);
                    }
                })

        }
    }
    $: if (room_data.used_by !== undefined && room_data.used_by >= 0) {
        console.log(room_data)
        requestUserData();
    }
</script>

<style>
    #room-details-modal-background {
        position: fixed;
        display: flex;
        width: 100%;
        height: 100vh;
        background: rgba(0, 0, 0, 0.3);
        justify-content: center;
        align-items: center;
    }

    #rd-modal {
        width: 60%;
        height: 30vh;
        border-radius: 10px;
        background-color: white;
        padding: 1vh 1vw;
        border: 3px solid var(--theme-color);
        user-select: none;
    }

    #room-title {
        display: flex;
        height: 10vh;
        font-size: 1.4rem;
        justify-content: space-around;
        align-items: center;
        border-bottom: 2px dotted var(--theme-color);
    }

    #room-numbert {
        text-transform: uppercase;
        font-size: 2.7rem;
    }

    #user-details {
        display: flex;
        flex-direction: column;
        justify-content: space-between;
        align-items: center;
    }

</style>

<div on:click={handleBakcgroundClick} style="display: {show_modal ? "flex" : "none"};" id="room-details-modal-background">
    <div id="rd-modal">
        <div id="room-title">
            <div id="room-numbert">room {room_data.number}</div>
            <div id="room-classt">{ROOMS_CLASSES[room_data.class]} class room</div>
        </div>
        <div id="room-state">
            <span class="state-title">Room State: {user_data.user === undefined ? "Free" : "Busy" }</span>
            <div id="user-details">
                {#if user_data.user !== undefined}
                    <span class="user-data">name: {user_data.user.name}</span>
                    <span class="user-data">phone: {user_data.user.phone}</span>
                    <span class="user-data">email: {user_data.user.email}</span>
                {/if}
            </div>
        </div>
    </div>
</div>