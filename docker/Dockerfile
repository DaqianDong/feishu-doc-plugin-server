FROM ubuntu:22.04

RUN apt-get update
RUN apt-get install tzdata
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

COPY ../bin/server /data/feishu/server

# backend
ENV PORT=8081

WORKDIR /data/feishu

CMD ["server"]
  +++++++ REPLACE
