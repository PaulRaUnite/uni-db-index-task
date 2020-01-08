<template>
    <va-app
            bg-color="#F4F5F7"
            page-bg-color="#FFFFFF"
            desktop-margin="0"
            desktop-minimum-width="0"
            desktop-sidebar-width="220"
            desktop-minibar-width="0"
            desktop-topbar-height="50"
            mobile-sidebar-width="0"
            mobile-minibar-width="0"
            mobile-topbar-height="50"
            :rtl="false"
            :reverse="false"
            :split="false"
            :sidebar-priority="false"
            :topbar-priority="false"
            :topbar-padded="false">

        <va-topbar theme="dark">
            <div slot="left">
                <div slot="trigger">
                    <va-button type="dark" @click="() => this.$router.push('/')">
                        <va-icon type="home"/>
                    </va-button>
                </div>
            </div>
            <div slot="right">
                <inline-account/>
            </div>
        </va-topbar>

        <va-sidebar
                :compact="false"
                :text-links="false">
            <va-sidebar-group
                    :items="items"
                    title="Menu"
                    :show-toggle="false"/>
        </va-sidebar>

        <transition>
            <router-view/>
        </transition>
    </va-app>
</template>

<script>
    import InlineAccount from "@/components/InlineAccount";

    export default {
        name: "AppLayout",
        components: {InlineAccount},
        data() {
            return {}
        },
        computed: {
            items() {
                if (this.$store.state.token !== null) {
                    return [
                        {name: 'Inventory', route: '/inventory', icon: 'store-alt'},
                        {name: 'Complaints', route: '/complaints', icon: 'comments'},
                    ]
                } else {
                    return [
                        {name: 'Login', route: '/auth#signin', icon: 'sign-in-alt'},
                        {name: 'Inventory', route: '/inventory', icon: 'store-alt'},
                    ]
                }
            },
        }
    }
</script>
<style lang="scss">
    #nav {
        text-align: center;
        padding: 30px;

        a.router-link-exact-active {
            color: #42b983;
        }
    }
</style>
