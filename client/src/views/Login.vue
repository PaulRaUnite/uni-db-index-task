<template>
    <va-page article="false">
        <va-page-header>
            <div slot="breadcrumb">
                <va-breadcrumb separator="/">
                    <va-breadcrumb-item to="/login">Login</va-breadcrumb-item>
                </va-breadcrumb>
            </div>

            <div slot="title">
                <span>Login</span>
            </div>

            <div slot="bottom">
                <va-input autofocus placeholder="Search users" width="md" icon="search" icon-style="solid" size="large"
                          v-model="search"/>
            </div>
        </va-page-header>

        <div class="container">
        <va-card class="card" v-for="u in users" :elevation="2" :pading="12">
            <p>Id: {{u.id}}</p>
            <p>Name: {{u.name}}</p>
            <va-button v-on:click="() => {handleLogin(u.id)}">Login as {{u.name}}</va-button>
        </va-card></div>
    </va-page>
</template>

<script>
    import {login} from "../_helpers/login"

    export default {
        name: "Login",
        data() {
            return {
                search: "",
                users: null,
                error: null,
            }
        },
        watch: {
            search: function (value) {
                return fetch(`http://api.localhost/user?filter[name]=${value}`).then(resp => {
                    if (!resp.ok) {
                        return Promise.reject(resp.statusText)
                    }
                    return resp.json().then(data => {
                        return data.data
                    })
                }).catch(reason => {
                    this.error = reason;
                }).then(data => {
                    this.users = data.map((v, i, _) => {
                        return {
                            id: v.id,
                            ...v.attributes
                        }
                    });
                })
            }
        },
        methods: {
            handleLogin: function (id) {
                login(id).then(token => this.$store.commit("login", token)
                )
            },
        }
    }
</script>

<style scoped>
    .container {
        width: 800px;
    }

    .card + .card {
        margin-top: 20px;
    }
</style>