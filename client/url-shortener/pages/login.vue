<template>
    <div class="login-form-container">
        <va-form tag="form" @submit.prevent="handleSubmit" class="login-form">
            <va-input v-model="username" label="Username" class="input" />
            <va-input v-model="password" class="mt-3 input" type="password" label="Password" />
            <va-button type="submit" class="mt-3">
                Login
            </va-button>
        </va-form>
    </div>
</template>

<script lang="ts">
export default {
    data() {
        return {
            username: "",
            password: "",
        };
    },
    methods: {
        handleSubmit() {
            alert("-- form submit --");
        },
    },
};
</script>

<style scoped>
body {
    margin: 0;
    overflow: hidden;
}

.login-form-container {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 90vh;
}

.login-form {
    width: 90%;
    max-width: 400px;
    padding: 20px;
    border: 1px solid #ddd;
    border-radius: 4px;
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
}

.input {
    margin-bottom: 10px;
}
</style>
