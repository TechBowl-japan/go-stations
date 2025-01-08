# 環境構築

1. Node.js（v18以降を推奨します。）
2. Yarn (v1)
3. Git
4. Visual Studio Code（VSCode）
5. Railway VSCode 拡張機能

上記をインストールする必要があります。インストールできているかの確認やインストール方法は、
[Railway 準備編](https://www.notion.so/techbowl/Railway-ceba695d5014460e9733c2a46318cdec) をご確認いただき、挑戦の準備をしましょう。※ GitHub Codespaces についての資料はスキップしてください。

## トラブルシューティング

### go test で 404 というエラーが返ってきます。

main.goなどで handler の登録を確認してみましょう。
テストの関係上 router.NewRouter のメソッド内部で追加するようにしましょう。

### DBに接続して中身が見れないのですが？

次のような結果が返ってきていれば、正常です。

```
$ sqlite3 .sqlite3/todo.db
SQLite version 3.32.3 2020-06-18 14:16:19
Enter ".help" for usage hints.
sqlite> .tables
todos
```

もし、 `todos` が作成されていないようであれば、次のコマンドを実行しましょう。

```
$ sqlite3 .sqlite3/todo.db < db/schema.sql
```

これで、 `todos` が作成されていれば、問題なく接続できます。
