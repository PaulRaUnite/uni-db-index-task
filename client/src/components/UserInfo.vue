<template>
    <va-card :elevation="elevation" :padding="padding" class="card">
        <div class="title">
            <h2>User info</h2>
            <va-loading v-if="loading" color="blue" size="sm" class="spinner"/>
        </div>
        <va-form type="vertical">
            <va-form-item label="Name">
                <va-input
                        buttons
                        :loading="loading"
                        placeholder="Name"
                        @confirm="update_name"
                        @cancel="cancel_name_change"
                        v-model="new_name">
                </va-input>
            </va-form-item>
        </va-form>
        <h2>Password</h2>
        <va-form type="vertical">
            <va-form-item label="Old password" need>
                <va-input placeholder="Password"/>
            </va-form-item>
            <va-form-item label="New password" need>
                <va-input placeholder="Password"/>
            </va-form-item>
            <va-form-item>
                <va-button>Change password</va-button>
            </va-form-item>
        </va-form>
    </va-card>
</template>

<script>
    import {jsoner} from "../_helpers/jsoner";
    import {withjwt} from "../_helpers/withjwt";

    export default {
        name: "user-info",
        props: {
            elevation: {
                type: Number,
            },
            padding: {
                type: Number,
            },
            username: {
                type: String,
            }
        },
        data() {
            return {
                loading: false,
                new_name: null,
                user: null,
            }
        },
        created() {
            this.fetch_data()
        },
        methods: {
            fetch_data() {
                this.loading = true;
                jsoner(withjwt(`http://api.localhost/user/${this.username}`, this.$store.state.token))
                    .then((json) => {
                        this.user = {
                            id: json.data.id,
                            ...json.data.attributes
                        };
                        this.new_name = this.user.name;
                    })
                    .catch((error) => {
                        if (error.status === 401) {
                            this.notification.warning({
                                    title: "Unauthorized access."
                                }
                            );
                            this.$router.push("/inventory")
                        } else if (error.status === 404) {
                            this.$router.push("/inventory")
                        } else {
                            this.notification.danger({
                                    title: "Error is occurred during logging.",
                                    message: error.details + ". Try to reload page later.",
                                }
                            )
                        }
                    }).then(() => {
                    this.loading = false;
                })
            },
            update_name() {

            },
            cancel_name_change() {
                this.new_name = this.user.name
            },
        },
    }
</script>

<style scoped>

    .spinner {
        margin: 6px 0;
        display: inline-block;
    }

    .card {
        padding: 10px 30px 20px 30px;
        margin-bottom: 14px;
    }

    .title {
        height: 34px;
    }
    h2 {
        margin: 5px 0;
    }
</style>