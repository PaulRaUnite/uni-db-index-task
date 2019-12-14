<template>
    <va-page>
        <va-page-header>
            <div slot="breadcrumb">
                <va-breadcrumb separator="/">
                    <va-breadcrumb-item :to="{ path: '/users' }">Users</va-breadcrumb-item>
                    <va-breadcrumb-item>
                        {{username}}
                    </va-breadcrumb-item>
                </va-breadcrumb>
            </div>

            <div slot="title">
                <span>User page</span>
            </div>

            <div v-if="user_itelf" slot="subtitle">
                <span>Your account</span>
            </div>
        </va-page-header>
        <va-row :gutter="16">
            <va-column :xs="12" :sm="6" :md="4">
                <user-info :padding="padding" :elevation="elevation" :username="username"/>
            </va-column>
            <va-column :xs="12" :sm="6" :md="4">
                <user-orders :elevation="elevation" :padding="padding" :username="username"/>
            </va-column>
            <va-column :xs="12" :sm="6" :md="4">
                <user-complaints :elevation="elevation" :padding="padding" :username="username"/>
            </va-column>
        </va-row>


    </va-page>
</template>

<script>
    import UserInfo from "@/components/UserInfo";
    import UserOrders from "@/components/UserOrders";
    import UserComplaints from "@/components/UserComplaints";

    export default {
        name: "user-page",
        components: {UserComplaints, UserOrders, UserInfo},
        data() {
            return {
                user: null,
                elevation: 3,
                padding: 20,
            }
        },
        computed: {
            user_itelf() {
                return this.username === this.$store.state.username
            },
            username() {
                return this.$route.params.username
            }
        },
    }
</script>

