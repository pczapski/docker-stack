#!/usr/bin/env sh

vault auth enable oidc

vault write auth/oidc/role/demo \
    user_claim="email" \
    allowed_redirect_uris="http://127.0.0.1:8200/ui/vault/auth/oidc/oidc/callback"  \
    groups_claim="groups" \
    policies="default"

vault write auth/oidc/config \
    oidc_client_id="vault" \
    oidc_client_secret="demo" \
    default_role="demo" \
    oidc_discovery_url="http://keycloak:8090/auth/realms/demo"