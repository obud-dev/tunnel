# 1、在web目录打包react项目，输出dist目录
# 2、把dist目录复制到server目录下
# 3、在server目录打包golang项目，输出server文件

# 使用 Node.js 进行 React 项目构建
FROM node:20.16.0 AS react-build

WORKDIR /app

# 复制 web 目录到工作目录
COPY ./web /app

# 安装依赖并构建项目
RUN npm install && npm run build

# 使用 Go 进行后端构建
FROM golang:alpine AS golang-build

WORKDIR /app

# 复制 go.mod 和 go.sum 并安装依赖
COPY ./go.mod ./go.sum ./
RUN go mod download

# 复制整个 server 目录到工作目录
COPY ./server ./server

# 复制从前一阶段构建的 React 项目的 dist 目录到 Go 项目
COPY --from=react-build /app/dist ./server/dist

# 进入 server 目录并打包 Go 项目
WORKDIR /app/server
RUN go build -o server .

# 最后的镜像基于alpin容器
FROM alpine

WORKDIR /app

# 将构建好的服务器二进制文件和dist目录复制到alpin容器中
COPY --from=golang-build /app/server/server /app/server
COPY --from=golang-build /app/server/dist /app/dist

# 安装 libc6-compat 以确保与 Go 二进制兼容
RUN apk add --no-cache libc6-compat

# 暴露服务端口，以便外部可以访问
EXPOSE 8000 5429

# 如果需要可以使用以下命令运行服务
CMD ["/app/server"]
