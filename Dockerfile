FROM ubuntu:latest

WORKDIR /app

COPY app /app/app
COPY web /app/web

EXPOSE 7540

CMD ["/app/app"]
