FROM golang:1.16 as builder

COPY . /app
WORKDIR /app
RUN make dc

FROM busybox:1.33.0

COPY --from=builder /app/bin/dc /dc

EXPOSE 1337

ENTRYPOINT ["./dc"]

CMD "-h"
