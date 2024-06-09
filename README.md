# Todo API

## API Server

- air を利用したホットリロード
  - windows の場合は`air -c .air.windows.toml`
  - それ以外は`air`
- `GET /todos` で全取得
- `POST /todos` で追加
- `GET /todos/{id}` で特定の ID のみ取得
- `PUT /todos/{id}` で更新
- `DELETE /todos/{id}` で削除
