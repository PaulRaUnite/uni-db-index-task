<template>
    <va-dropdown>
        <div slot="trigger">
            <va-button type="dark" round>
                <va-icon :type="icon_type" size="1.25em"/>
            </va-button>
        </div>
        <div class="dropdown-sp">
            <sign-in theme="primary" v-if="!$store.state.logged"/>
            <va-card v-if="$store.state.logged">
                <va-button block @click="() => $router.push(`/user/${$store.state.username}`)"
                           icon-before="user-cog">
                    Account
                </va-button>
                <va-button block @click="() => $router.push(`/cart`)"
                           icon-before="shopping-cart">
                    In cart: {{$store.state.cart_total}}&#163;
                </va-button>
                <va-button block icon-before="sign-out-alt" @click="logout">
                    Sign out
                </va-button>
            </va-card>
        </div>
    </va-dropdown>
</template>

<script>
    import SignIn from "@/components/SignIn";

    export default {
        name: "InlineAccount",
        components: {SignIn},
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
                this.$store.commit("signout");
                this.$router.push("/inventory")
            }
        }
    }
</script>

<style module>
    .dropdown-sp {
        width: 400px;
    }
</style>