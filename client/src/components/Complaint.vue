<template>
    <va-card :elevation="2" :pading="12" class="container">
        <p>ID: {{complaint.id}}</p>
        <div>
            <p>Body: {{complaint.body}}</p>
            <p>Requester: {{complaint.user.name}}</p>
        </div>
        <input v-if="!complaint.answer" placeholder="type answer">
        <div v-else>
            <p>{{complaint.answer}}</p>
            <p>Reviewer: {{complaint.reviewer.name}}</p>
        </div>
    </va-card>
</template>

<script>
    export default {
        name: "Complaint",
        props: ["complaint"]
    }
</script>

<style scoped>
    .container {
        margin-bottom: 20px;
    }
</style>