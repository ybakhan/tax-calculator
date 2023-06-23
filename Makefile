build:
	@go build -o bin/taxcalculator  ./main 

run:
	@docker-compose up --build tax-calculator

test:
	@go test ./... -v -count=1 -tags=unit

coverage:
	@go test -coverprofile=coverage.out ./... -v -count=1 -tags=unit
	@go tool cover -func=coverage.out
	@go tool cover -html=coverage.out -o coverage.html

testIT:
	@docker-compose up --build integration-test

testIT-local:
	@INTERVIEW_SERVER=http://localhost:5000 TAX_CALCULATOR_SERVER=http://localhost:8081 go test ./integration -tags=integration -v -count=1

