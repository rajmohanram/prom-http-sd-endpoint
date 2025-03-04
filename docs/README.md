# Prometheus HTTP Service Discovery

This is a simple service discovery mechanism for Prometheus that uses HTTP to discover targets. 

## Usage

To use this service discovery mechanism, you need to run the `prometheus-http-sd` containerized application. This application will expose a `/<jobs.name>` endpoint that will return a list of targets in the Prometheus static target format.
