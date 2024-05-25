FROM golang:1.22 as builder
WORKDIR /go/src/app
COPY . ./
RUN go mod download && go mod verify
RUN CGO_ENABLED=0 go build -o /go/bin/netdbg

FROM gcr.io/distroless/static-debian11
COPY --from=builder /go/bin/netdbg /
CMD ["/netdbg"]
