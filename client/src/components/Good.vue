<template>
    <va-card elevation="3">
        <h3>{{description}}</h3>
        <p>Amount: {{amount}}</p>
        <va-button>Add to cart</va-button>
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
            }
        },
    }
</script>

<style scoped>

</style>