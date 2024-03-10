FROM debian:12.5-slim

RUN apt-get update
RUN apt-get install ca-certificates -y
RUN apt-get upgrade -y

RUN useradd -m app

USER app
WORKDIR /home/app
COPY backend .

EXPOSE 8080

CMD ["./backend"]
