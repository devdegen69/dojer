// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2024-11-01',
  devtools: { enabled: true },
  modules: ['@nuxt/ui', '@vueuse/nuxt', '@unocss/nuxt'],
	runtimeConfig: {
		public: {
			apiUrl: process.env.API_URL || 'http://0.0.0.0:8033',
		}
	}
})
