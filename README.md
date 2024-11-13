<p align="center">
  <img src="assets/icon.png" width="100" height="100" />
</p>
<h1 align="center">
  Burden - Load Testing CLI Tool
</h1>

Burden is a CLI-based load testing tool designed to simulate requests to a specified URL or API collection. The tool measures key performance metrics, helping you evaluate the efficiency and resilience of your application under load.

Features
--------

*   Import collections to test multiple routes in your application.
*   Configure the number of concurrent users and specify either the total number of requests or test duration.
*   Measure critical metrics:
    *   Throughput
    *   Response Time
    *   Latency
    *   Errors
    *   Resource Utilization
    *   Concurrency
    *   Peak Load
    *   Downtime
*   Basic and detailed output modes to suit your testing needs.

Installation
------------

Download the binary or clone the repository, navigate to the project directory, and build the executable.

    git clone https://github.com/vladkanatov/burden.git
    cd burden
    go build -o burden cmd/main.go

Usage
-----

### Basic Usage

    ./burden --url http://example.com/api --users 10 --requests 1000

This command runs a load test on `http://example.com/api` with 10 concurrent users sending a total of 1000 requests.

### Options

*   `--url` : The target URL for the load test (required if `--collection` is not used).
*   `--collection` : Path to the request collection file for testing multiple routes (optional).
*   `--users` : Number of concurrent users to simulate (default: 1).
*   `--requests` : Total number of requests to send (default: 100).
*   `--max-errors` : Maximum allowable errors before stopping the test (-1 to disable).
*   `--detailed` : Display detailed metrics (optional).

### Example

    ./burden --collection ./path/to/collection.json --users 20 --requests 2000 --max-errors 10 --detailed

Output Metrics
--------------

*   **Throughput** : Requests per second.
*   **Response Time** : Average time taken for a request.
*   **Latency** : Average time to receive a response after sending a request.
*   **Errors** : Total number of errors encountered during the test.
*   **Resource Utilization** : Resource usage percentage (CPU, Memory).
*   **Concurrency** : Number of requests being processed concurrently.
*   **Peak Load** : Highest observed load during the test.
*   **Downtime** : Total time the service was unavailable during the test.

By default, Burden outputs three key metrics: Throughput, Response Time, and Latency. Use `--detailed` for full metric reporting.

License
-------

Burden is open-source software licensed under the MIT License.