#!/bin/bash

set -e
set -x

swagger-cli bundle api/openapi/index.yaml -o api/openapi/summary.yaml --type yaml

oapi-codegen -generate types -o ./internal/infrastructure/http/api_dto/api_dto.go -package api_dto ./api/openapi/summary.yaml