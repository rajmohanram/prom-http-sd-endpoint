# prom-http-sd-endpoint
 Prometheus HTTP Service Discovery Endpoint

## Usage

Refer to [documentation](./docs/) for more information.

## Configuration
```yaml
jobs:
  - name: node-exporter
    targets:
      - 1.1.1.1:9001
      - 2.2.2.2:9001
    labels:
      env: prod
      dc: dc1
  - name: podman-exporter
    targets:
      - 3.3.3.3:9002
      - 4.4.4.4:9002
    labels:
      env: prod
      dc: dc2

```

## Output (/node-exporter)
```json
[
  {
    "labels": {
      "dc": "dc1",
      "env": "prod"
    },
    "targets": [
      "1.1.1.1:9001",
      "2.2.2.2:9001"
    ]
  }
]
```
