import { useEffect, useState } from "react";

const WebSocketTodoList: React.FC<{ token: string; todos: any[] }> = ({
  token,
}) => {
  const [todos, setWsTodos] = useState([] as any[]);
  const [ws, setWs] = useState(null as WebSocket | null);

  const action = (event: string, data: any) => {
    switch (event) {
      case "list":
        setWsTodos(data);
        break;
      case "create":
        setWsTodos((prev) => [...prev, data]);
        break;
      case "update":
        setWsTodos((prev) =>
          prev.map((t) => {
            if (t.id === data.id) {
              return data;
            }
            return t;
          })
        );
        break;
      case "delete":
        setWsTodos((prev) => prev.filter((t) => t.id !== data.id));
        break;
      default:
        break;
    }
  };

  useEffect(() => {
    // Open a WebSocket connection
    const connect = () => {
      const wsInstance = new WebSocket(`ws://localhost:8080/ws?token=${token}`);
      wsInstance.onopen = () => {
        console.log("WebSocket connection opened");
        wsInstance.send(JSON.stringify({ event: "get_todos" }));
      };
      wsInstance.onmessage = (event) => {
        const message = JSON.parse(event.data);
        const data = JSON.parse(message.data);
        action(message.event, data);
      };
      wsInstance.onclose = () => {
        console.log("WebSocket connection closed");
        setWs(null);
        setTimeout(() => connect(), 3000);
      };
    };

    connect();

    return () => {
      if (ws) {
        ws.close();
      }
    };
  }, []);

  return (
    <ul>
      {todos.map((todo: any) => (
        <li key={todo.id}>
          {todo.title} {todo.completed && "✅"}
        </li>
      ))}
    </ul>
  );
};

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
      <hr />
      {todos.length > 0 && <WebSocketTodoList token={token} todos={todos} />}
    </>
  );
};
