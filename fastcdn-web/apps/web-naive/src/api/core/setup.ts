import { baseRequestClient } from '#/api/request';

export namespace SetupApi {
  /** 数据库连接测试参数 */
  export interface DatabaseTestParams {
    host: string;
    port: number;
    name: string;
    username: string;
    password: string;
  }

  /** 数据库连接测试返回值 */
  export interface DatabaseTestResult {
    success: boolean;
    message: string;
  }

  /** 系统安装参数 */
  export interface InstallParams {
    // API配置
    apiHost: string;
    apiPort: number;
    apiProtocol: string;
    apiType: string;
    nodeId?: string;
    secret?: string;
    
    // 数据库配置
    databaseHost: string;
    databasePort: number;
    databaseName: string;
    databaseUsername: string;
    databasePassword: string;
    
    // 管理员配置
    adminUsername: string;
    adminPassword: string;
    adminEmail: string;
  }

  /** 系统安装返回值 */
  export interface InstallResult {
    success: boolean;
    message: string;
  }
}

/**
 * 测试数据库连接
 */
export async function testDatabaseConnectionApi(data: SetupApi.DatabaseTestParams) {
  console.log(data);
  return baseRequestClient.post<SetupApi.DatabaseTestResult>('/setup/db_test', data);
}

/**
 * 系统安装
 */
export async function installSystemApi(data: SetupApi.InstallParams) {
  return baseRequestClient.post<SetupApi.InstallResult>('/setup/install', data);
}

/**
 * 创建API节点
 */
export async function createApiNodeApi(data: { host: string; port: number; protocol: string }) {
  return baseRequestClient.post('/setup/create-api-node', data);
}