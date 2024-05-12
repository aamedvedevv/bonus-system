up: start_postgres accrual gophermart


start_postgres:
	docker-compose -f docker-compose.yaml up -d

accrual:
	./cmd/accrual/accrual_darwin_arm64 &

gophermart:
	go run ./cmd/gophermart/main.go -d "host=127.0.0.1 port=5432 user=postgres sslmode=disable password=1234" -a "127.0.0.1:8081" &