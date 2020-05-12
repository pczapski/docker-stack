#!/usr/bin/env sh
vault secrets enable pki
vault secrets tune -max-lease-ttl=8760h pki