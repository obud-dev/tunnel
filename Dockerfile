# 1、在web目录打包react项目，输出dist目录
# 2、把dist目录复制到server目录下
# 3、在server目录打包golang项目，输出server文件

# 使用 Node.js 进行 React 项目构建
FROM node:20.16.0 as react-build

WORKDIR /app

# 复制 web 目录到工作目录
COPY ./web /app

# 安装依赖并构建项目
RUN pnpm install && pnpm run build

# 使用 Go 进行后端构建
FROM golang:1.22.5 as golang-build

WORKDIR /app/server

# 复制 server 目录到工作目录
COPY ./server /app/server

# 复制从前一阶段构建的 React 项目的 dist 目录到 Go 项目
COPY --from=react-build /app/dist /app/server/dist

# 打包 Go 项目
RUN go build -o server .

# 如果需要可以使用以下命令运行服务
CMD ["./server"]