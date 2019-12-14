<template>
    <va-app
            bg-color="#F4F5F7"
            page-bg-color="#FFFFFF"
            desktop-margin="0"
            desktop-minimum-width="0"
            desktop-sidebar-width="0"
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
        </va-topbar>

        <transition>
            <router-view/>
        </transition>
    </va-app>
</template>

<script>
    import SignIn from "@/components/SignIn";

    export default {
        name: "BlankLayout",
        components: {SignIn},
    }
</script>

<style scoped>

</style>