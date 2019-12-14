<template>
    <va-card :elevation="elevation" :padding="padding" class="card">
        <div class="title">
            <h2>Orders</h2>
            <va-badge type="info" class="counter" margin="8px">{{orders.length}}</va-badge>
            <va-loading v-if="loading" color="blue" size="sm" class="spinner"/>
        </div>
        <va-collapse :accordion="false">
            <va-card v-for="o in orders">
                <div>
                    <va-lozenge class="lozenge" type="primary" :uppercase="true">order {{o.id}}</va-lozenge>
                    <va-lozenge class="lozenge" type="default" :uppercase="true">{{o.date}}</va-lozenge>
                    <va-lozenge class="lozenge" type="success" :uppercase="true">{{o.status}}</va-lozenge>
                </div>
            <va-collapse-panel header="Details">
                <va-table :hover="true">
                    <table>
                        <thead>
                        <tr>
                            <th>Name</th>
                            <th>Amount</th>
                            <th>Price</th>
                        </tr>
                        </thead>
                        <tbody>
                        <tr v-for="p in o.invoice_parts">
                            <td><router-link :to="`/inventory/${p.good.id}`">{{p.good.description}}</router-link></td>
                            <td>{{p.quantity}}</td>
                            <td>{{p.good.price}}&#163;</td>
                        </tr>
                        </tbody>
                    </table>
                </va-table>
                <div>
                    <h3>Total: {{o.total_price.toFixed(2)}}&#163;</h3>
                </div>
            </va-collapse-panel>
            </va-card>
        </va-collapse>
    </va-card>
</template>
<script>
    import {get_orders} from "@/_helpers/orders";

    export default {
        name: 'user-orders',
        data() {
            return {
                loading: false,
                orders: [],
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
                get_orders(this.username, this.$store.state.token)
                    .then((data) => {
                        console.log(data);
                        this.orders = data.map((v, i, _) => {
                            v.date = new Date(v.date*1000).toDateString();
                            v.total_price = v.invoice_parts.reduce((pre, curr, _) => {
                                return Number.parseFloat(curr.good.price) + pre
                            }, 0);
                            return v;
                        });
                        console.log(this.orders);
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
</style>