#!/bin/bash

# JWT 설정 변수
ISSUER="CRDP03"
SUBJECT="user01"
EXPIRY_DAYS=3650  # 10년

PRIKEY="-----BEGIN EC PARAMETERS-----
BggqhkjOPQMBBw==
-----END EC PARAMETERS-----
-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIE7CQOZPvG7j5MJO82o+WgmqmZkHNqWllLazvkwa5+KJoAoGCCqGSM49
AwEHoUQDQgAELZo4vTZ1ypjIZB/KzOAeRbS52Z0GsP5nYXcddc2xP16Rm3+bLdib
0OxWsCm1ltEtg9rM+dXQRXCKlGkMgqnsVw==
-----END EC PRIVATE KEY-----"

# ECDSA 공개키 생성
PUBKEY=$(echo "$PRIKEY" | openssl ec -pubout 2>/dev/null)

# JWT 헤더, 페이로드
HEADER=$(echo -n '{"alg":"ES256","typ":"JWT"}' | base64 -w0 | tr '+/' '-_' | tr -d '=')
PAYLOAD=$(echo -n "{\"exp\":$(date +%s -d "+${EXPIRY_DAYS} days"),\"iss\":\"${ISSUER}\",\"sub\":\"${SUBJECT}\"}" | base64 -w0 | tr '+/' '-_' | tr -d '=')

# 서명 메시지
MESSAGE="${HEADER}.${PAYLOAD}"

# 임시 파일 생성
KEYFILE=$(mktemp)
SIGFILE=$(mktemp)
trap "rm -f $KEYFILE $SIGFILE" EXIT

# 개인키 저장
echo "$PRIKEY" > "$KEYFILE"

# DER 서명 생성
echo -n "$MESSAGE" | openssl dgst -sha256 -sign "$KEYFILE" > "$SIGFILE" 2>/dev/null

# 환경 변수에 시그니처 파일 경로 설정
export SIG_FILE="$SIGFILE"

# Python을 사용하여 DER을 r||s로 변환
RS_SIG=$(python3 << 'PYTHON_SCRIPT'
import base64
import os

# 환경 변수에서 시그니처 파일 경로 받기
sig_file = os.environ.get('SIG_FILE')

# DER 서명 파일 읽기
with open(sig_file, 'rb') as f:
    der_sig = f.read()

# DER 구조 파싱: SEQUENCE { INTEGER r, INTEGER s }
offset = 0

# SEQUENCE 태그 (0x30)
if der_sig[offset:offset+1] != b'\x30':
    raise ValueError("Invalid DER signature")
offset += 1

# SEQUENCE 길이
seq_len = der_sig[offset]
offset += 1

# 첫 번째 INTEGER (r) 파싱
if der_sig[offset:offset+1] != b'\x02':
    raise ValueError("Invalid r tag")
offset += 1
r_len = der_sig[offset]
offset += 1
r_bytes = der_sig[offset:offset+r_len]
offset += r_len

# 두 번째 INTEGER (s) 파싱
if der_sig[offset:offset+1] != b'\x02':
    raise ValueError("Invalid s tag")
offset += 1
s_len = der_sig[offset]
offset += 1
s_bytes = der_sig[offset:offset+s_len]

# 32바이트로 패딩
r_padded = (b'\x00' * 32 + r_bytes)[-32:]
s_padded = (b'\x00' * 32 + s_bytes)[-32:]

# r||s 결합
rs_sig = r_padded + s_padded

# Base64URL 인코딩
encoded = base64.urlsafe_b64encode(rs_sig).decode().rstrip('=')
print(encoded)
PYTHON_SCRIPT
)

# 환경 변수 해제
unset SIG_FILE

# 최종 JWT
JWT="${MESSAGE}.${RS_SIG}"

# 출력
echo "=== Public Key ==="
echo "$PUBKEY"
echo ""
echo "=== JWT Token ==="
echo "$JWT"
