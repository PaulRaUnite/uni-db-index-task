<template>
    <va-card elevation="3">
        <va-button :disabled="$store.state.token === null" slot="topRight" icon-before="cart-plus" @click="() => this.$refs.modal.open()">Add to cart</va-button>
        <va-modal title="Add to cart" ref="modal" @confirm="addToCart">
            <div slot="body">
                <p>Will be added: {{amountToAdd}}</p>
                <va-range step="1" min="0" :max="amount" v-model="amountToAdd"/>
            </div>
        </va-modal>
        <h3>{{description}}</h3>
        <p>Amount: {{amount}}</p>
    </va-card>
</template>

<script>
    export default {
        name: "Good",
        props: {
            id: {
                type: Number,
                required: true,
            },
            description: {
                type: String,
                required: true,
            },
            amount: {
                type: Number,
                required: true,
            },
            price: {
                type: Number,
                required: true
            },
        },
        data() {
            return {
                amountToAdd: 0,
            }
        },
        methods: {
            addToCart() {
                if (this.amountToAdd !== 0) {
                    this.$store.commit("add_to_cart", {
                        id: this.id,
                        price: this.price,
                        amount: this.amountToAdd,
                        description: this.description
                    })
                }
                this.$refs.modal.close()
            }
        }
    }
</script>

<style scoped>

</style>