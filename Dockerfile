FROM node:20-alpine as nodejs
COPY ./web /build/web
WORKDIR /build/web
RUN npm config set registry https://registry.npmmirror.com
RUN npm install
RUN npm run build

FROM golang:1.21-alpine as go
ENV GOPROXY=https://goproxy.cn,direct
COPY ./md /build/md
COPY --from=nodejs /build/web/dist /build/md/web
WORKDIR /build/md
RUN go build


FROM alpine:latest
COPY --from=go /build/md/md /md/
COPY /md/data/ /var/data
ENV reg=false
EXPOSE 4001
RUN chmod +x /md/md
CMD cp -R /var/data /md && /md/md -p 4001 -log /md/logs -data /md/data -reg=${reg} -pg_host=${pg_host} -pg_port=${pg_port} -pg_user=${pg_user} -pg_password=${pg_password} -pg_db=${pg_db}

