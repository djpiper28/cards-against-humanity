FROM golang:latest as build

RUN apt-get update
RUN apt-get install ca-certificates make npm -y
RUN apt-get upgrade -y
RUN npm i -g pnpm

WORKDIR /build

# Set context to ../
COPY . .
RUN make backend -j

FROM build as release

RUN useradd -m app
WORKDIR /home/app
USER app

ENV PORT 8080
ENV GIN_MODE release

EXPOSE 8080
COPY --from=build /build/backend/backend .
COPY --from=build /build/backend/packs/cah-all-compact.json packs/cah-all-compact.json
CMD ["./backend"]
