// Allow user to list all namespaces secrets (include nested namespaces)
// IMPORTANT! secret/metadata
path "secret/metadata/*" {
  capabilities = ["list"]
}
// Disallow user to list secrets in login namespace
// IMPORTANT! secret/metadata
path "secret/metadata/devops" {
  capabilities = ["deny"]
}
// Allow user to read and update secrets in all namespaces under example
// + in trading-bridge/+/ is a directory wildcard, it matches all envs
// example/dev/
// example/cicd/ etc
// * in /example/+/* is glob selector for all secrets
// IMPORTANT! secret/data
path "secret/data/user/application/example/+/*" {
  capabilities = ["read", "update"]
}

// Allow user only to read secrets in prod namespaces under example
// IMPORTANT! secret/data
path "secret/data/user/application/example/prod/*" {
  capabilities = ["read"]
}

// Allow user only to read secrets in public namespaces
path "secret/data/public/*" {
  capabilities = ["read"]
}

path "secret/metadata/public/*" {
  capabilities = ["list"]
}
