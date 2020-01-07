<template>

    <va-page>
        <va-page-header>
            <div slot="breadcrumb">
                <va-breadcrumb separator="/">
                    <va-breadcrumb-item to="/inventory">Inventory</va-breadcrumb-item>
                </va-breadcrumb>
            </div>

            <div slot="title">
                <span>Inventory</span>
            </div>

            <div slot="actions">
            </div>

            <div slot="bottom">
                <va-input placeholder="Search by description" width="md" v-model="descriptionFilter" clearable
                          @keyup.enter.native="filterChanged"/>
            </div>
        </va-page-header>

        <va-loading v-if="loading" center/>
        <div>
            <va-row :gutter="15">
                <va-column :xs="12" :sm="6" :md="4" v-for="g in goods" class="good-card">
                    <good :id="g.id" :description="g.description" :amount="g.amount"/>
                </va-column>
            </va-row>
            <va-pagination :total="totalNumber" :per-page="limit" :value="pageNumber" @change="fetchData"/>
        </div>
    </va-page>
</template>

<script>
    import {get_goods, get_goods_count} from "@/_helpers/goods";
    import Good from "@/components/Good";

    export default {
        name: "Inventory",
        components: {
            "good": Good,
        },
        data() {
            return {
                loading: false,
                descriptionFilter: null,
                goods: [],
                pageNumber: 1,
                totalNumber: 0,
                limit: 12,
                error: null
            }
        },
        created() {
            this.filterChanged()
        },
        methods: {
            filterChanged() {
                this.fetchData({pageNumber: this.pageNumber, perPage: this.limit})
            },
            fetchData(e) {
                this.loading = true;
                get_goods_count(this.$store.state.token, this.descriptionFilter).then(result => {
                    this.totalNumber = result;
                }).catch(reason => {
                    this.notification.warning({
                        title: "Error occured",
                        message: reason,
                    })
                });
                get_goods(this.$store.state.token, this.descriptionFilter, e.pageNumber, e.perPage)
                    .then(result => {
                        this.goods = result;
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
    .good-card {
        margin-bottom: 20px;
    }
</style>