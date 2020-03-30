# SnowedIN

## Use case

Let alerts sent from AlertManager be templated and forwared to a configurable endpoint

## Usage

Configuration is done through `config.yaml`, command line arguments and/or environment variables, notable `SERVICENOW_USERNAME` and `SERVICENOW_PASSWORD`.
Command line arguments always take precedent over environment variables.

The fields to be posted to the ServiceNow webhook can be specified under `default_incident` in the `config.yaml`.

