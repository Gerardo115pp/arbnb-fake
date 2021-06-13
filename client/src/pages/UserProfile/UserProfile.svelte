<script>
    import Button from '../../components/Button.svelte';
    import { ROOMS_CLASSES } from '../../enums/rooms';
    import { server_name } from "../../server_info";
    import { push } from 'svelte-spa-router';
    import { onMount } from "svelte";

    export let params = {};

    let user_data = {};
    let user_name = "";

    let { sk:session_key } = params;

    onMount(() => requestUserData());

    const cancelReservation = room_data => {
        const headers = new Headers();
        headers.set("X-sk", session_key);

        const form_data = new FormData();
        form_data.append('room', room_data.number);

        const request = new Request(`${server_name}/reservations`, {method: 'DELETE', body:form_data, headers: headers});

        fetch(request)
            .then(promise => {
                if(promise.ok) {
                    requestUserData();
                }
            })
    }

    const requestUserData = () => {
        const headers = new Headers();
        headers.set("X-sk", session_key);

        const request = new Request(`${server_name}/user`, {method: 'GET', headers: headers});

        fetch(request)
            .then(promise => {
                if (promise.ok) {
                    promise.json().then(response => user_data = response);
                }
            })
    }

    const logout = () => {
        const headers = new Headers();
        headers.set('X-sk', session_key);

        const request = new Request(`${server_name}/logout`, { method: 'PATCH', headers: headers });
        fetch(request)
            .then(promise => {
                if(promise.ok) {
                    push("/");
                }
            })
    }



    $: if (user_data.user !== undefined) {
        user_name = user_data.user.username
    }

</script>

<style>
    #user-profile-bar {
        display: flex;
        padding: 1vh 2vw;
        justify-content: space-between;
        align-items: center;
    }

    #username-title {
        background: var(--theme-gradiant);
        font-size: 1.3rem;
        font-weight: bolder;
        color: white;
        padding: .4% 1%;
        border-radius: 5px;
    }

    #logout-btn {
        cursor: pointer;
        font-size: 1.2rem;
    }

    #username-title::before {
        content: '@';
    }

    #main-content {
        display: flex;
        height: 90vh;
    }

    #reservations-container {
        box-sizing: border-box;
        width: 70%;
        display: flex;
        justify-content: space-around;
        align-items: center;
        flex-direction: column;
        padding-top: 5vh;
    }
    
    #reservations {
        overflow-y: auto;
        width: 60%;
        height: 40vh;
        background: var(--theme-gradiant);
        border-radius: 5px;
        padding: 1vh 1vw;
    }

    .reservation-container {
        cursor: pointer;
        box-sizing: border-box;
        display: flex;
        justify-content: space-between;
        padding: 1vh 6vw;
        font-size: 1.2rem;
        color: white;
        border-bottom: 2px dotted white;
        transition: all .3s ease-in;
    }

    .reservation-container:hover {
        background-color: rgba(255, 32, 32, 0.3);
    }

    .reservation-container .room-class {
        width: 20%;
        text-align: center;
        border-radius: 5px;
    }

    .reservation-container .class-luxury {
        background-color: goldenrod;
    }
    .reservation-container .class-high-end {
        background-color: springgreen;
    }
    .reservation-container .class-standard {
        background-color: deepskyblue;
    }
    .reservation-container .class-slave {
        background-color: tomato;
    }

    #user-data {
        width: 30%;
        height: 91vh;
        padding: 5vh 2vw;
        box-shadow: 0 60px 40px 10px rgba(0, 0, 0, 0.3);
    }

    #ud-container {
        display: flex;
        height: 30vh;
        background-color: var(--dimtheme-color);
        flex-direction: column;
        font-size: 1.3rem;
        justify-content: space-between;
        color: white;
        border-top: 2px solid var(--theme-color);
        border-bottom: 2px solid var(--theme-color);
        padding: 4vh 1vw;
    }
</style>

<div id="user-profile-page">
    <nav id="user-profile-bar">
        <div on:click={logout} id="logout-btn">logout</div>
        <div id="username-title">{user_name}</div>
    </nav>
    <main id="main-content">
        {#if user_data.user !== undefined && user_data.reservations !== undefined }
            <section id="reservations-container">
                <div id="reservations">
                    {#each user_data.reservations as r, h}
                        <div on:click={() => cancelReservation(user_data.reservations[h])} class="reservation-container">
                            <span class="room-name">room {r.number}</span>
                            <span class="class-{ROOMS_CLASSES[r.class].toLowerCase().replace(" ", "-")} room-class">{ROOMS_CLASSES[r.class]}</span>
                        </div>
                    {/each}
                </div>
                <div id="reservation-controls">
                    <Button onClick={() => push(`/reservations/${session_key}`)} label="make a reservation"/>
                </div>
            </section>
            <aside id="user-data">
                <div id="ud-container">
                    <div class="ud-field">real name: {user_data.user.name}</div>
                    <div class="ud-field">phone: {user_data.user.phone}</div>
                    <div class="ud-field">email: {user_data.user.email}</div>
                    <div class="ud-field">reservations: {user_data.reservations.length}</div>
                </div>
            </aside>
        {/if}
    </main>
</div>