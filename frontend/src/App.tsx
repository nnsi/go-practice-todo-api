import { useState } from "react";
import { TodoApp } from "../components/TodoApp";

const LoginForm: React.FC<{ setToken: (token: string) => void }> = ({
  setToken,
}) => {
  return (
    <form
      action={async (formData) => {
        try {
          const req = await fetch("http://localhost:8080/login", {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify({
              login_id: formData.get("id"),
              password: formData.get("password"),
            }),
          });
          const data = await req.json();
          setToken(data.token);
          localStorage.setItem("token", data.token);
        } catch (e) {
          console.error(e);
        }
        return false;
      }}
    >
      <input type="text" name="id" />
      <input type="text" name="password" />
      <button type="submit">認証</button>
    </form>
  );
};

const SignUpForm: React.FC<{ setToken: (token: string) => void }> = ({
  setToken,
}) => {
  return (
    <form
      action={async (formData: FormData) => {
        try {
          const req = await fetch("http://localhost:8080/register", {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify({
              login_id: formData.get("login_id"),
              username: formData.get("username"),
              password: formData.get("password"),
            }),
          });
          const res = await req.json();
          console.log(res.token);
          setToken(res.token);
          localStorage.setItem("token", res.token);
        } catch (e) {
          console.error(e);
        }
        return false;
      }}
    >
      <input type="text" name="username" placeholder="username" />
      <input type="text" name="login_id" placeholder="id" />
      <input type="text" name="password" placeholder="password" />
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
      {!token && (
        <>
          <LoginForm setToken={setToken} />
          <hr />
          <SignUpForm setToken={setToken} />
        </>
      )}
    </>
  );
};

export default App;
