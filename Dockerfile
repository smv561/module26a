FROM golang:latest AS golang_stage
RUN mkdir -p /go/src/pipeline
WORKDIR /go/src/pipeline
ADD pipe.go .
ADD go.mod .
#RUN go install .
RUN go build -o pipe .

FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /root/
#COPY --from=golang_stage /go/src/pipeline .
COPY --from=golang_stage /go/src/pipeline/pipe .
#ENTRYPOINT ./pipe
CMD ["./pipe"]