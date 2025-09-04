<script lang="ts" setup>
import type { VbenFormSchema } from '@vben/common-ui';
import type { Recordable } from '@vben/types';

import { computed, ref } from 'vue';
import { useRouter } from 'vue-router';

import { AuthenticationRegister, z } from '@vben/common-ui';
import { $t } from '@vben/locales';

defineOptions({ name: 'Install' });

const router = useRouter();
const loading = ref(false);

const formSchema = computed((): VbenFormSchema[] => {
    return [
        {
            component: 'VbenInput',
            componentProps: {
                placeholder: '请输入站点名称',
            },
            fieldName: 'siteName',
            label: '站点名称',
            rules: z
                .string()
                .min(1, { message: '请输入站点名称' })
                .default('FastCDN'),
        },
        {
            component: 'VbenInput',
            componentProps: {
                placeholder: '请输入站点URL',
            },
            fieldName: 'siteUrl',
            label: '站点URL',
            rules: z
                .string()
                .min(1, { message: '请输入站点URL' })
                .url({ message: '请输入有效的URL' }),
        },
        {
            component: 'VbenInput',
            componentProps: {
                placeholder: '请输入管理员用户名',
            },
            fieldName: 'adminUsername',
            label: '管理员用户名',
            rules: z
                .string()
                .min(1, { message: '请输入管理员用户名' }),
        },
        {
            component: 'VbenInputPassword',
            componentProps: {
                passwordStrength: true,
                placeholder: '请输入管理员密码',
            },
            fieldName: 'adminPassword',
            label: '管理员密码',
            renderComponentContent() {
                return {
                    strengthText: () => '密码强度',
                };
            },
            rules: z
                .string()
                .min(6, { message: '密码至少6位' }),
        },
        {
            component: 'VbenInputPassword',
            componentProps: {
                placeholder: '请确认管理员密码',
            },
            dependencies: {
                rules(values) {
                    const { adminPassword } = values;
                    return z
                        .string({
                            required_error: '请确认密码',
                        })
                        .min(1, { message: '请确认密码' })
                        .refine((value) => value === adminPassword, {
                            message: '两次密码输入不一致',
                        });
                },
                triggerFields: ['adminPassword'],
            },
            fieldName: 'confirmPassword',
            label: '确认密码',
        },
        {
            component: 'VbenSelect',
            componentProps: {
                options: [
                    { label: 'MySQL', value: 'mysql' },
                    { label: 'PostgreSQL', value: 'postgresql' },
                    { label: 'SQLite', value: 'sqlite' },
                ],
                placeholder: '请选择数据库类型',
            },
            fieldName: 'databaseType',
            label: '数据库类型',
            rules: z
                .string()
                .min(1, { message: '请选择数据库类型' })
                .default('mysql'),
        },
        {
            component: 'VbenInput',
            componentProps: {
                placeholder: '请输入数据库主机',
            },
            fieldName: 'databaseHost',
            label: '数据库主机',
            rules: z
                .string()
                .min(1, { message: '请输入数据库主机' })
                .default('localhost'),
        },
        {
            component: 'VbenInputNumber',
            componentProps: {
                placeholder: '请输入数据库端口',
                min: 1,
                max: 65535,
            },
            fieldName: 'databasePort',
            label: '数据库端口',
            rules: z
                .number()
                .min(1, { message: '端口号必须大于0' })
                .max(65535, { message: '端口号不能超过65535' })
                .default(3306),
        },
        {
            component: 'VbenInput',
            componentProps: {
                placeholder: '请输入数据库名称',
            },
            fieldName: 'databaseName',
            label: '数据库名称',
            rules: z
                .string()
                .min(1, { message: '请输入数据库名称' })
                .default('fastcdn'),
        },
        {
            component: 'VbenInput',
            componentProps: {
                placeholder: '请输入数据库用户名',
            },
            fieldName: 'databaseUsername',
            label: '数据库用户名',
            rules: z
                .string()
                .min(1, { message: '请输入数据库用户名' }),
        },
        {
            component: 'VbenInputPassword',
            componentProps: {
                placeholder: '请输入数据库密码',
            },
            fieldName: 'databasePassword',
            label: '数据库密码',
            rules: z
                .string()
                .min(1, { message: '请输入数据库密码' }),
        },
    ];
});

async function handleInstall(values: Recordable<any>) {
    if (loading.value) return;
    
    loading.value = true;
    
    try {
        // 这里应该调用安装API
        console.log('Installing with data:', values);
        
        // 模拟安装过程
        await new Promise(resolve => setTimeout(resolve, 2000));
        
        // 安装完成后跳转到登录页面
        router.push('/auth/login');
    } catch (error) {
        console.error('Installation failed:', error);
    } finally {
        loading.value = false;
    }
}
</script>

<template>
    <AuthenticationRegister
        :form-schema="formSchema"
        :loading="loading"
        :submit-button-options="{
            content: $t('page.auth.installButton'),
        }"
        @submit="handleInstall"
    />
</template>