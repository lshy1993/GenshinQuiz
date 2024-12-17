# 基于官方 Node.js 镜像
FROM node:22

# 设置工作目录
WORKDIR /app

# 复制依赖文件并安装依赖
COPY package.json package-lock.json ./
RUN npm install

# 复制应用代码
COPY . .

# 启动后端服务器
CMD ["node", "index.js"]
