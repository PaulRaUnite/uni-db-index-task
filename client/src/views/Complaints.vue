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
                    <va-icon type="plus" icon-style="regular" margin="0 7px 0 0"></va-icon>
                    Create complaint
                </va-button>
            </div>

            <div slot="bottom">
                <va-input placeholder="Filter" width="md"></va-input>
            </div>
        </va-page-header>

        <complaint v-for="c in complaints" v-bind:complaint="c"/>
    </va-page>
</template>

<script>
    import Complaint from "@/components/Complaint";

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
                this.error = this.complaints = null;
                this.loading = true;
                fetch(`http://api.localhost/complaint`).then(resp => {
                    if (!resp.ok) {
                        return Promise.reject(resp.statusText)
                    }
                    return resp.json()
                }).catch(reason => {
                    this.error = reason;
                }).then(result => {
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