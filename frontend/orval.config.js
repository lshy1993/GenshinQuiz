module.exports = {
  api: {
    input: './openapi/openapi.yaml', // OpenAPI YAML 文件路径
    output: {
      mode: 'split', // 按 operationId 拆分生成文件
      target: './src/api/', // 生成的 TypeScript 客户端代码路径
      schemas: './src/api/dto/', // 生成的模式文件路径
      client: 'swr', // 使用 SWR 作为 HTTP 客户端
      mock: false, // 不生成 mock 数据
      override: {
        useDates: true,
        mutator: {
          path: './src/api/fetcher/fetcher.ts',
          name: 'Fetcher',
        },
        operationName: (operation, route, verb) => {
          // 转换操作 ID 为驼峰命名
          const operationId = operation.operationId;
          if (operationId) {
            return operationId.replace(/-./g, (match) => match.charAt(1).toUpperCase());
          }
          return operation.operationId;
        },
      },
    },
  },
};
