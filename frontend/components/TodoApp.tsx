import { useEffect, useState } from "react";

export const TodoApp: React.FC<{ token: string }> = ({ token }) => {
  const [todos, setTodos] = useState([] as any[]);

  useEffect(() => {
    fetch("http://localhost:8080/todos", {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
      .then((response) => response.json())
      .then((data) => setTodos(data));
  }, []);

  return (
    <>
      <form
        action={async (formData: FormData) => {
          try {
            const req = await fetch("http://localhost:8080/todos", {
              method: "POST",
              headers: {
                "Content-Type": "application/json",
                Authorization: `Bearer ${token}`,
              },
              body: JSON.stringify({ title: formData.get("title") }),
            });
            const data = await req.json();
            setTodos([...todos, data]);
          } catch (e) {
            console.error(e);
          }
          return false;
        }}
      >
        <input type="text" name="title" />
        <button type="submit">Add</button>
      </form>
      <ul>
        {todos.map((todo: any) => (
          <li key={todo.id}>
            <span
              onClick={async () => {
                try {
                  await fetch(`http://localhost:8080/todos/${todo.id}`, {
                    method: "PUT",
                    headers: {
                      "Content-Type": "application/json",
                      Authorization: `Bearer ${token}`,
                    },
                    body: JSON.stringify({ completed: !todo.completed }),
                  });
                  setTodos(
                    todos.map((t) => {
                      if (t.id === todo.id) {
                        return { ...t, completed: !t.completed };
                      }
                      return t;
                    })
                  );
                } catch (e) {
                  console.error(e);
                }
              }}
              style={todo.completed ? { textDecoration: "line-through" } : {}}
            >
              {todo.title}
            </span>
            <button
              onClick={async () => {
                try {
                  await fetch(`http://localhost:8080/todos/${todo.id}`, {
                    method: "DELETE",
                    headers: {
                      Authorization: `Bearer ${token}`,
                    },
                  });
                  setTodos(todos.filter((t) => t.id !== todo.id));
                } catch (e) {
                  console.error(e);
                }
              }}
            >
              DELETE
            </button>
          </li>
        ))}
      </ul>
    </>
  );
};
