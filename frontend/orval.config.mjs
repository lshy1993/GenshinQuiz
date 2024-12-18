export default {
  api: {
    input: './openapi/api-docs.yaml', // OpenAPI YAML 文件路径
    output: {
      mode: 'tags', // 按 tags 生成文件
      target: './src/api', // 生成的 TypeScript 客户端代码路径
      client: 'axios', // 使用 Axios 作为 HTTP 客户端
    },
  },
};
