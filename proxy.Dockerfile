FROM nginx

WORKDIR /app
COPY ./devProxy/mime.types .
COPY ./devProxy/docker.nginx.conf nginx.conf 
COPY proxy.nginx.conf /etc/nginx/nginx.conf

EXPOSE 80
