<template>

    <va-page>
        <va-page-header>
            <div slot="breadcrumb">
                <va-breadcrumb separator="/">
                    <va-breadcrumb-item to="/cart">Cart</va-breadcrumb-item>
                </va-breadcrumb>
            </div>

            <div slot="title">
                <span>Cart</span>
            </div>

            <div slot="actions">
            </div>

            <div slot="bottom">
            </div>
        </va-page-header>

        <div>
            <va-row :gutter="15">
                <va-column :xs="12" :sm="6" :md="4">
                    <va-table hover size="lg">
                        <table>
                            <thead>
                            <tr>
                                <th>Name</th>
                                <th>Price</th>
                                <th>Amount</th>
                                <th>Total</th>
                            </tr>
                            </thead>
                            <tbody>
                            <tr v-for="(item, key, index) in $store.state.cart">
                                <td>{{item.description}}</td>
                                <td>{{item.price}}</td>
                                <td>{{item.amount}}</td>
                                <td>{{item.price*item.amount}}</td>
                            </tr>
                            <tr v-if="Object.keys($store.state.cart).length !== 0">
                                <td>All</td>
                                <td></td>
                                <td></td>
                                <td>{{$store.state.cart_total}}</td>
                            </tr>
                            </tbody>
                        </table>
                    </va-table>
                </va-column>
                <va-column :xs="12" :sm="6" :md="4">
                    <va-form type="vertical">
                        <va-form-item label="Country">
                            <va-select placeholder="Select destination country" search v-model="country"
                                       :options="countries"/>
                        </va-form-item>
                        <va-form-item>
                            <va-button type="primary"
                                       :disabled="country === '' || Object.keys($store.state.cart).length === 0"
                                       :loading="processingInvoice" @click="createInvoice">Make order
                            </va-button>
                            <va-button style="margin-left: 10px" @click="clearCart">Clear cart</va-button>
                        </va-form-item>
                    </va-form>
                </va-column>
            </va-row>
        </div>
    </va-page>
</template>

<script>
    import {get_countries} from "@/_helpers/countries";
    import {create_order} from "@/_helpers/orders";

    export default {
        name: "CartPage",
        data() {
            return {
                country: '',
                countries: [],
                processingInvoice: false,
            }
        },
        created() {
            this.fetchData()
        },
        methods: {
            fetchData() {
                get_countries(this.$store.state.token)
                    .then(result => {
                        this.countries = result.map((v) => {
                            return {value: v.readable_name, label: v.readable_name}
                        });
                    })
                    .catch(reason => {
                        this.notification.warning({
                            title: "Error occured",
                            message: reason,
                        })
                    })
            },
            createInvoice() {
                this.processingInvoice = true;
                let parts = [];
                for (let [key, item] of Object.entries(this.$store.state.cart)) {
                    parts.push({quantity: item.amount, good_id: item.id})
                }
                create_order(this.$store.state.token, {country: this.country}, parts)
                    .then(() => {
                        this.clearCart();
                    })
                    .catch(reason => {
                        this.notification.warning({
                            title: "Error occured",
                            message: reason,
                        })
                    })
                    .finally(() => {
                        this.processingInvoice = false;
                    })
            },
            clearCart() {
                this.$store.commit("clear_cart");
                this.$router.push("/inventory")
            }
        },

    }
</script>

<style scoped>

</style>