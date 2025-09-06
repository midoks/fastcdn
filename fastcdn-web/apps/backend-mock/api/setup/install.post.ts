export default defineEventHandler(async (event) => {
  const body = await readBody(event);
  
  // 模拟系统安装
  const {
    apiHost, apiPort, apiProtocol, apiType,
    nodeId, secret,
    databaseHost, databasePort, databaseName, databaseUsername, databasePassword,
    adminUsername, adminPassword, adminEmail
  } = body;
  
  // 简单的参数验证
  if (!apiHost || !apiPort || !databaseHost || !databasePort || !databaseName || !databaseUsername || !adminUsername || !adminPassword) {
    return {
      code: 1,
      message: '请填写完整的安装信息',
      data: {
        success: false,
        message: '请填写完整的安装信息'
      }
    };
  }
  
  // 模拟管理员用户名冲突检查
  if (adminUsername === 'root') {
    return {
      code: 1,
      message: '系统安装失败',
      data: {
        success: false,
        message: '管理员用户名不能为root，请选择其他用户名'
      }
    };
  }
  
  // 模拟密码强度检查
  if (adminPassword.length < 6) {
    return {
      code: 1,
      message: '系统安装失败',
      data: {
        success: false,
        message: '管理员密码长度不能少于6位'
      }
    };
  }
  
  // 模拟安装过程延迟
  await new Promise(resolve => setTimeout(resolve, 2000));
  
  // 模拟成功响应
  return {
    code: 0,
    message: '系统安装成功',
    data: {
      success: true,
      message: 'FastCDN系统安装完成，请使用管理员账号登录'
    }
  };
});