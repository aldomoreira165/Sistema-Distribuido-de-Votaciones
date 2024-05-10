FROM rust:latest as builder
WORKDIR /usr/src/app
COPY . .
RUN cargo build --release
FROM ubuntu:latest
RUN apt-get update && apt-get install -y libssl-dev
COPY --from=builder /usr/src/app/target/release/app-cliente /usr/local/bin/app-cliente
EXPOSE 8000
CMD ["app-cliente"]