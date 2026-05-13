FROM golang:1.22

WORKDIR /app
COPY . .

WORKDIR /app/golden_repo
RUN go mod tidy
RUN go build -o /app/ratelimiter ./cmd/ratelimiter

WORKDIR /app/behavioral_tests
RUN pip install pytest

CMD ["pytest", "-q"]
