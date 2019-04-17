ARG VERSION=1.12
FROM golang:$VERSION as builder

WORKDIR /src
ENV GO111MODULE=on

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . .

WORKDIR /src
RUN CGO_ENABLED=1 go build -ldflags '-extldflags "-static"' -o /release/app
RUN cd /release && ls -aGlhSr

FROM gcr.io/distroless/base as runner
COPY --from=builder /release /

CMD ["/app", "-m"]
