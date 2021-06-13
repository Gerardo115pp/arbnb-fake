<script>
    import  FieldData, { FieldStates } from '../../classes/FieldData';
    import Button from '../../components/Button.svelte';
    import HotelLogo from '../../resources/review.svg';
    import Input from '../../components/Input.svelte';
    import { server_name } from '../../server_info';
    import { push } from 'svelte-spa-router';

    
    let is_form_ready = false;
    let field_inputs = [
        new FieldData("login-username-field", /^[A-Za-z\d_\s\-]+$/, "username"),
        new FieldData("login-password-field", /^[^\s]+$/, "password", "password")
    ];

    const logUser = () => {
        if (is_form_ready) {
            const username = field_inputs[0].getFieldValue();
            const password = field_inputs[1].getFieldValue();

            const request = new Request(`${server_name}/login?username=${username}&password=${password}`);
            fetch(request)
                .then(promise => {
                    if(promise.ok) {
                        promise.json().then(response => {
                            if(response.is_manager) {
                                push(`/home-manager/${response.response}`);
                            } else {
                                push(`/user-profile/${response.response}`);
                            }
                        });
                    } else if (promise.status === 400) {
                        alert("password doesnt match");
                    } else if (promise.status === 404) {
                        alert("user doesnt exists");
                    } else {
                        alert("you should never see me, odd... o.o");
                    }
                })
        } else {
            alert("the login fields are not completed");
        }
    }

    const verifyForm = () => {
        let fd_value = "";
        let verification_result = is_form_ready;
        for(let fd of field_inputs) {
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
        field_inputs = [...field_inputs];
    }

</script>

<style>
    #login-upper-banner {
        display: flex;
        height: 40vh;
        background: var(--theme-gradiant);
        flex-direction: column;
        color: white;
        align-items: center;
        border-radius: 0 0 0 200px;


    }

    #lub-logo-container {
        width: 10vw;
        fill: white;
        margin: 5vh 0 0 0;
    }

    #login-form-container {
        display: flex;
        height: 60vh;
        justify-content: center;
        flex-direction: column;
        align-items: center;
    }

    #form-fields-container {
        width: 40vw;
    }

    .field-input-container:first-child {
        margin-bottom: 5vh;
    }

    #form-login-controls {
        display: flex;
        position: relative;
        width: 30vw;
        margin: 4vh 0;
        justify-content: flex-start;
        align-items: center;
    }

    .center-btn {
        position: absolute;
        left: 13vw;
    }


</style>

<div id="login-page">
    <div id="login-upper-banner">
        <div id="lub-logo-container">
            {@html HotelLogo}
        </div>
        <div id="lub-title">Arbnb fake</div>
    </div>
    <div id="login-form-container">
        <div id="form-fields-container">
            {#each field_inputs as fd}
                <div class="field-input-container">
                    <Input onBlur={verifyForm} onEnterPressed={verifyForm} field_data={fd}/>
                </div>
            {/each}
        </div>
        <div id="form-login-controls">
            <div class="login-button-container">
                <Button onClick={() => push("/sign-up")} isClear={true} label="Sign up"/>
            </div>
            <div class="login-button-container center-btn">
                <Button onClick={logUser} padding="1.5vh .5vw" isRounded={true} label="Login"/>
            </div>
        </div>
    </div>
</div>