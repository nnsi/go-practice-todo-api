export const SignupForm: React.FC<{ setToken: (token: string) => void }> = ({
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
      <button type="submit">ユーザー登録</button>
    </form>
  );
};
