<template>
    <va-card :elevation="2" :pading="12" class="container">
        <va-form type="vertical">
            <va-form-item label="ID">
                {{complaint.id}}
            </va-form-item>
            <va-form-item label="Created at">
                {{complaint.created_at}}
            </va-form-item>
            <va-form-item label="Complaint or question">
                {{complaint.body}}
            </va-form-item>
            <va-form-item label="Requester">
                {{complaint.user.name}}
            </va-form-item>
            <div class="card" v-if="!complaint.answer">
                <va-form-item label="Answer">
                    <va-input :clearable="true" placeholder="type answer" buttons :loading="sending_answer"
                              @confirm="send_answer"
                              @cancel="() => this.answer = ''"
                              v-model="answer"/>
                </va-form-item>
            </div>
            <div v-else>
                <va-form-item label="Answer">
                    <p>{{complaint.answer}}</p>
                </va-form-item>
                <va-form-item label="Reviewer">
                    <p>{{complaint.reviewer.name}}</p>
                </va-form-item>
                <va-form-item label="Reviewed at">
                    {{complaint.reviewed_at}}
                </va-form-item>
            </div>
        </va-form>
    </va-card>
</template>

<script>
    import {review_complaint} from "@/_helpers/complaints";

    export default {
        name: "Complaint",
        props: ["complaint"],
        data() {
            return {
                sending_answer: false,
                answer: "",
            }
        },
        methods: {
            send_answer() {
                this.sending_answer = true;
                review_complaint(this.complaint.id, this.$store.state.token, this.answer)
                    .then((json) => {
                        this.sending_answer = false;
                        this.$emit('reviewed');
                        this.answer = ""
                    })
                    .catch((error) => {
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
                        this.sending_answer = false;
                    })
            }
        }
    }
</script>

<style scoped>
    .container {
        margin-bottom: 20px;
    }

    .card {
        width: 100%;
    }
</style>