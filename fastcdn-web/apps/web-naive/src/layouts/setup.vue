<script lang="ts" setup>
import { computed, ref } from 'vue';
import { useRouter } from 'vue-router';

import { $t } from '@vben/locales';

defineOptions({ name: 'Setup' });

const router = useRouter();
const loading = ref(false);
const currentStep = ref(1);
const totalSteps = 5;

// 表单数据
const formData = ref({
    // Step 2: API配置
    apiHost: 'localhost',
    apiPort: 8080,
    apiProtocol: 'http',
    apiType: 'new',
    
    // Step 3: 数据库配置
    databaseHost: 'localhost',
    databasePort: 3306,
    databaseName: 'fastcdn',
    databaseUsername: '',
    databasePassword: '',
    
    // Step 4: 管理员配置
    adminUsername: '',
    adminPassword: '',
    confirmPassword: '',
    adminEmail: ''
});



// 步骤配置
const steps = computed(() => [
    { key: 'introduction', title: $t('page.auth.setup.steps.introduction') },
    { key: 'apiConfig', title: $t('page.auth.setup.steps.apiConfig') },
    { key: 'database', title: $t('page.auth.setup.steps.database') },
    { key: 'admin', title: $t('page.auth.setup.steps.admin') },
    { key: 'complete', title: $t('page.auth.setup.steps.complete') }
]);




// 下一步
function nextStep() {
    if (currentStep.value < totalSteps) {
        currentStep.value++;
    }
}

// 上一步
function prevStep() {
    if (currentStep.value > 1) {
        currentStep.value--;
    }
}

// 测试数据库连接
function testDatabaseConnection() {
    // 这里应该调用API测试数据库连接
    console.log('Testing database connection...', {
        host: formData.value.databaseHost,
        port: formData.value.databasePort,
        name: formData.value.databaseName,
        username: formData.value.databaseUsername,
    });
}

// 完成安装
async function completeInstallation() {
    if (loading.value) return;
    
    loading.value = true;
    
    try {
        // 这里应该调用安装API
        console.log('Installing with data:', formData.value);
        
        // 模拟安装过程
        await new Promise(resolve => setTimeout(resolve, 2000));
        
        // 安装完成，进入最后一步
        currentStep.value = totalSteps;
    } catch (error) {
        console.error('Installation failed:', error);
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
                    {{ $t('page.auth.setup.title') }}
                </h2>
                <p style="color: #666; font-size: 14px;">
                    {{ $t('page.auth.setup.description') }}
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
                        <p style="margin-bottom: 15px;">{{ $t('page.auth.setup.step1.welcome') }}</p>
                        <p style="margin-bottom: 15px;">{{ $t('page.auth.setup.step1.requirements') }}</p>
                        
                        <div style="background-color: #fff3cd; border: 1px solid #ffeaa7; border-radius: 4px; padding: 15px; margin-bottom: 15px;">
                            <p style="font-size: 14px; margin: 0;">{{ $t('page.auth.setup.step1.disclaimer') }}</p>
                        </div>
                        
                        <div style="background-color: #d1ecf1; border: 1px solid #bee5eb; border-radius: 4px; padding: 15px; margin-bottom: 20px;">
                            <p style="font-size: 14px; margin: 0;">{{ $t('page.auth.setup.step1.agreement') }}</p>
                        </div>
                    </div>
                    
                    <div style="text-align: right;">
                        <button
                            @click="nextStep"
                            style="background-color: #007bff; color: white; padding: 10px 20px; border: none; border-radius: 4px; cursor: pointer; font-size: 14px;"
                            @mouseover="$event.target.style.backgroundColor='#0056b3'"
                            @mouseout="$event.target.style.backgroundColor='#007bff'"
                        >
                            {{ $t('page.auth.setup.buttons.start') }}
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
                                    <td style="width: 140px;padding: 12px; border: 1px solid #ddd; color: #333; font-size: 14px; background-color: #f8f9fa;">{{ $t('page.auth.setup.step2.api_type') }}</td>
                                    <td style="padding: 8px; border: 1px solid #ddd;">
                                        <select v-model="formData.apiType" style="width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px; font-size: 14px;">
                                            <option value="new">自动启动新API节点</option>
                                            <option value="old">使用已安装节点</option>
                                        </select>
                                    </td>
                                </tr>
                                <tr>
                                    <td style="padding: 12px; border: 1px solid #ddd; color: #333; font-size: 14px; background-color: #f8f9fa;">{{ $t('page.auth.setup.step2.protocol') }}</td>
                                    <td style="padding: 8px; border: 1px solid #ddd;">
                                        <select v-model="formData.apiProtocol" style="width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px; font-size: 14px;">
                                            <option value="http">HTTP</option>
                                            <option value="https">HTTPS</option>
                                        </select>
                                    </td>
                                </tr>
                                <tr>
                                    <td style="padding: 12px; border: 1px solid #ddd; color: #333; font-size: 14px; background-color: #f8f9fa;">{{ $t('page.auth.setup.step2.host') }}</td>
                                    <td style="padding: 8px; border: 1px solid #ddd;">
                                        <input v-model="formData.apiHost" type="text" style="width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px; font-size: 14px;" />
                                    </td>
                                </tr>
                                <tr>
                                    <td style="padding: 12px; border: 1px solid #ddd; color: #333; font-size: 14px; background-color: #f8f9fa;">{{ $t('page.auth.setup.step2.port') }}</td>
                                    <td style="padding: 8px; border: 1px solid #ddd;">
                                        <input v-model="formData.apiPort" type="number" style="width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px; font-size: 14px;" />
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
                            {{ $t('page.auth.setup.buttons.previous') }}
                        </button>
                        <button
                            @click="nextStep"
                            style="background-color: #007bff; color: white; padding: 10px 20px; border: none; border-radius: 4px; cursor: pointer; font-size: 14px;"
                            @mouseover="$event.target.style.backgroundColor='#0056b3'"
                            @mouseout="$event.target.style.backgroundColor='#007bff'"
                        >
                            {{ $t('page.auth.setup.buttons.next') }}
                        </button>
                    </div>
                </div>

                <!-- Step 3: 数据库配置 -->
                <div v-if="currentStep === 3" style="margin-bottom: 20px;">
                    <div style="text-align: center; margin-bottom: 20px;">
                        <p style="color: #666; font-size: 14px;">{{ $t('page.auth.setup.step3.description') }}</p>
                    </div>
                    
                    <div style="max-width: 400px; margin: 0 auto 20px;">
                        <div style="margin-bottom: 15px;">
                            <label style="display: block; margin-bottom: 5px; color: #333; font-size: 14px;">{{ $t('page.auth.setup.step3.host') }}</label>
                            <input v-model="formData.databaseHost" type="text" style="width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px; font-size: 14px;" />
                        </div>
                        <div style="margin-bottom: 15px;">
                            <label style="display: block; margin-bottom: 5px; color: #333; font-size: 14px;">{{ $t('page.auth.setup.step3.port') }}</label>
                            <input v-model="formData.databasePort" type="number" style="width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px; font-size: 14px;" />
                        </div>
                        <div style="margin-bottom: 15px;">
                            <label style="display: block; margin-bottom: 5px; color: #333; font-size: 14px;">{{ $t('page.auth.setup.step3.database') }}</label>
                            <input v-model="formData.databaseName" type="text" style="width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px; font-size: 14px;" />
                        </div>
                        <div style="margin-bottom: 15px;">
                            <label style="display: block; margin-bottom: 5px; color: #333; font-size: 14px;">{{ $t('page.auth.setup.step3.username') }}</label>
                            <input v-model="formData.databaseUsername" type="text" style="width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px; font-size: 14px;" />
                        </div>
                        <div style="margin-bottom: 15px;">
                            <label style="display: block; margin-bottom: 5px; color: #333; font-size: 14px;">{{ $t('page.auth.setup.step3.password') }}</label>
                            <input v-model="formData.databasePassword" type="password" style="width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px; font-size: 14px;" />
                        </div>
                        
                        <div style="margin-top: 15px;">
                            <button
                                @click="testDatabaseConnection"
                                style="width: 100%; background-color: #28a745; color: white; padding: 10px; border: none; border-radius: 4px; cursor: pointer; font-size: 14px;"
                                @mouseover="$event.target.style.backgroundColor='#218838'"
                                @mouseout="$event.target.style.backgroundColor='#28a745'"
                            >
                                {{ $t('page.auth.setup.step3.testConnection') }}
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
                            {{ $t('page.auth.setup.buttons.previous') }}
                        </button>
                        <button
                            @click="nextStep"
                            style="background-color: #007bff; color: white; padding: 10px 20px; border: none; border-radius: 4px; cursor: pointer; font-size: 14px;"
                            @mouseover="$event.target.style.backgroundColor='#0056b3'"
                            @mouseout="$event.target.style.backgroundColor='#007bff'"
                        >
                            {{ $t('page.auth.setup.buttons.next') }}
                        </button>
                    </div>
                </div>

                <!-- Step 4: 管理员配置 -->
                <div v-if="currentStep === 4" style="margin-bottom: 20px;">
                    <div style="text-align: center; margin-bottom: 20px;">
                        <p style="color: #666; font-size: 14px;">{{ $t('page.auth.setup.step4.description') }}</p>
                    </div>
                    
                    <div style="max-width: 400px; margin: 0 auto 20px;">
                        <div style="margin-bottom: 15px;">
                            <label style="display: block; margin-bottom: 5px; color: #333; font-size: 14px;">{{ $t('page.auth.setup.step4.username') }}</label>
                            <input v-model="formData.adminUsername" type="text" style="width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px; font-size: 14px;" />
                        </div>
                        <div style="margin-bottom: 15px;">
                            <label style="display: block; margin-bottom: 5px; color: #333; font-size: 14px;">{{ $t('page.auth.setup.step4.password') }}</label>
                            <input v-model="formData.adminPassword" type="password" style="width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px; font-size: 14px;" />
                        </div>
                        <div style="margin-bottom: 15px;">
                            <label style="display: block; margin-bottom: 5px; color: #333; font-size: 14px;">{{ $t('page.auth.setup.step4.confirmPassword') }}</label>
                            <input v-model="formData.confirmPassword" type="password" style="width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px; font-size: 14px;" />
                        </div>
                        <div style="margin-bottom: 15px;">
                            <label style="display: block; margin-bottom: 5px; color: #333; font-size: 14px;">{{ $t('page.auth.setup.step4.email') }}</label>
                            <input v-model="formData.adminEmail" type="email" style="width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px; font-size: 14px;" />
                        </div>
                    </div>
                    
                    <div style="display: flex; justify-content: space-between;">
                        <button
                            @click="prevStep"
                            style="background-color: #6c757d; color: white; padding: 10px 20px; border: none; border-radius: 4px; cursor: pointer; font-size: 14px;"
                            @mouseover="$event.target.style.backgroundColor='#545b62'"
                            @mouseout="$event.target.style.backgroundColor='#6c757d'"
                        >
                            {{ $t('page.auth.setup.buttons.previous') }}
                        </button>
                        <button
                            @click="completeInstallation"
                            :disabled="loading"
                            style="background-color: #007bff; color: white; padding: 10px 20px; border: none; border-radius: 4px; cursor: pointer; font-size: 14px;"
                            @mouseover="$event.target.style.backgroundColor='#0056b3'"
                            @mouseout="$event.target.style.backgroundColor='#007bff'"
                        >
                            <span v-if="loading">安装中...</span>
                            <span v-else>{{ $t('page.auth.setup.buttons.install') }}</span>
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
                            {{ $t('page.auth.setup.step5.success') }}
                        </h3>
                        <p style="color: #666; font-size: 14px;">{{ $t('page.auth.setup.step5.description') }}</p>
                    </div>
                    
                    <button
                        @click="goToLogin"
                        style="background-color: #007bff; color: white; padding: 12px 32px; border: none; border-radius: 4px; cursor: pointer; font-size: 16px;"
                        @mouseover="$event.target.style.backgroundColor='#0056b3'"
                        @mouseout="$event.target.style.backgroundColor='#007bff'"
                    >
                        {{ $t('page.auth.setup.step5.loginButton') }}
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
/* 自定义样式 */
</style>