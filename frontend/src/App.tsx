import { useState } from "react";
import { TodoApp } from "../components/TodoApp";

const LoginForm: React.FC<{ setToken: (token: string) => void }> = ({
  setToken,
}) => {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  return (
    <form
      onSubmit={async (e) => {
        e.preventDefault();
        try {
          const req = await fetch("http://localhost:8080/login", {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify({ login_id: username, password }),
          });
          const data = await req.json();
          console.log(data.token);
          setToken(data.token);
          localStorage.setItem("token", data.token);
        } catch (e) {
          console.error(e);
        }
      }}
    >
      <input
        type="text"
        value={username}
        onChange={(e) => setUsername(e.target.value)}
      />
      <input
        type="password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
      />
      <button type="submit">認証</button>
    </form>
  );
};

const App: React.FC = () => {
  const [token, setToken] = useState(localStorage.getItem("token"));

  const handleLogout = () => {
    localStorage.removeItem("token");
    setToken(null);
  };

  return (
    <>
      {token && (
        <>
          <TodoApp token={token} />
          <hr />
          <button onClick={handleLogout}>ログアウト</button>
        </>
      )}
      {!token && <LoginForm setToken={setToken} />}
    </>
  );
};

export default App;
