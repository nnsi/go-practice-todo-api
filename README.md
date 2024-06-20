# Todo API

## API Server

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

## Frontend

- React 19 + Vite
- `localhost:5173`で起動

```bash
cd frontend
npm install
npm run dev
```
