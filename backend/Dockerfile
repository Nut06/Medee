FROM dhi.io/golang:1.26.1-debian13-dev as build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping ./cmd/api

FROM build as run-test-stage
RUN go test -v ./...

FROM dhi.io/golang:1.26.1-debian13 as runner

WORKDIR /

COPY --from=build /docker-gs-ping /docker-gs-ping

EXPOSE 5173

USER nonroot:nonroot

ENTRYPOINT [ "/docker-gs-ping" ]