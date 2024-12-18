module.exports = {
  api: {
    input: '../openapi/api-docs.yaml', // OpenAPI YAML 文件路径
    output: {
      mode: 'tags', // 按照 tag 生成文件（每个 API 会生成独立的文件）
      target: './src/api/', // 生成的客户端代码存放路径
      client: 'axios', // 使用 axios 作为 HTTP 客户端
    },
  },
};
