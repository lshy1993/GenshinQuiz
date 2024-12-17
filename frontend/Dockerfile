# 使用官方的 Bun 镜像
FROM oven/bun:latest

# 设置工作目录
WORKDIR /app

# 复制 package.json 和其他依赖文件
COPY package.json bun.lockb ./

# 使用 Bun 安装依赖
RUN bun install

# 复制应用代码
COPY . .

# 启动开发服务器
CMD ["bun", "run", "dev"]