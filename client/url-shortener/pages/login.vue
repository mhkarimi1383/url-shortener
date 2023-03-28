<template id="login-page">
    <div>
        <va-form tag="form" @submit.prevent="handleSubmit">
            <va-input class="form-item first-item" v-model="username" label="Username" />
            <va-input class="form-item" v-model="password" type="password" label="Password" />
            <va-button class="form-item last-item" type="submit">
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
        async handleSubmit() {
            const config = useRuntimeConfig();
            const { data } = await useFetch(config.public.baseURL + "/api/login/", {
                method: "POST",
                body: {
                    "name": this.username,
                    "password": this.password,
                }
            });
            console.log("token");
            console.log(data.value.token);
        },
    },
};
</script>

<style scoped>
form {
    display: flex;
    align-items: center;
    height: 100vh;
    flex-direction: column;
    margin-left: 20%;
    margin-right: 20%;
}


/* I Don't know how, But it works */
.last-item {
    margin-bottom: 30%;
}

.first-item {
    margin-top: 30%;
}

@media only screen and (max-width : 768px) {
    .last-item {
        margin-bottom: 100%;
    }

    .first-item {
        margin-top: 100%;
    }
}

.form-item {
    width: 100%;
}
</style>
