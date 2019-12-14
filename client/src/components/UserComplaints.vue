<template>
    <va-card :elevation="elevation" :padding="padding" class="card">
        <h2>Orders
            <va-badge type="info" class="counter" margin="0px">{{complaints.length}}</va-badge>
            <va-loading v-if="loading" color="blue" size="sm" class="spinner"/>
        </h2>
        <va-collapse :accordion="false">
            <va-card v-for="c in complaints">
                <div>
                    <va-lozenge class="lozenge" type="default" :uppercase="true">{{c.date}}</va-lozenge>
                    <va-lozenge class="lozenge" type="success" :uppercase="true">{{c.status}}</va-lozenge>
                </div>
                <va-collapse-panel header="Details">
                </va-collapse-panel>
            </va-card>
        </va-collapse>
    </va-card>
</template>
<script>
    import {get_complaints} from "@/_helpers/complaints";

    export default {
        name: 'user-complaints',
        data() {
            return {
                loading: false,
                complaints: [],
            }
        },
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
        created() {
            this.fetch_data()
        },
        methods: {
            fetch_data() {
                this.loading = true;
                get_complaints(this.username, this.$store.state.token)
                    .then((data) => {
                        this.complaints = data;
                    })
                    .catch((error) => {
                        console.log(error)
                        if (error.status === 401) {
                            this.notification.warning({
                                    title: "Unauthorized access."
                                }
                            );
                            this.$router.push("/inventory")
                        } else {
                            this.notification.danger({
                                    title: "Error is occurred during logging.",
                                    message: error.details + ". Try to reload page later.",
                                }
                            )
                        }
                    })
                    .then(() => {
                        this.loading = false;
                    })
            }
        }
    }
</script>
<style scoped>
    .spinner {
        display: inline-block;
    }

    .card {
        padding: 10px 30px 20px 30px;
        margin-bottom: 14px;
    }

    .lozenge {
        margin-right: 6px;
    }

    .counter {
        text-align: center;
    }
</style>