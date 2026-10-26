FROM golang:1.27.2-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY *.go ./
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /gowatch .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /gowatch /gowatch
ENV GOWATCH_ADDR=0.0.0.0:8080
EXPOSE 8080
ENTRYPOINT ["/gowatch"]
