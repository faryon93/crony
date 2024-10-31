# ----------------------------------------------------------------------------------------
# Image: Builder
# ----------------------------------------------------------------------------------------
FROM golang:1.23-alpine AS builder

# setup the environment
ENV TZ=Europe/Berlin

# install dependencies
RUN apk --update --no-cache add git gcc musl-dev
WORKDIR /work
ADD ./ ./

# build the go binary
RUN go build -ldflags \
        '-X "main.BuildTime='$(date -Iminutes)'" \
         -X "main.GitCommit='$(git rev-parse --short HEAD)'" \
         -X "main.GitBranch='$(git rev-parse --abbrev-ref HEAD)'" \
         -X "main.BuildNumber='$CI_BUILDNR'" \
         -s -w' \
         -v -o /tmp/crony && \
        chown nobody:nobody /tmp/crony

# ----------------------------------------------------------------------------------------
# Image: Deployment
# ----------------------------------------------------------------------------------------
FROM alpine:latest

# setup the environment
ENV TZ=Europe/Berlin

RUN apk --update --no-cache add ca-certificates tzdata bash && mkdir /crony

# add relevant files to container
COPY --from=builder /tmp/crony /usr/bin/crony

WORKDIR /crony
ENTRYPOINT ["/usr/bin/crony"]
