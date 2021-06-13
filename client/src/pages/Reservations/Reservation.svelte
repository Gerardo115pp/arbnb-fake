<script>
    import Button from '../../components/Button.svelte';
    import { ROOMS_CLASSES } from '../../enums/rooms';
    import { server_name } from '../../server_info';
    import { pop } from 'svelte-spa-router';
    import Datepicker from 'svelte-calendar';
    import { onMount } from 'svelte';



    export let params = {};
    
    const one_day = 1000*60*60*24;
    const { sk:session_key } = params;

    let is_check_in_set, is_check_out_set;
    let billing_description;
    let check_in_date = new Date(Date.now());
    let check_out_date = new Date(check_in_date.getDate()+1); // today + 1 day
    let costs_details;
    let room_selector;
    let selected_room;
    let rooms = [];

    $: if(rooms.length > 0) {
        selected_room = rooms[0];
    }
    
    onMount(() => {
        const headers = new Headers();
        headers.set("X-sk", session_key);

        fetch(`${server_name}/reservations?guest=a`, {method: 'GET', headers: headers})
            .then(promise => {
                if (promise.ok) {
                    promise.json().then(response => {
                        rooms = response;
                    });
                }
            })
    });

    const getNDays = () => {
        const unix_timestamp = check_out_date - check_in_date;
        let days = unix_timestamp/one_day;
        
        return Math.ceil(days);
    }

    const getDescription = () => {
        if (selected_room && costs_details !== undefined) {
            return `${ROOMS_CLASSES[selected_room.class]} room at $${costs_details.costs_per_day.toLocaleString('en')} MXN for ${getNDays()} days`;
        } else {
            return "";
        }
    }

    const handleRoomChange = () => {
        const room_index =  parseInt(room_selector.value);
        if (selected_room === undefined || selected_room.number !== rooms[room_index].number) {

            let old_room = selected_room;
            selected_room = rooms[room_index];

            if (old_room === undefined || old_room.class !== selected_room.class || costs_details === undefined) {
                requestCostDetails();
            }

        }
    }

    const makeReservation = () => {
        const headers = new Headers();
        headers.set("X-sk", session_key);
        // headers.set("Content-Type", "x-www-form-urlencoded")

        const form_data = new FormData();
        form_data.append('room_id', selected_room.number);
        form_data.append('from', check_in_date.valueOf());
        form_data.append('to', check_out_date.valueOf());

        const request = new Request(`${server_name}/reservations`, {method: 'PATCH', body: form_data, headers: headers});

        fetch(request)
            .then(promise => {
                if(promise.ok) {
                    pop();
                }
            })
    }

    const requestCostDetails = () => {
        fetch(`${server_name}/budget?class=${selected_room.class}&days=${getNDays()}`)
            .then(promise => {
                if(promise.ok) {
                    promise.json().then(response => costs_details = response);
                }
            })
    }



    $:billing_description = getDescription(), check_in_date, check_out_date, costs_details, selected_room;
    
    $: if (selected_room !== undefined && check_in_date !== check_out_date) {
        requestCostDetails(); // call requestCostDetails every time any date is changed
    }
</script>

<style>
    #reservations-page {
        display: flex;
        height: 100vh;
        background: var(--theme-gradiant);
    }

    #reservations-page > div {
        width: 50vw;
        height: 100vh;
    }

    #price-container {
        display: flex;
        color: white;
        font-size: 4rem;
        justify-content: center;
        align-items: center;
        user-select: none;
    }

    #resume-container {
        background-color: white;
        box-shadow: 0 0 60px 40px rgba(0, 0, 0, 0.2);
    }   

    #resume-title {
        font-size: 1.8rem;
        margin: 4vh 0;
        text-align: center;
    }

    #reservation-form {
        box-sizing: border-box;
        display: flex;
        height: 80vh;
        flex-direction: column;
        justify-content: space-around;
        align-items: center;
        padding: 2vh 2vw;
    }

    #room-picker {
        border: none;
        display: flex;
        justify-content: space-around;
    }

    #room-picker option {
        display: flex;
        justify-content: space-around;
        color: aqua;
    }

    .reservation-field {
        cursor: pointer;
        background: var(--theme-gradiant);
        color: white;
        padding: 1vh 2vw;
        font-size: 1.2rem;
        border-radius: 5px;

    }

    #billing-section {
        display: flex;
        width: 100%;
        height: 20vh;
        flex-direction: column;
        justify-content: space-around;
        align-items: center;
    }
</style>

<div id="reservations-page">
    <div id="price-container">
        <span id='price-label'>
            {#if costs_details !== undefined}
                $ {costs_details.total.toLocaleString('en')} MXN
            {:else}
                Loading
            {/if}
        </span>
    </div>
    <div id="resume-container">
        <div id="resume-title">Make a reservation</div>
        <div id="reservation-form">
            <div id="room-picker">
                <!-- svelte-ignore a11y-no-onchange -->
                <select bind:this={room_selector} on:change={handleRoomChange} name="rooms" id="room-picker" class="reservation-field">
                    {#each rooms as room, h}
                        <option class="room-option" value="{h}"> 
                            <span class="room-name">Room {room.number} - </span>
                            <span class="room-class">{ROOMS_CLASSES[room.class]}</span>
                        </option>
                    {/each}
                </select>
            </div>
            <Datepicker highlightColor='var(--theme-color)' start={new Date(Date.now())} bind:selected={check_in_date} bind:dateChosen={is_check_in_set}>
                <div class="reservation-field" id="reservation-check-in">
                    Check in: {check_in_date.toLocaleDateString()} 
                </div>
            </Datepicker>
            <Datepicker highlightColor='var(--theme-color)' start={new Date(Date.now() + one_day)} bind:selected={check_out_date} bind:dateChosen={is_check_out_set}>
                <div class="reservation-field" id="reservation-check-out">
                    Check out: {check_out_date.toLocaleDateString()}
                </div>
            </Datepicker>
            <section id="billing-section">
                <span id="billing-description">Description: {billing_description}</span>
                <Button onClick={makeReservation} width="40%" label="Make reservation"/>
            </section>
        </div>
    </div>
</div>