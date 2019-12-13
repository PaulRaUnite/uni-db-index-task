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

        <div id="nav">
            <span>Logged: {{this.$store.state.token != null}}</span>
            <va-breadcrumb separator="/">
                <va-breadcrumb-item to="/">
                    <router-link to="/">app</router-link>
                </va-breadcrumb-item>
                <va-breadcrumb-item v-if="$route.path==='/login'" to="/login">
                    <router-link to="/login">login</router-link>
                </va-breadcrumb-item>
                <va-breadcrumb-item v-if="$route.path==='/complaints'" to="/complaints">
                    <router-link to="/complaints">complaints</router-link>
                </va-breadcrumb-item>
            </va-breadcrumb>
        </div>
        <va-topbar theme="blue">
            <div slot="left">

                <va-dropdown>
                    <div slot="trigger">
                        <va-button type="primary-dark" round>
                            <va-icon type="bars" color="white"></va-icon>
                        </va-button>
                    </div>
                    <li>
                        <router-link :to="'/'">Home</router-link>
                    </li>
                    <li>
                        <router-link :to="'/documentation'">Documentation</router-link>
                    </li>
                </va-dropdown>

                <span style="font-weight:700;margin:0 20px 0 10px;">
                    Documentation
                </span>

                <va-dropdown style="margin-right: 10px;">
                    <div slot="trigger">
                        <va-button type="primary-dark">
                            Templates
                            <va-icon type="angle-down" margin="0 2px 0 10px"></va-icon>
                        </va-button>
                    </div>
                    <li>
                        <router-link to="/templates/signin">Sign-in</router-link>
                    </li>
                </va-dropdown>

            </div>
            <div slot="right">
                <va-input :clearable="true" icon="search" icon-style="solid" style="margin-right:7px;"></va-input>
                <va-button type="primary-dark" round @click="openAside">
                    <va-icon type="cog" size="1.25em"></va-icon>
                </va-button>
            </div>
        </va-topbar>

        <va-sidebar
                theme="blue"
                :compact="false"
                :text-links="false">
            <va-sidebar-group
                    :items="[{name:'Login', route:'/login'}]"
                    title="Basics"
                    :show-toggle="false"/>
            <va-sidebar-group
                    :items="[{name:'View all', route:'/complaints'}]"
                    title="Complaints"
                    :show-toggle="false"/>
            <va-sidebar-group
                    title="Debug"
                    :items="debug_values"/>
        </va-sidebar>

        <transition>
            <router-view/>
        </transition>
    </va-app>
</template>

<script>
    export default {
        name: "App",
        data() {
            return {}
        },
        computed: {
            debug_values() {
                let items = [];
                if (this.$store.state.user_id !== null) {
                    items.push(
                        {
                            name: this.$store.state.user_id.toString()
                        }
                    )
                }
                return items
            }
        }
    }
</script>
<style lang="scss">
    #app {
        font-family: 'Avenir', Helvetica, Arial, sans-serif;
        -webkit-font-smoothing: antialiased;
        -moz-osx-font-smoothing: grayscale;
        color: #2c3e50;
    }

    #nav {
        text-align: center;
        padding: 30px;

        a.router-link-exact-active {
            color: #42b983;
        }
    }
</style>
