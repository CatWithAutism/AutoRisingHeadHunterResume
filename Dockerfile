FROM golang:latest
LABEL authors="catwithaut1sm"
WORKDIR /app
COPY src/go.mod src/go.sum ./
RUN go mod download
COPY src/*.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /hh
CMD ["/hh"]