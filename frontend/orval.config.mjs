export default {
  api: {
    input: './openapi/api-docs.yaml', // OpenAPI YAML 文件路径
    output: {
      mode: 'operationId', // 按 operationId 生成文件
      target: './src/api', // 生成的 TypeScript 客户端代码路径
      client: 'axios', // 使用 Axios 作为 HTTP 客户端
      transformOperationId: (operationId) => {
        return operationId.replace(/-./g, (match) => match.charAt(1).toUpperCase());
      }, // 转换操作 ID 为驼峰命名
    },
  },
};
