export default defineEventHandler(async (event) => {
  const body = await readBody(event);
  
  // 模拟数据库连接测试
  const { host, port, name, username, password } = body;
  
  // 简单的参数验证
  if (!host || !port || !name || !username) {
    return {
      code: 1,
      message: '请填写完整的数据库连接信息',
      data: {
        success: false,
        message: '请填写完整的数据库连接信息'
      }
    };
  }
  
  // 模拟连接测试逻辑
  if (host === 'invalid-host') {
    return {
      code: 1,
      message: '数据库连接失败',
      data: {
        success: false,
        message: '无法连接到数据库服务器'
      }
    };
  }
  
  // 模拟成功响应
  return {
    code: 0,
    message: '数据库连接测试成功',
    data: {
      success: true,
      message: '数据库连接正常'
    }
  };
});