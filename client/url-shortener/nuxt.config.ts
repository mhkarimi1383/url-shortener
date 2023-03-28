// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  modules: [
    '@vuestic/nuxt',
  ],
  runtimeConfig: {
    public: {
      baseURL: process.env.BASE_URL || "http://127.0.0.1:8080",
    },
  },
  vuestic: {
    config: {
      // Vuestic config here
      colors: {
        currentPresetName: "dark",
      },
    },
  },
})
