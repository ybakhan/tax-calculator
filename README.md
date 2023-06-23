# Tax Calculator Assignment

REST API for calculating taxes


## Build 

To build the project, run the following command:

```bash
make build
```

Binary is /bin/taxcalculator

## Unit Test 

To run unit tests, run the following command:

```bash
make test
```

## Code Coverage 

To see code coverage of unit test, run the following command:

```bash
make coverage
```

Detailed report of coverage is coverage.html

## Integration Test 

To run integration tests, run the following command:

```bash
make testIT
```

## Run Tax Calculator

To run tax calculator, run the following command:

```bash
make run
```

To get taxes for a year and salary call the following api 

```plaintext 
http://localhost:8081/tax/{year}?s={salary}
```

e.g To get taxes for year 2022 and salary 80000 call

```plaintext 
http://localhost:8081/tax/2022?s=80000
```