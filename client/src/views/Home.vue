<template>
    <va-page>
        <va-page-header>
            <div slot="title">
                <span>Home</span>
            </div>
        </va-page-header>
        <h2>Hi!</h2>
        <p>The page is a summary of the project.</p>
        <p>Currently the following was implemented:</p>
        <ul>
            <li><router-link to="/auth#signin">Log in</router-link></li>
            <li><router-link to="/auth#signup">Sign up</router-link></li>
            <li><router-link to="/user/paulra">Account</router-link></li>
            <li><router-link to="/home">Home</router-link></li>
            <li><router-link to="/complaints">Complaints</router-link></li>
            <li><router-link to="/inventory">Inventory</router-link></li>
        </ul>
    </va-page>
</template>

<script>
    export default {
        name: 'home'
    }
</script>
