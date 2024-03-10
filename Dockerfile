FROM golang:latest

RUN apt-get update
RUN apt-get install ca-certificates make -y
RUN apt-get upgrade -y

RUN useradd -m app

USER app
WORKDIR /home/app

COPY . . 
RUN make -j

EXPOSE 8080

CMD ["./backend"]
