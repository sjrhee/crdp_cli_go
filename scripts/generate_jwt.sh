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

# DER 서명 생성
TMPDIR=$(mktemp -d)
trap "rm -rf $TMPDIR" EXIT

DER_SIG=$(echo -n "$MESSAGE" | openssl dgst -sha256 -sign <(printf "%s" "$PRIKEY") 2>/dev/null)
echo -n "$DER_SIG" > "$TMPDIR/sig.der"

# DER을 r||s 형식으로 변환 (각 32바이트씩, 총 64바이트)
# OpenSSL을 사용하여 DER을 r||s로 변환
RS_SIG=$(python3 -c "
import sys, base64, subprocess

# DER 서명을 r||s로 변환
der_sig = open(sys.argv[1], 'rb').read()

# DER 파싱: SEQUENCE { INTEGER r, INTEGER s }
def parse_der_int(data, offset):
    if data[offset] != 0x02:  # INTEGER tag
        raise ValueError('Not an INTEGER')
    offset += 1
    length = data[offset]
    offset += 1
    # 상위 바이트가 0x00이면 제거
    value_bytes = data[offset:offset+length]
    if len(value_bytes) > 1 and value_bytes[0] == 0x00:
        value_bytes = value_bytes[1:]
    return value_bytes, offset + length

# 서명 파싱
offset = 2  # SEQUENCE 헤더 건너뛰기
r_bytes, offset = parse_der_int(der_sig, offset)
s_bytes, offset = parse_der_int(der_sig, offset)

# r, s를 32바이트씩으로 패딩
r_bytes = r_bytes.rjust(32, b'\x00')[-32:]
s_bytes = s_bytes.rjust(32, b'\x00')[-32:]

# r||s 형식으로 연결 (총 64바이트)
rs = r_bytes + s_bytes

# Base64URL 인코딩
encoded = base64.urlsafe_b64encode(rs).decode().rstrip('=')
print(encoded)
" "$TMPDIR/sig.der")

# 최종 JWT
JWT="${MESSAGE}.${RS_SIG}"

# 출력
echo "=== Public Key ==="
echo "$PUBKEY"
echo ""
echo "=== JWT Token ==="
echo "$JWT"
