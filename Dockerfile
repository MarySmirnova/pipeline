FROM golang AS compiling
RUN mkdir -p /go/pipeline
WORKDIR /go/pipeline
COPY . .
RUN go build -o bin/app .

FROM ubuntu AS building
WORKDIR /opt
COPY --from=compiling /go/pipeline/bin/app .
CMD ["./app"]