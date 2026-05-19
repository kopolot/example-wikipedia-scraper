import { defineConfig } from 'vitest/config'
import { defineVitestProject } from '@nuxt/test-utils/config'

export default defineConfig({
    root: __dirname,
    test: {
        testTimeout: 10000, // globalny timeout 10s
        projects: [
            {
                test: {
                    name: 'unit',
                    include: ['test/{e2e,unit}/**/*.{test,spec}.ts'],
                    environment: 'node',
                },
            },
            await defineVitestProject({
                test: {
                    name: 'nuxt',
                    include: ['test/nuxt/**/*.{test,spec}.ts'],
                    environment: 'nuxt',
                    testTimeout: 20000, // timeout 20s dla testów Nuxt
                },
            }),
        ],
    },
})
