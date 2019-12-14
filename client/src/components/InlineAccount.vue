<template>
    <va-dropdown>
        <div slot="trigger">
            <va-button type="dark" round>
                <va-icon :type="icon_type" size="1.25em"/>
            </va-button>
        </div>
        <sign-in theme="primary" v-if="!$store.state.logged"/>
        <va-button v-if="$store.state.logged" block @click="() => $router.push(`/user/${$store.state.username}`)"
                   icon-before="user-cog">
            Account
        </va-button>
        <va-button v-if="$store.state.logged" block icon-before="sign-out-alt" @click="logout">
            Sign out
        </va-button>
    </va-dropdown>
</template>

<script>
    import SignIn from "@/components/SignIn";

    export default {
        name: "InlineAccount",
        components: {SignIn},
        comments: {"sign-in": SignIn},
        computed: {
            icon_type() {
                if (!this.$store.state.logged) {
                    return "user-circle"
                } else {
                    return "user-astronaut"
                }
            }
        },
        methods: {
            logout() {
                this.$store.commit("signout")
            }
        }
    }
</script>

<style module>
    .dropdown {
        width: 300px;
    }
</style>