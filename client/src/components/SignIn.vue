<template>
    <va-form ref="form" type="vertical">
        <va-form-item label="Username" need>
            <va-input :theme="theme" name="username" v-model="form.username"
                      :rules="[{type:'required', tip:'Please input your username'}]"/>
        </va-form-item>
        <va-form-item label="Password" need>
            <va-input :theme="theme" :type="password_input_type" name="username" v-model="form.password"
                      :rules="[{type:'required', tip:'Please input password of your account'}]"/>

            <va-tooltip
                    trigger="hover"
                    content="Show password"
                    placement="top"
                    effect="tooltip-fade-top">
                <va-toggle :theme="theme" class="toggle" v-model="show_password"/>
            </va-tooltip>
        </va-form-item>
        <va-form-item>
            <va-button block type="primary" :loading="logging_in" @click="login">Sign In</va-button>
        </va-form-item>
    </va-form>
</template>

<script>
    import {login} from "@/_helpers/login";

    export default {
        name: "sign-in",
        props: {
            theme: {
                type: String,
                default: 'primary',
                required: false,
                validator(v) {
                    return [
                        'default',
                        'primary',
                        'success',
                        'warning',
                        'danger',
                        'purple'
                    ].includes(v)
                }
            }
        },
        data() {
            return {
                show_password: false,
                logging_in: false,

                form: {
                    username: null,
                    password: null,
                },
            }
        },
        computed: {
            password_input_type() {
                if (this.show_password) {
                    return 'text'
                } else {
                    return 'password'
                }
            }
        },
        methods: {
            login() {
                this.$refs.form.validateFields((result) => {
                    if (!result.isvalid) {
                        return
                    }
                    this.logging_in = true;
                    login(this.form.username, this.form.password).then((jwt) => {
                        this.$store.commit("signin", {token: jwt, username: this.form.username});
                        this.$router.push("/inventory")
                    }).catch((error) => {
                        console.log(error.status);
                        if (error.status === 401) {
                            this.notification.warning({
                                    title: "Invalid credentials."
                                }
                            )
                        } else {
                            this.notification.danger({
                                    title: "Error is occurred during logging.",
                                    message: error.details,
                                }
                            )
                        }
                    }).then(() => {
                        this.logging_in = false;
                    })

                });

            }
        }
    }
</script>

<style module>
</style>