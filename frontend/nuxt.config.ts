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
      wsUrl: process.env.NUXT_PUBLIC_WS_URL || 'ws://localhost',
      apiBase: process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:8080'
    }
  },
  modules: [
    '@nuxt/eslint',
    '@nuxt/image',
    '@nuxt/scripts',
    '@nuxt/ui',
    '@nuxtjs/i18n',
    'usebootstrap',
    'nuxt-bootstrap-icons',
    ...(process.env.NODE_ENV !== 'production' ? [
      '@nuxt/test-utils',
      '@nuxt/test-utils/module'
    ] : [])
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