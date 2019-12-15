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
            </div>

            <div slot="bottom">
                <va-input placeholder="Filter" width="md"></va-input>
            </div>
        </va-page-header>

        <va-loading v-if="loading" center/>
        <va-row :gutter="15">
            <va-column :xs="12" :sm="6" :md="6">
                <h2>Not reviewed</h2>
                <complaint v-for="c in unreviewed" v-on:reviewed="fetchData" :complaint="c"/>
            </va-column>
            <va-column :xs="12" :sm="6" :md="6">
                <h2>Reviewed</h2>
                <complaint v-for="c in reviewed" :complaint="c"/>
            </va-column>
        </va-row>
    </va-page>
</template>

<script>
    import Complaint from "@/components/Complaint";
    import {all_complaints} from "@/_helpers/complaints";

    export default {
        name: "Complaints",
        components: {Complaint},
        data() {
            return {
                loading: false,
                complaints: [],
                error: null
            }
        },
        created() {
            this.fetchData()
        },
        computed: {
            reviewed() {
                return this.complaints.filter((v,i,_) => v.answer)
            },
            unreviewed() {
                return this.complaints.filter((v,i,_) => !v.answer)
            }
        },
        methods: {
            fetchData() {
                this.loading = true;
                all_complaints(this.$store.state.token)
                    .then(result => {
                        this.complaints = result;
                    })
                    .catch(reason => {
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