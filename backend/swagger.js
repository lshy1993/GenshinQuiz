// swagger.js
const swaggerUi = require('swagger-ui-express');
const swaggerJsDoc = require('swagger-jsdoc');
const fs = require('fs');
const yaml = require('js-yaml');
const path = require('path');

// 动态加载 YAML 文件
const swaggerFilePath = path.resolve(__dirname, './openapi/api-docs.yaml'); // YAML文件的路径
let swaggerDocument;

try {
  const fileContents = fs.readFileSync(swaggerFilePath, 'utf8');
  swaggerDocument = yaml.load(fileContents); // 将 YAML 文件转换为 JavaScript 对象
} catch (error) {
  console.error('Failed to load Swagger YAML file:', error);
  process.exit(1); // 如果加载失败，终止进程
}

const swaggerOptions = {
  definition: swaggerDocument,
  apis: ['./routes/*.js'], // 配置路由文件
};

const swaggerDocs = swaggerJsDoc(swaggerOptions);

// 通过 Swagger UI 提供 API 文档界面
const setupSwagger = (app) => {
  app.use('/api-docs', swaggerUi.serve, swaggerUi.setup(swaggerDocs));
};

module.exports = setupSwagger;
