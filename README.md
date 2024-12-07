## 概要
このリポジトリは、Go言語のWebフレームワークであるGin, Echo, Fiber, Chiを比較するためのプロジェクトです。

## 環境
- Go 1.23.2
- MacOSのみでの動作確認

## プロジェクト構成
``repository``にはInMemoryでデータを参照・永続化するための処理があり、``shared``には共通で利用する処理があります。  
それ以外のディレクトリはそれぞれのフレームワークごとの処理とテストコードがあります。
```
.
├── main.go
├── chiapp
│   ├── routes.go
│   └── routes_test.go
├── echoapp
│   ├── routes.go
│   └── routes_test.go
├── fiberapp
│   ├── routes.go
│   └── routes_test.go
├── ginapp
│   ├── routes.go
│   └── routes_test.go
├── repository
│   └── todos.go
└── shared
    ├── server.go
    └── token.go
```

## 実行方法
プロジェクトをクローンして以下コマンドを実行してください。
#### サーバーの起動
``-server``と``-addr``のオプションを指定してサーバーを起動します。  
``-server``には``gin``, ``echo``, ``fiber``, ``chi``のいずれかを指定します。
``-addr``には``host:port``の形式でアドレスを指定します。
```bash
go run main.go -server gin -addr localhost:8080
```
上記を実行すると、``localhost:8080``でサーバーが起動しますので、ブラウザで``http://localhost:8080/api/v1/todos?token=token``にアクセスするとjson形式でデータが返ってきます。  
``token``は``shared/token.go``に定義されているものを利用してください。

#### テストの実行
以下のコマンドでテストを実行します。
```bash
go test app/...
```


## ルール

### コミットメッセージ
- コミットメッセージは必ず日本語で書く
- コミットメッセージは以下のフォーマットで書く
  - `[add/feat/fix/remove/refactor/chore/...] メッセージ`
  - 例: `[add] 新しいファイルを追加`
