#!/bin/bash

echo "=== Basic認証テストを開始します ==="

# 1. ルートエンドポイントへのアクセス（認証なし）
echo -e "\n1. ルートエンドポイントテスト（認証なし）"
curl -i http://localhost:8080/

# 2. /healthzエンドポイントへの正常アクセス
echo -e "\n2. /healthz 正常系テスト - 正しい認証情報"
curl -i -u test:test http://localhost:8080/healthz

# 3. /healthzエンドポイントへの異常アクセス（誤った認証情報）
echo -e "\n3. /healthz 異常系テスト - 誤った認証情報"
curl -i -u wrong:wrong http://localhost:8080/healthz

# 4. /healthzエンドポイントへの異常アクセス（空の認証情報）
echo -e "\n4. /healthz 異常系テスト - 空の認証情報"
curl -i -u : http://localhost:8080/healthz

# 5. /healthzエンドポイントへの異常アクセス（認証情報なし）
echo -e "\n5. /healthz 異常系テスト - 認証情報なし"
curl -i http://localhost:8080/healthz

echo -e "\n=== テスト完了 ==="