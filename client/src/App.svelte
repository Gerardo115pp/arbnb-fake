<script>
    import Router from 'svelte-spa-router';
    import { routes } from './routes';
</script>

<style>
    @import url('https://fonts.googleapis.com/css2?family=Rhodium+Libre&display=swap');

    :global(:root) {
        --theme-color: #7000FF;
        --dimtheme-color: #C600CA;
        --theme-gradiant: linear-gradient(5.45deg, #7000FF -45.19%, #C600CA 74.22%);
        --color-gradiant: linear-gradient(12.69deg, #FF0DD8 35.54%, #7000FF 85.07%, #C69AFF 136.08%);
        --dark-color: #ffffff;
        --danger:#e72d2d;
        --ready:#25ff3d;
        --modal-shadows: #0000001a;
        --main-font: 'Rhodium Libre';

    }

    :global(body) {
        font-family: 'Rhodium Libre';
        overflow-y: hidden;
        background-color: var(--dark-color);
        color: var(--theme-color);
        margin: 0;
    }

    :global(select) {
        appearance: none;
        background-color: transparent;
        border: none;
        padding: 0 1em 0 0;
        margin: 0;
        width: 100%;
        font-family: inherit;
        font-size: inherit;
        cursor: inherit;
        line-height: inherit;
    }

    :global(select::-ms-expand) {
        display: none;
    }

    #app-container {
        width: 100%;
        height: 100vh;
    }
</style>

<div id="app-container">
    <Router {routes}/>
</div>