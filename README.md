# GO-api
2025.04 ~ APIエンドポイントの作成

## フォルダ構成
<pre>
GO-api/
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
├── cmd/
│   └── api-server/
│       └── main.go
└── internal/
    ├── db/
    │   ├── access_log.go
    │   └── init.go
    ├── geoapi/
    │   └── client.go
    ├── handler/
    │   ├── access_logs.go
    │   ├── address.go
    │   └── hello.go
    └── model/
        ├── location_logic.go
        └── types.go
</pre>

## 実装ポイント
- apiエンドポイントの作成  
　webサーバーとしてクライアント（ブラウザ）からリクエストを受け取り、レスポンスを返す基礎実装
- 外部apiを叩く  
　外部apiを叩き、リクエストを送り、レスポンスを受け取り適切な処理を実装
- dbとの連携  

## 課題・今後の実装試み
- 認証について  
　認証の実装や、cookieやセッションについても実装したい
- アーキテクチャについて  
　ソフトウェアを構築する上での適切なアーキテクチャや、GOでのベストプラクティスを学び真似したい
- テストコード  
 ユニットテスト、結合テストなどのコードを書けるようになりたい
 
