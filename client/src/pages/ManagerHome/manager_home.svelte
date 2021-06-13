<script>
    import RoomDetailsModal from './page_components/RoomDetailsModal.svelte';
    import HotelComponent from './page_components/HotelComponent.svelte';
    import SearchNavBar from './page_components/SearchNavBar.svelte';

    export let params = {};

    let show_room_modal = false;
    let room_selected = {};
    let used_rooms = 0;

    let { sk:session_key } = params;
    
    const roomSelected = room_data => {
        room_selected = room_data;
        show_room_modal = !show_room_modal;
    }

</script>

<style>
    #mhome-main-container {
        box-sizing: border-box;
        border-top: 1px solid white;
    }

    #hotel-component-container {
        display: flex;
        height: 87.9vh;
        justify-content: center;
        align-items: flex-end;
        
    }
</style>

<RoomDetailsModal {session_key} room_data={room_selected} bind:show_modal={show_room_modal} />
<div id="mhome-main-container">
    <SearchNavBar {used_rooms}/>
    <div id="hotel-component-container">
        <HotelComponent rooms_selected={roomSelected} bind:used_rooms={used_rooms} {session_key}/>
    </div>
</div>