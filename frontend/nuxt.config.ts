export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },
  app: {
    baseURL: '/dashboard/',
  },
  devServer: {
    host: 'frontend',
  },
  runtimeConfig: {
    public: {
      wsUrl: 'ws://localhost',
      apiBase: 'http://localhost:8080'
    }
  },
  modules: [
    '@nuxt/eslint',
    '@nuxt/image',
    '@nuxt/scripts',
    '@nuxt/test-utils',
    '@nuxt/ui',
    '@nuxt/test-utils/module',
    '@nuxtjs/i18n',
    'usebootstrap',
    'nuxt-bootstrap-icons'
  ],
  i18n: {
    locales: [
      { code: 'en', name: 'English', file: 'en.json' },
      { code: 'pl', name: 'Polski', file: 'pl.json' }
    ],
    defaultLocale: 'pl',
  },
  sourcemap: {
    server: true,
    client: true,
  },
  css: [
    '~/assets/scss/main.scss'
  ],
  vite: {
    css: {
      preprocessorOptions: {
        scss: {
          additionalData: '@use "~/assets/scss/variables.scss" as *;',
        },
      },
    },
  },
  robots: {
    robotsTxt: false
  }
})