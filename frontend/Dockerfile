# ---- build ----
FROM node:20-alpine AS builder
WORKDIR /app
# 清空代理环境变量，避免容器内走宿主机 127.0.0.1 代理（不可达）
ENV HTTP_PROXY=
ENV HTTPS_PROXY=
ENV ALL_PROXY=
ENV http_proxy=
ENV https_proxy=
ENV all_proxy=
COPY package*.json ./
RUN npm config set registry https://registry.npmmirror.com && npm install --no-audit --no-fund
COPY . .
RUN npm run build

# ---- runtime ----
FROM nginx:1.27-alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
