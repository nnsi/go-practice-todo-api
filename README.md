# Todo API

## API Server

```bash
$ docker-compose build
$ docker-compose up -d
$ docker-compose exec backend air
```

### API 仕様

- `POST /login` でアクセストークンを取得
- `POST /register` でユーザー登録
- `GET /todos` で全取得
- `POST /todos` で追加
- `GET /todos/{id}` で特定の ID のみ取得
- `PUT /todos/{id}` で更新
- `DELETE /todos/{id}` で削除

## Frontend

- React 19 + Vite
- `localhost:5173`で起動

```bash
$ cd frontend
$ npm install
$ npm run dev
```
