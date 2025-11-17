#!/bin/bash

# JWT 설정 변수
ISSUER="CRDP03"
SUBJECT="user01"
EXPIRY_DAYS=3650  # 10년

PRIKEY="
-----BEGIN EC PARAMETERS-----
BggqhkjOPQMBBw==
-----END EC PARAMETERS-----
-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIE7CQOZPvG7j5MJO82o+WgmqmZkHNqWllLazvkwa5+KJoAoGCCqGSM49
AwEHoUQDQgAELZo4vTZ1ypjIZB/KzOAeRbS52Z0GsP5nYXcddc2xP16Rm3+bLdib
0OxWsCm1ltEtg9rM+dXQRXCKlGkMgqnsVw==
-----END EC PRIVATE KEY-----
"

# ECDSA 공개키 생성
PUBKEY="$(openssl ec -in <(printf "%s" "$PRIKEY") -pubout 2>/dev/null)"

# JWT 헤더, 페이로드
HEADER=$(echo -n '{"alg":"ES256","typ":"JWT"}' | base64 -w0 | tr '+/' '-_' | tr -d '=')
PAYLOAD=$(echo -n "{\"exp\":$(date +%s -d "+${EXPIRY_DAYS} days"),\"iss\":\"${ISSUER}\",\"sub\":\"${SUBJECT}\"}" | base64 -w0 | tr '+/' '-_' | tr -d '=')

# 서명 메시지
MESSAGE="${HEADER}.${PAYLOAD}"

# DER 서명 생성 및 Base64URL 변환
SIGNATURE=$(echo -n "$MESSAGE" | openssl dgst -sha256 -sign <(printf "%s" "$PRIKEY") | base64 -w0 | tr '+/' '-_' | tr -d '=')

# 최종 JWT
JWT="${MESSAGE}.${SIGNATURE}"

# 출력
echo "=== Public Key ==="
echo "$PUBKEY"
echo ""
echo "=== JWT Token ==="
echo "$JWT"
