// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  modules: [
    '@vuestic/nuxt'
  ],
  vuestic: {
    config: {
      // Vuestic config here
      colors: {
        currentPresetName: "dark",
      },
    }
  },
})
