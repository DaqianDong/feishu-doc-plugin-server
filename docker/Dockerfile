FROM alpine:3.18

# 替换为阿里云镜像源并安装 tzdata
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk add --no-cache tzdata

ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# 确保你的 bin 目录在构建上下文内（建议将 bin 移至 Dockerfile 同级目录）
COPY bin/server /data/feishu/server

WORKDIR /data/feishu

CMD ["./server"]