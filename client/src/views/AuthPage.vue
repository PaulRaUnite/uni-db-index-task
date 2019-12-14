<template>
    <div class="container">
        <va-card elevation="3" class="auth">
            <div class="icon"><va-icon type="sign-in-alt" color="coral" size="4em"/></div>
            <h1>{{header}}</h1>
            <va-tabs @changed="onTabChange" :options="{useUrlFragment: true}">
                <va-tab name="SignIn" class="tab">
                    <sign-in/>
                </va-tab>
                <va-tab name="SignUp" class="tab">
                    <sign-up/>
                </va-tab>
            </va-tabs>
        </va-card>
    </div>
</template>

<script>
    import SignIn from "@/components/SignIn";
    import SignUp from "@/components/SignUp";

    export default {
        name: "AuthPage",
        components: {SignIn, SignUp},
        data() {
            return {
                header: "SignIn"
            }
        },
        methods: {
            onTabChange(event) {
                const tab_name = event.tab.$props.name;
                if (tab_name === "SignUp") {
                    this.header = "Sign Up"
                } else {
                    this.header = "Sign In"
                }
            }
        }
    }
</script>

<style scoped>
    .container {
        position: fixed;
        inset: 0;
        overflow: auto;
        height: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .auth {
        width: 440px;
        height: 460px;
        padding: 20px;
    }
    .tab {
        padding: 20px 10px 10px;
    }

    .icon {
        text-align: center;
    }
    .auth h1 {
        text-align: center;
    }
</style>