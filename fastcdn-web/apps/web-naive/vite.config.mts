import { defineConfig } from '@vben/vite-config';
import type { UserConfig } from 'vite';
import { resolve } from 'node:path';

interface ConfigReturn {
    application: Record<string, any>;
    vite: UserConfig;
}

const config: ReturnType<typeof defineConfig> = defineConfig(async (): Promise<ConfigReturn> => {
    return {
        application: {},
        vite: {
            resolve: {
                alias: {
                    '#': resolve(__dirname, './src'),
                },
            },
            build: {
                rollupOptions: {
                    output: {
                        assetFileNames: (assetInfo) => {
                            const extType = assetInfo.name?.split('.').pop();
                            if (extType === 'css') {
                                return 'static/css/[name]-[hash].[ext]';
                            }
                            return 'static/[name]-[hash].[ext]';
                        },
                        chunkFileNames: 'static/js/[name]-[hash].js',
                        entryFileNames: 'static/js/index-[name]-[hash].js',
                    },
                },
            },
            server: {
                proxy: {
                    '/api': {
                        changeOrigin: true,
                        rewrite: (path: string) => path.replace(/^\/api/, ''),
                        // mock代理目标地址
                        target: 'http://localhost:5320/api',
                        ws: true,
                    },
                },
            },
        },
    };
});

export default config;
