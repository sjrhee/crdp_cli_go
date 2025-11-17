#!/bin/bash

# 임시 디렉토리
TMPDIR=$(mktemp -d)
trap "rm -rf $TMPDIR" EXIT

# ECDSA P-256 키쌍 생성
openssl ecparam -name prime256v1 -genkey -noout -out "$TMPDIR/priv.pem"
openssl ec -in "$TMPDIR/priv.pem" -pubout -out "$TMPDIR/pub.pem" 2>/dev/null

# JWT 헤더, 페이로드
HEADER=$(echo -n '{"alg":"ES256","typ":"JWT"}' | base64 -w0 | tr '+/' '-_' | tr -d '=')
PAYLOAD=$(echo -n "{\"exp\":$(date +%s -d '+1 day'),\"iss\":\"CRDP03\",\"sub\":\"user01\"}" | base64 -w0 | tr '+/' '-_' | tr -d '=')

# 서명 메시지
MESSAGE="${HEADER}.${PAYLOAD}"

# DER 서명 생성 및 Base64URL 변환
SIGNATURE=$(echo -n "$MESSAGE" | openssl dgst -sha256 -sign "$TMPDIR/priv.pem" | base64 -w0 | tr '+/' '-_' | tr -d '=')

# 최종 JWT
JWT="${MESSAGE}.${SIGNATURE}"

# 출력
echo "=== Public Key ==="
cat "$TMPDIR/pub.pem"
echo ""
echo "=== JWT Token ==="
echo "$JWT"
