<script lang="ts" setup>
import { computed, ref, watch } from 'vue';
import { useRouter } from 'vue-router';

import { $t } from '@vben/locales';
import { message } from '#/adapter/naive';
import { testDatabaseConnectionApi, installSystemApi, createApiNodeApi } from '#/api/core/setup';

defineOptions({ name: 'Setup' });

const router = useRouter();
const loading = ref(false);
const currentStep = ref(1);
const totalSteps = 5;

// 表单数据
const formData = ref({
    // Step 2: API配置
    apiHost: '127.0.0.1',
    apiPort: 10001,
    apiProtocol: 'http',
    apiType: 'new',
    nodeId: '',
    secret: '',
    
    // Step 3: 数据库配置
    databaseHost: '127.0.0.1',
    databasePort: 3306,
    databaseName: 'fastcdn',
    databaseUsername: '',
    databasePassword: '',
    
    // Step 4: 管理员配置
    adminUsername: 'admin',
    adminPassword: '',
    confirmPassword: '',
    adminEmail: ''
});

// 监听apiType变化，当选择"使用已安装节点"时清空相关字段
watch(() => formData.value.apiType, (newType) => {
    if (newType === 'old') {
        formData.value.nodeId = '';
        formData.value.secret = '';
        formData.value.apiHost = '';
        formData.value.apiPort = '';
    } else if (newType === 'new'){
        if (formData.value.apiPort === ''){
            formData.value.apiPort = "10001";
        }

        if (formData.value.apiHost === ''){
            formData.value.apiHost = "127.0.0.1";
        }
    }
});


// 步骤配置
const steps = computed(() => [
    { key: 'introduction', title: $t('page.setup.steps.introduction') },
    { key: 'apiConfig', title: $t('page.setup.steps.apiConfig') },
    { key: 'database', title: $t('page.setup.steps.database') },
    { key: 'admin', title: $t('page.setup.steps.admin') },
    { key: 'complete', title: $t('page.setup.steps.complete') }
]);


// 下一步
async function nextStep() {
    if (currentStep.value < totalSteps) {
        let plus = true;

        if (currentStep.value === 3){
            let db_t = await testDatabaseConnection();
            if (!db_t){
                plus = false;
            }
        }

        if (plus){
            currentStep.value++;
        }
    }
}

// 上一步
function prevStep() {
    if (currentStep.value > 1) {
        currentStep.value--;
    }
}

// 测试数据库连接
async function testDatabaseConnection() {
    if (loading.value) return false;

    if (formData.value.databaseHost === ""){
        message.error("数据库主机不能为空!");
        return false;
    }
    if (formData.value.databasePort === ""){
        message.error("数据库端口不能为空!");
        return false;
    }

    if (formData.value.databaseName === ""){
        message.error("数据库名称不能为空!");
        return false;
    }
    
    if (formData.value.databaseUsername === ""){
        message.error("数据库用户名不能为空!");
        return false;
    }

    if (formData.value.databasePassword === ""){
        message.error("数据库密码不能为空!");
        return false;
    }

    loading.value = true;
    
    try {
        const result = await testDatabaseConnectionApi({
            hostname: formData.value.databaseHost,
            port: formData.value.databasePort,
            dbname: formData.value.databaseName,
            username: formData.value.databaseUsername,
            password: formData.value.databasePassword,
        });
        
        if (result.data.status===0) {
            message.success('数据库连接测试成功！');
            return true;
        } else {
            message.error(`${result.data.message}`);
        }
    } catch (error: any) {
        console.error('Database connection test failed:', error);
        message.error(`数据库连接测试失败：${error.message || '未知错误'}`);
    } finally {
        loading.value = false;
    }

    return false;
}

// 创建API节点
async function createApiNode() {
    if (formData.value.apiType !== 'new') return;
    
    try {
        await createApiNodeApi({
            host: formData.value.apiHost,
            port: formData.value.apiPort,
            protocol: formData.value.apiProtocol,
        });
        message.success('API节点创建成功！');
    } catch (error: any) {
        console.error('API node creation failed:', error);
        message.error(`API节点创建失败：${error.message || '未知错误'}`);
        throw error;
    }
}

// 完成安装
async function completeInstallation() {
    if (loading.value) return;

    if (formData.value.adminUsername === ""){
        message.error("管理员用户名不能为空!");
        return;
    }

    if (formData.value.adminPassword === ""){
        message.error("管理员密码不能为空!");
        return;
    }
    
    loading.value = true;
    
    try {    
        // 调用安装API
        const result = await installSystemApi({
            api_host: formData.value.apiHost,
            api_port: formData.value.apiPort,
            api_protocol: formData.value.apiProtocol,
            api_type: formData.value.apiType,

            node_id: formData.value.nodeId,
            secret: formData.value.secret,

            hostname: formData.value.databaseHost,
            port: formData.value.databasePort,
            dbname: formData.value.databaseName,
            username: formData.value.databaseUsername,
            password: formData.value.databasePassword,

            admin_username: formData.value.adminUsername,
            admin_password: formData.value.adminPassword,
        });
        
        if (result.success) {
            message.success('系统安装成功！');
            // 安装完成，进入最后一步
            currentStep.value = totalSteps;
        } else {
            message.error(`系统安装失败：${result.data.message}`);
        }
    } catch (error: any) {
        console.error('Installation failed:', error);
        message.error(`系统安装失败：${error.message || '未知错误'}`);
    } finally {
        loading.value = false;
    }
}

// 跳转到登录页面
function goToLogin() {
    router.push('/auth/login');
}
</script>

<template>
    <div style="min-height: 100vh; background-color: #f9f9f9; padding: 20px; font-family: Arial, sans-serif;">
        <div style="max-width: 800px; margin: 0 auto;">
            <!-- 头部 -->
            <div style="text-align: center; margin-bottom: 20px;">
                <h2 style="font-size: 28px; color: #333; margin-bottom: 10px;">
                    {{ $t('page.setup.title') }}
                </h2>
                <p style="color: #666; font-size: 14px;">
                    {{ $t('page.setup.description') }}
                </p>
            </div>

            <!-- 步骤指示器 -->
            <div style="margin-bottom: 20px;">
                <div 
                    :style="{
                        display: 'flex',
                        justifyContent: 'center',
                        alignItems: 'center',
                        gap: '0px',
                        border: '1px solid rgba(34, 36, 38, 0.15)'
                    }"
                >
                    <div v-for="(step, index) in steps" :key="step.key" style="display: flex; align-items: center;">
                        
                        <span
                            :style="{
                                marginLeft: '3px',
                                fontSize: '15px',
                                padding: '10px 15px 10px 15px',
                                color: index + 1 <= currentStep ? '#007bff' : '#666'
                            }"
                        >
                            {{ index + 1 }}. {{ step.title }}
                        </span>
                        <!--
                        <div v-if="index < steps.length - 1" style="width: 40px; height: 2px; background-color: #ddd; margin-left: 15px;"></div>
                        -->
                    </div>
                </div>
            </div>

            <!-- 内容区域 -->
            <div style="background-color: white; border-radius: 2px; padding: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1);">
                <!-- Step 1: 介绍 -->
                <div v-if="currentStep === 1" style="margin-bottom: 20px;">
                    <div style="color: #555; line-height: 1.6;">
                        <p style="margin-bottom: 15px;">{{ $t('page.setup.step1.welcome') }}</p>
                        <p style="margin-bottom: 15px;">{{ $t('page.setup.step1.requirements') }}</p>
                        
                        <div style="background-color: #fff3cd; border: 1px solid #ffeaa7; border-radius: 4px; padding: 15px; margin-bottom: 15px;">
                            <p style="font-size: 14px; margin: 0;">{{ $t('page.setup.step1.disclaimer') }}</p>
                        </div>
                        
                        <div style="background-color: #d1ecf1; border: 1px solid #bee5eb; border-radius: 4px; padding: 15px; margin-bottom: 20px;">
                            <p style="font-size: 14px; margin: 0;">{{ $t('page.setup.step1.agreement') }}</p>
                        </div>
                    </div>
                    
                    <div style="text-align: right;">
                        <button
                            @click="nextStep"
                            style="background-color: #007bff; color: white; padding: 10px 20px; border: none; border-radius: 4px; cursor: pointer; font-size: 14px;"
                            @mouseover="$event.target.style.backgroundColor='#0056b3'"
                            @mouseout="$event.target.style.backgroundColor='#007bff'"
                        >
                            {{ $t('page.setup.buttons.start') }}
                        </button>
                    </div>
                </div>

                <!-- Step 2: API配置 -->
                <div v-if="currentStep === 2" style="margin-bottom: 20px;">
                    <div style="width: 100%; margin: 0 auto 20px;">
                        <table style="width: 100%; border-collapse: collapse; border: 1px solid #ddd;">
                            <!--
                            <thead>
                                <tr style="background-color: #f8f9fa;">
                                    <th style="padding: 12px; text-align: left; border: 1px solid #ddd; color: #333; font-size: 14px; font-weight: bold;">配置项</th>
                                    <th style="padding: 12px; text-align: left; border: 1px solid #ddd; color: #333; font-size: 14px; font-weight: bold;">值</th>
                                </tr>
                            </thead>
                            -->
                            <tbody>
                                <tr>
                                    <td style="width: 140px;padding: 12px; border: 1px solid #ddd; color: #333; font-size: 14px; background-color: #f8f9fa;">{{ $t('page.setup.step2.api_type') }}</td>
                                    <td style="padding: 8px; border: 1px solid #ddd;">
                                        <div style="display: flex; gap: 20px; align-items: center;">
                                            <label style="display: flex; align-items: center; cursor: pointer;">
                                                <input v-model="formData.apiType" type="radio" value="new" style="margin-right: 8px;" />
                                                <span style="font-size: 14px; color: #333;">自动启动新API节点</span>
                                            </label>
                                            <label style="display: flex; align-items: center; cursor: pointer;">
                                                <input v-model="formData.apiType" type="radio" value="old" style="margin-right: 8px;" />
                                                <span style="font-size: 14px; color: #333;">使用已安装节点</span>
                                            </label>
                                        </div>
                                    </td>
                                </tr>
                                <tr v-if="formData.apiType === 'old'">
                                    <td style="padding: 12px; border: 1px solid #ddd; color: #333; font-size: 14px; background-color: #f8f9fa;">{{ $t('page.setup.step2.protocol') }}</td>
                                    <td style="padding: 8px; border: 1px solid #ddd;">
                                        <select v-model="formData.apiProtocol" style="width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px; font-size: 14px;">
                                            <option value="http">HTTP</option>
                                            <option value="https">HTTPS</option>
                                        </select>
                                    </td>
                                </tr>
                                <tr>
                                    <td style="padding: 12px; border: 1px solid #ddd; color: #333; font-size: 14px; background-color: #f8f9fa;">{{ $t('page.setup.step2.host') }}</td>
                                    <td style="padding: 8px; border: 1px solid #ddd;">
                                        <input v-model="formData.apiHost" type="text" style="width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px; font-size: 14px;" />
                                    </td>
                                </tr>
                                <tr>
                                    <td style="padding: 12px; border: 1px solid #ddd; color: #333; font-size: 14px; background-color: #f8f9fa;">{{ $t('page.setup.step2.port') }}</td>
                                    <td style="padding: 8px; border: 1px solid #ddd;">
                                        <input v-model="formData.apiPort" type="number" style="width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px; font-size: 14px;" />
                                    </td>
                                </tr>
                                <tr v-if="formData.apiType === 'old'">
                                    <td style="padding: 12px; border: 1px solid #ddd; color: #333; font-size: 14px; background-color: #f8f9fa;">{{ $t('page.setup.step2.node_id') }}</td>
                                    <td style="padding: 8px; border: 1px solid #ddd;">
                                        <input v-model="formData.nodeId" type="text" style="width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px; font-size: 14px;" />
                                    </td>
                                </tr>
                                <tr v-if="formData.apiType === 'old'">
                                    <td style="padding: 12px; border: 1px solid #ddd; color: #333; font-size: 14px; background-color: #f8f9fa;">{{ $t('page.setup.step2.secret') }}</td>
                                    <td style="padding: 8px; border: 1px solid #ddd;">
                                        <input v-model="formData.secret" type="text" style="width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px; font-size: 14px;" />
                                    </td>
                                </tr>
                            </tbody>
                        </table>
                    </div>
                    
                    <div style="display: flex; justify-content: space-between;">
                        <button
                            @click="prevStep"
                            style="background-color: #6c757d; color: white; padding: 10px 20px; border: none; border-radius: 4px; cursor: pointer; font-size: 14px;"
                            @mouseover="$event.target.style.backgroundColor='#545b62'"
                            @mouseout="$event.target.style.backgroundColor='#6c757d'"
                        >
                            {{ $t('page.setup.buttons.previous') }}
                        </button>
                        <button
                            @click="nextStep"
                            style="background-color: #007bff; color: white; padding: 10px 20px; border: none; border-radius: 4px; cursor: pointer; font-size: 14px;"
                            @mouseover="$event.target.style.backgroundColor='#0056b3'"
                            @mouseout="$event.target.style.backgroundColor='#007bff'"
                        >
                            {{ $t('page.setup.buttons.next') }}
                        </button>
                    </div>
                </div>

                <!-- Step 3: 数据库配置 -->
                <div v-if="currentStep === 3" style="margin-bottom: 20px;">
                    
                    <div style="max-width: 100%; margin: 0 auto 20px;">
                        <table style="width: 100%; border-collapse: collapse; border: 1px solid #ddd; margin-bottom: 20px;">
                            <tbody>
                                <tr>
                                    <td style="width:140px;padding: 12px; border: 1px solid #ddd; color: #333; font-size: 14px; background-color: #f8f9fa;">{{ $t('page.setup.step3.host') }}</td>
                                    <td style="padding: 8px; border: 1px solid #ddd;">
                                        <input v-model="formData.databaseHost" type="text" style="width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px; font-size: 14px;" />
                                    </td>
                                </tr>
                                <tr>
                                    <td style="padding: 12px; border: 1px solid #ddd; color: #333; font-size: 14px; background-color: #f8f9fa;">{{ $t('page.setup.step3.port') }}</td>
                                    <td style="padding: 8px; border: 1px solid #ddd;">
                                        <input v-model="formData.databasePort" type="number" style="width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px; font-size: 14px;" />
                                    </td>
                                </tr>
                                <tr>
                                    <td style="padding: 12px; border: 1px solid #ddd; color: #333; font-size: 14px; background-color: #f8f9fa;">{{ $t('page.setup.step3.database') }}</td>
                                    <td style="padding: 8px; border: 1px solid #ddd;">
                                        <input v-model="formData.databaseName" type="text" style="width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px; font-size: 14px;" />
                                    </td>
                                </tr>
                                <tr>
                                    <td style="padding: 12px; border: 1px solid #ddd; color: #333; font-size: 14px; background-color: #f8f9fa;">{{ $t('page.setup.step3.username') }}</td>
                                    <td style="padding: 8px; border: 1px solid #ddd;">
                                        <input v-model="formData.databaseUsername" type="text" style="width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px; font-size: 14px;" />
                                    </td>
                                </tr>
                                <tr>
                                    <td style="padding: 12px; border: 1px solid #ddd; color: #333; font-size: 14px; background-color: #f8f9fa;">{{ $t('page.setup.step3.password') }}</td>
                                    <td style="padding: 8px; border: 1px solid #ddd;">
                                        <input v-model="formData.databasePassword" type="password" style="width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px; font-size: 14px;" />
                                    </td>
                                </tr>
                            </tbody>
                        </table>
                        
                        <div style="margin-top: 15px;">
                            <button
                                @click="testDatabaseConnection"
                                style="width: 100%; background-color: #28a745; color: white; padding: 10px; border: none; border-radius: 4px; cursor: pointer; font-size: 14px;"
                                @mouseover="$event.target.style.backgroundColor='#218838'"
                                @mouseout="$event.target.style.backgroundColor='#28a745'"
                            >
                                {{ $t('page.setup.step3.testConnection') }}
                            </button>
                        </div>
                    </div>
                    
                    <div style="display: flex; justify-content: space-between;">
                        <button
                            @click="prevStep"
                            style="background-color: #6c757d; color: white; padding: 10px 20px; border: none; border-radius: 4px; cursor: pointer; font-size: 14px;"
                            @mouseover="$event.target.style.backgroundColor='#545b62'"
                            @mouseout="$event.target.style.backgroundColor='#6c757d'"
                        >
                            {{ $t('page.setup.buttons.previous') }}
                        </button>
                        <button
                            @click="nextStep"
                            style="background-color: #007bff; color: white; padding: 10px 20px; border: none; border-radius: 4px; cursor: pointer; font-size: 14px;"
                            @mouseover="$event.target.style.backgroundColor='#0056b3'"
                            @mouseout="$event.target.style.backgroundColor='#007bff'"
                        >
                            {{ $t('page.setup.buttons.next') }}
                        </button>
                    </div>
                </div>

                <!-- Step 4: 管理员配置 -->
                <div v-if="currentStep === 4" style="margin-bottom: 20px;">
                    
                    <div style="max-width: 100%; margin: 0 auto 20px;">
                        <table style="width: 100%; border-collapse: collapse; border: 1px solid #ddd;">
                            <tbody>
                                <tr>
                                    <td style="width:140px;padding: 12px; border: 1px solid #ddd; color: #333; font-size: 14px; background-color: #f8f9fa;">{{ $t('page.setup.step4.username') }}</td>
                                    <td style="padding: 8px; border: 1px solid #ddd;">
                                        <input v-model="formData.adminUsername" type="text" style="width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px; font-size: 14px;" />
                                    </td>
                                </tr>
                                <tr>
                                    <td style="padding: 12px; border: 1px solid #ddd; color: #333; font-size: 14px; background-color: #f8f9fa;">{{ $t('page.setup.step4.password') }}</td>
                                    <td style="padding: 8px; border: 1px solid #ddd;">
                                        <input v-model="formData.adminPassword" type="password" style="width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px; font-size: 14px;" />
                                    </td>
                                </tr>
                                <tr>
                                    <td style="padding: 12px; border: 1px solid #ddd; color: #333; font-size: 14px; background-color: #f8f9fa;">{{ $t('page.setup.step4.confirmPassword') }}</td>
                                    <td style="padding: 8px; border: 1px solid #ddd;">
                                        <input v-model="formData.confirmPassword" type="password" style="width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px; font-size: 14px;" />
                                    </td>
                                </tr>
                            </tbody>
                        </table>
                    </div>
                    
                    <div style="display: flex; justify-content: space-between;">
                        <button
                            @click="prevStep"
                            style="background-color: #6c757d; color: white; padding: 10px 20px; border: none; border-radius: 4px; cursor: pointer; font-size: 14px;"
                            @mouseover="$event.target.style.backgroundColor='#545b62'"
                            @mouseout="$event.target.style.backgroundColor='#6c757d'"
                        >
                            {{ $t('page.setup.buttons.previous') }}
                        </button>
                        <button
                            @click="completeInstallation"
                            :disabled="loading"
                            style="background-color: #007bff; color: white; padding: 10px 20px; border: none; border-radius: 4px; cursor: pointer; font-size: 14px;"
                            @mouseover="$event.target.style.backgroundColor='#0056b3'"
                            @mouseout="$event.target.style.backgroundColor='#007bff'"
                        >
                            <span v-if="loading">安装中...</span>
                            <span v-else>{{ $t('page.setup.buttons.install') }}</span>
                        </button>
                    </div>
                </div>

                <!-- Step 5: 完成安装 -->
                <div v-if="currentStep === 5" style="text-align: center; margin-bottom: 20px;">
                    <div style="width: 64px; height: 64px; background-color: #d4edda; border-radius: 50%; display: flex; align-items: center; justify-content: center; margin: 0 auto 20px;">
                        <svg style="width: 32px; height: 32px; color: #28a745;" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
                        </svg>
                    </div>
                    
                    <div style="margin-bottom: 20px;">
                        <h3 style="font-size: 24px; font-weight: bold; color: #333; margin-bottom: 10px;">
                            {{ $t('page.setup.step5.success') }}
                        </h3>
                        <p style="color: #666; font-size: 14px;">{{ $t('page.setup.step5.description') }}</p>
                    </div>
                    
                    <button
                        @click="goToLogin"
                        style="background-color: #007bff; color: white; padding: 12px 32px; border: none; border-radius: 4px; cursor: pointer; font-size: 16px;"
                        @mouseover="$event.target.style.backgroundColor='#0056b3'"
                        @mouseout="$event.target.style.backgroundColor='#007bff'"
                    >
                        {{ $t('page.setup.step5.loginButton') }}
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
/* 自定义样式 */
</style>