FROM golang:alpine as builder

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    GIN_MODE=release \
    PORT=8799

# install makefile tools
RUN apk --no-cache add make gcc libc-dev git

# create and use work directory
WORKDIR /app

#copy go models file
COPY go.mod go.sum ./
RUN go mod download
# copy project source file to work dir
COPY . .

# exec Makefile
RUN make
# compile go source file to app dir
RUN go build -o ./executable

# use pure alpine as runtime environment
FROM alpine AS runner
RUN apk --no-cache add ca-certificates

COPY --from=builder /app/executable .
COPY --from=builder /app/config ./config
COPY --from=builder /app/docs   ./docs

ARG env
ENV env=${env:-prod}

# replace source with aliyun
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && apk update
# add time zone setting
RUN apk update && apk add tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone && apk del tzdata

EXPOSE 8799
# run go executable file
# use shell form to support env param
ENTRYPOINT ./executable "env=$env"