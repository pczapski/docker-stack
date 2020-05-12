#!/usr/bin/env sh

vault auth enable oidc

vault write auth/oidc/config \
    oidc_client_id="vault" \
    oidc_client_secret="demo" \
    default_role="user" \
    oidc_discovery_url="http://keycloak:8090/auth/realms/demo"

vault write auth/oidc/role/user \
    user_claim="email" \
    allowed_redirect_uris="http://127.0.0.1:8200/ui/vault/auth/oidc/oidc/callback" \
    groups_claim="groups" \
    policies="default, user"

vault auth list -format=json | grep -Eo 'auth_oidc.+[a-zA-Z0-9]' > accessor

vault write -field=id identity/group name="devops" type="external" \
        policies="devops" \
        metadata=responsibility="Manage K/V Secrets" > groupId

vault write identity/group-alias name="devops" mount_accessor=$(cat accessor) canonical_id=$(cat groupId)

vault write -field=id identity/group name="users" type="external" \
        policies="users" \
        metadata=responsibility="Manage K/V Secrets" > groupId

vault write identity/group-alias name="users" mount_accessor=$(cat accessor) canonical_id=$(cat groupId)
