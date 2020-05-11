curl \
    --request PUT \
    --data @/init-data/api-gateway.json \
    http://127.0.0.1:8500/v1/kv/app/api-gateway
curl \
    --request PUT \
    --data @/init-data/worker.json \
    http://127.0.0.1:8500/v1/kv/app/worker