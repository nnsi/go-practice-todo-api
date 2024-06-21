# Todo API

## バックエンド(Go)

```bash
docker-compose build
docker-compose up -d
```

API サーバーのログを確認する場合は`docker-compose logs -f backend`で確認可能

### API 仕様

- `POST /login` でアクセストークンを取得
  - login_id string, password string
- `POST /register` でユーザー登録
  - username string, login_id string, password string
- `GET /todos` で全取得
- `POST /todos` で追加
  - title string
- `GET /todos/{id}` で特定の ID のみ取得
- `PUT /todos/{id}` で更新
  - title? string, completed? boolean
- `DELETE /todos/{id}` で削除

## フロントエンド(React 19 + Vite)

```bash
cd frontend
npm install
npm run dev
```

## デプロイ先設定

- フロントエンドは GitHub Pages、バックエンドは Fly.io へデプロイを行うよう設定している
  - GitHub の Setting で Actions -> General から Workflow permissions を Read and write に設定する
  - Fly.io はデプロイ後に`flyctl ips allocate-v4 --shared -a go-practice-todo-app`とコマンドを打ってサービスに IP を付与する
- DB の接続情報は Fly.io の secret に記載

```sh
flyctl secrets set DB_HOST=<your-db-host> DB_USER=<your-db-user> DB_PASSWORD=<your-db-password> DB_NAME=<your-db-name> DB_PORT=<your-db-port>
```

## 動作サンプル

- https://nnsi.github.io/go-practice-todo-api/
- API:https://go-practice-todo-app.fly.dev/
