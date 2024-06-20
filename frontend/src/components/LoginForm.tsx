export const LoginForm: React.FC<{ setToken: (token: string) => void }> = ({
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
      <button type="submit">ログイン</button>
    </form>
  );
};
