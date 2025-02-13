# helm-aggregator

A service that aggregates Helm repositories into a single entry point

## Description

This utility can aggregate multiple Helm repositories into one.

## Usage

Config file config.yaml

```yaml
repos:
  - name: wiremind
    url: https://wiremind.github.io/wiremind-helm-charts

port: "8080"
```
