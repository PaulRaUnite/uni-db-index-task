<template>
    <va-card :elevation="elevation" :padding="padding" class="card">
        <div class="title">
            <h2>Complaints</h2>
            <va-badge type="info" class="counter" margin="8px">{{complaints.length}}</va-badge>
            <va-loading v-if="loading" color="blue" size="sm" class="spinner"/>
            <va-button style="float: right" icon-before="plus" type="primary" @click="showModal">Create complaint
            </va-button>
        </div>
        <va-modal title="Create complaint" width="600px" ref="customModal" :backdrop-clickable="true"
                  @confirm="createComplaint">
            <div slot="body">
                <va-form type="vertical">
                    <va-form-item label="Reason to complain/question to ask">
                        <va-textarea v-model="complaint_body" class="body"/>
                    </va-form-item>
                </va-form>
            </div>
            <div slot="footer">
                <div style="text-align: right;">
                    <ctx-button type="primary" :state="confirm_button_state" @click="createComplaint">
                        Confirm
                    </ctx-button>
                    <va-button @click="() => this.$refs.customModal.close()">Close</va-button>
                </div>
            </div>
        </va-modal>
        <va-collapse :accordion="false">
            <va-card v-for="c in complaints">
                <div>
                    <va-lozenge class="lozenge" type="default" :uppercase="true">{{c.created_at}}</va-lozenge>
                    <va-lozenge class="lozenge" :type="c.status === 'reviewed' ? 'success' : 'default'"
                                :uppercase="true">{{c.status}}
                    </va-lozenge>
                </div>
                <va-collapse-panel header="Details">
                    <va-form type="vertical">
                        <va-form-item label="Your complaint/question">
                            <p>{{c.body}}</p>
                        </va-form-item>
                        <div v-if="c.reviewer">
                            <va-form-item label="Answer">
                                <p>{{c.answer}}</p>
                            </va-form-item>
                            <va-form-item label="Reviewer">
                                <p>{{c.reviewer.name}}</p>
                            </va-form-item>
                            <va-form-item label="Reviewed at">
                                <p>{{c.reviewed_at}}</p>
                            </va-form-item>
                        </div>
                    </va-form>
                </va-collapse-panel>
            </va-card>
        </va-collapse>
    </va-card>
</template>
<script>
    import {get_complaints, new_complaint} from "@/_helpers/complaints";
    import CtxButton from "@/components/CtxButton.vue";

    export default {
        name: 'user-complaints',
        components: {
            "ctx-button": CtxButton,
        },
        data() {
            return {
                loading: false,
                complaints: [],
                complaint_body: "",
                confirm_button_state: "default"
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
            showModal() {
                this.complaint_body = "";
                this.$refs.customModal.open();
                this.confirm_button_state = "default";
            },
            createComplaint() {
                this.confirm_button_state = "loading";
                new_complaint(this.username, this.$store.state.token, this.complaint_body)
                    .then(() => {
                        this.confirm_button_state = "success";
                        this.fetch_data();
                        return new Promise((resolve => setTimeout(resolve, 1000))).then(() => {
                            this.$refs.customModal.close();
                        })
                    }).catch((error) => {
                    console.log(error);
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
                    this.confirm_button_state = "failure"
                })
            },
            fetch_data() {
                this.loading = true;
                get_complaints(this.username, this.$store.state.token)
                    .then((data) => {
                        this.complaints = data;
                        this.loading = false;
                    })
                    .catch((error) => {
                        console.log(error);
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
                        this.loading = false;
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
        margin: 6px 0;
    }

    .card {
        padding: 10px 30px 20px 30px;
        margin-bottom: 14px;
    }

    .lozenge {
        margin-right: 6px;
    }

    .counter {
        display: inline-block;
        float: left;
    }

    .title {
        height: 34px;
    }

    h2 {
        float: left;
        display: inline-block;
        margin: 5px 0;
    }

    .body {
        width: 100%;
    }
</style>