#!/usr/bin/env sh

vault kv put secret/devops/infrastructure/cluster/values.yaml value=example

vault kv put secret/user/application/example/staging/values.yaml value=example

vault kv put secret/public/values.yaml value=example
