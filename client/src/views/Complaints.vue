<template>
    <va-page>
        <va-page-header>
            <div slot="breadcrumb">
                <va-breadcrumb separator="/">
                    <va-breadcrumb-item to="/complaints">Complaints</va-breadcrumb-item>
                </va-breadcrumb>
            </div>

            <div slot="title">
                <span>Complaints</span>
            </div>

            <div slot="actions">
                <va-button type="primary">
                    <va-icon type="plus" margin="0 7px 0 0"></va-icon>
                    Create complaint
                </va-button>
            </div>

            <div slot="bottom">
                <va-input placeholder="Filter" width="md"></va-input>
            </div>
        </va-page-header>

        <va-loading v-if="loading"/>
        <complaint v-for="c in complaints" v-bind:complaint="c"/>
    </va-page>
</template>

<script>
    import Complaint from "@/components/Complaint";
    import {jsoner} from "@/_helpers/jsoner";
    import {withjwt} from "@/_helpers/withjwt";

    export default {
        name: "Complaints",
        components: {Complaint},
        data() {
            return {
                loading: false,
                complaints: null,
                error: null
            }
        },
        created() {
            this.fetchData()
        },
        watch: {
            '$route': 'fetchData'
        },
        methods: {
            newComplaint(text) {
                postMessage()
            },
            fetchData() {
                this.loading = true;
                jsoner(withjwt(`http://api.localhost/complaint`, this.$store.state.token))
                    .then(result => {
                        this.complaints = result.data.map((v, i, _) => {
                            let user_id = v.relationships.user.data.id;
                            let user = {
                                id: user_id,
                                name: result.included.find((v, i, _) => {
                                    return v.type === "customers" && v.id === user_id
                                }).attributes.name
                            };
                            let reviewer;
                            if (v.relationships.reviewer.data) {
                                const reviewer_id = v.relationships.reviewer.data.id;
                                reviewer = {
                                    id: reviewer_id,
                                    name: result.included.find((v, i, _) => {
                                        return v.type === "customers" && v.id === reviewer_id
                                    }).attributes.name
                                }
                            }
                            return {
                                id: v.id,
                                body: v.attributes.body,
                                answer: v.attributes.answer,
                                user: user,
                                reviewer: reviewer,
                            }
                        });
                    }).catch(reason => {
                    this.notification.warning({
                        title: "Error occured",
                        message: reason,
                    })
                }).finally(() => {
                    this.loading = false;
                })
            }
        }
    }
</script>

<style scoped>
    ul {

        padding: 0;
    }

    li {
        list-style: none;
    }
</style>