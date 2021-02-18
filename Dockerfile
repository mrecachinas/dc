FROM golang:1.16 as builder

COPY . /app
WORKDIR /app
RUN make dc

FROM busybox:1.33.0-glibc

COPY --from=builder /app/bin/dc /dc

ENTRYPOINT ["./dc"]

CMD "-h"