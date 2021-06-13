<script>
    import FieldData, { FieldStates } from '../../classes/FieldData';
    import Button from '../../components/Button.svelte';
    import Input from '../../components/Input.svelte';
    import { server_name } from '../../server_info'
    import { push } from 'svelte-spa-router'


    let is_form_ready = false;
    let form_inputs = [
        new FieldData('new-guest-username', /^[a-zA-Z\d_\-]+$/, 'username'),
        new FieldData('new-guest-name', /^[a-zA-Z\s]+$/, 'name'),
        new FieldData('new-guest-phone', /^[\-\d]+$/, 'phone'),
        new FieldData('new-guest-email', /^[a-zA-Z_\d]+@[a-z]+(\.[a-z]+)+$/, 'email'),
        new FieldData('new-guest-password', /^[^\n\s]+$/, 'password', 'password'),
    ]

        const registerUser = () => {
        if (is_form_ready) {
            const form_data = new FormData();
            form_inputs.forEach(fd => {
                form_data.append(fd.name, fd.getFieldValue());
            });

           const request = new Request(`${server_name}/register`, {method: "POST", body: form_data});

            fetch(request)
                .then(promise => {
                    if(promise.ok) {
                        push("/");
                    }
                })
        } else {
            alert("please fill all registration values");
        }
    }

    const verifyForm = () => {
        let fd_value = "";
        let verification_result = is_form_ready;
        for(let fd of form_inputs) {
            fd_value = fd.getFieldValue();
            verification_result = false;
            if (fd_value === "") {
                fd.state = FieldStates.NORMAL;
                break;
            } else if(!fd.isReady()) {
                fd.state = FieldStates.HAS_ERRORS;
                break;
            } else {
                verification_result = true;
                fd.state = FieldStates.READY;
            }
        }

        is_form_ready = verification_result;
        form_inputs = [...form_inputs];
    }

</script>

<style>
    #sign-up-container {
        position: relative;
        box-sizing: border-box;
        width: 100%;
        height: 100vh;
    }

    #background-container {
        background: var(--color-gradiant);
        height: 35vh;
    }

    #sign-up-from {
        position: absolute;
        width: 90vw;
        height: 80vh;
        background: white;
        border-radius: 150px 0 0 0;
        margin: 20vh 0 0 0;
        padding: 6vh 5vw;
        top: 0;
        box-shadow: 0 0 20px 10px rgba(0, 0, 0, 0.2);
    }

    #form-title {
        font-size: 1.4rem;
    }

    #form-fields {
        display: flex;
        flex-wrap: wrap;
        width: 90vw;
        justify-content: space-evenly;
        align-items: center;
    }

    .sign-up-field {
        width: 40vw;
        box-sizing: border-box;
        margin-top: 5vh;
    }

    .sign-up-field:first-child, .sign-up-field:nth-child(2) {
        width: 83vw;
    }

    .sign-up-field:nth-last-child(2) {
        margin-right: auto;
        margin-left: 3vw;
    }
    
    #sign-up-btn {
        width: 5vw;
        height: 5vh;
        margin: 5vh auto 0;
    }

    #form-controls {
        margin: 5vh 3vw;
    }
</style>

<div id="sign-up-container">
    <div id="background-container">

    </div>
    <div id="sign-up-from">
        <div id="form-title">Sign up</div>
        <div id="form-fields">
            {#each form_inputs as fd}
                <div class="sign-up-field">
                    <Input onBlur={verifyForm} onEnterPressed={verifyForm} isClear={true} field_data={fd}/>
                </div>
            {/each}
            <div id="sign-up-btn">
                <Button onClick={registerUser} isMaterial={true} font_size="1.2rem" padding="3.5vh .8vw" isRounded={true} label='Sumit'/>
            </div>
        </div>
        <div id="form-controls">
            <Button onClick={() => push("/")} isDimmed={true} width="5%" padding="0" label="Log in"/>
        </div>
    </div>
</div>