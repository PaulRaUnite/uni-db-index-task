<template>
    <va-form ref="form" type="vertical">
        <va-form-item label="Username" need>
            <va-input name="username" v-model="form.username"
                      :rules="[{type:'required', tip:'Please input your username'}]"/>
        </va-form-item>
        <va-form-item label="Name" need>
            <va-input name="name" v-model="form.name" :rules="[{type:'required', tip:'Please input your name'}]"/>
        </va-form-item>
        <va-form-item label="Password" need>
            <va-input name="username" v-model="form.password"
                      :rules="[{type:'required', tip:'Please input password of your account'}]"/>
        </va-form-item>
        <va-form-item>
            <va-button block class="button" type="primary" :loading="submitting" @click="submit">Sign Up</va-button>
        </va-form-item>
    </va-form>
</template>

<script>
    import {signup} from "@/_helpers/signup";

    export default {
        name: "Signing",
        data() {
            return {
                submitting: false,
                form: {
                    username: null,
                    name: null,
                    password: null,
                }
            }
        },
        methods: {
            submit() {
                this.$refs.form.validateFields((result) => {
                    if (!result.isvalid) {
                        return
                    }
                    this.submitting = true;
                    signup(this.form.username, this.form.password, this.form.name)
                        .then(
                            (_) => {
                                this.$router.push("/auth#signin");
                            }
                        ).catch((error) => {
                        if (error.status === "409") {
                            this.notification.warning({
                                    title: "Username conflict.",
                                    duration: 3000,
                                }
                            )
                        } else {
                            this.notification.danger({
                                    title: "Error is occurred during logging.",
                                    message: error.details,
                                    duration: 3000,
                                }
                            )
                        }
                    }).then(() => {
                        this.submitting = false;
                    })
                })
            }
        }
    }
</script>

<style scoped>
</style>