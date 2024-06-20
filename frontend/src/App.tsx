import { useState } from "react";
import { TodoApp } from "./components/TodoApp";
import { LoginForm } from "./components/LoginForm";
import { SignupForm } from "./components/SignupForm";

const App: React.FC = () => {
  const [token, setToken] = useState(localStorage.getItem("token"));

  const handleLogout = () => {
    localStorage.removeItem("token");
    setToken(null);
  };

  return (
    <>
      {token ? (
        <>
          <TodoApp token={token} />
          <hr />
          <button onClick={handleLogout}>ログアウト</button>
        </>
      ) : (
        <>
          <LoginForm setToken={setToken} />
          <hr />
          <SignupForm setToken={setToken} />
        </>
      )}
    </>
  );
};

export default App;
