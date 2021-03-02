FROM node:14.15.1 as jsbuilder

COPY ui/webapp /app
WORKDIR /app
RUN npm i && npm run build

FROM golang:1.16.0-buster as gobuilder

# TODO: Is there a way we don't need to
#       copy all of . and just the go stuff?
COPY . /app
COPY --from=jsbuilder /app /app/ui/webapp
WORKDIR /app
RUN make dc

FROM busybox:1.33.0

RUN adduser -D -H -u 10001 dcuser
USER dcuser

COPY --from=gobuilder /app/bin/dc /dc

EXPOSE 1337

ENTRYPOINT ["./dc"]

CMD "-h"
