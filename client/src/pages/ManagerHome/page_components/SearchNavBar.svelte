<script>
    export let used_rooms = 0;
</script>

<style>
    #search-nav-bar {
        display: flex;
        width: 70%;
        background: var(--color-gradiant);
        color: white;
        padding: .3%;
        margin: 3vh 0;
        justify-content: space-between;
        align-items: center;
        border-radius: 0 15px 15px 0;
        box-shadow: 0 10px 15px 6px rgba(0, 0, 0, 0.2);
    }

    #search-nav-bar input {
        height: 3.5vh;
        width: 40%;
    }

    #snb-nav-options {
        width: 40%;
        display: flex;
        justify-content: space-around;
        font-size: 1.2rem;
    }

    .nav-option {
        cursor: pointer;
    }
</style>

<div id="search-nav-bar">
    <input placeholder="Search" type="text">
    <div id="snb-nav-options">
        <div class="nav-option">
            <span>Reservations: {used_rooms}</span>
        </div>
    </div>
</div>