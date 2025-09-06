export default defineEventHandler(async (event) => {
  const body = await readBody(event);
  
  // 模拟API节点创建
  const { host, port, protocol } = body;
  
  // 简单的参数验证
  if (!host || !port || !protocol) {
    return {
      code: 1,
      message: '请填写完整的API节点信息',
      data: {
        success: false,
        message: '请填写完整的API节点信息'
      }
    };
  }
  
  // 模拟端口冲突检查
  if (port === 80 || port === 443) {
    return {
      code: 1,
      message: 'API节点创建失败',
      data: {
        success: false,
        message: '端口已被占用，请选择其他端口'
      }
    };
  }
  
  // 模拟成功响应
  return {
    code: 0,
    message: 'API节点创建成功',
    data: {
      success: true,
      message: 'API节点已成功创建并启动',
      nodeId: `node_${Date.now()}`,
      secret: `secret_${Math.random().toString(36).substring(2)}`
    }
  };
});