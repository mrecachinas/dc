FROM golang:1.15 as builder

COPY . /app
RUN make dc

FROM busybox

COPY --from=builder /app/dc /dc

ENTRYPOINT ["./dc"]

CMD "-h"