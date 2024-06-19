import { useEffect, useRef, useState } from "react";

const WebSocketTodoList: React.FC<{ token: string; todos: any[] }> = ({
  token,
}) => {
  const [todos, setWsTodos] = useState([] as any[]);
  const wsRef = useRef<WebSocket | null>(null);
  const wsManualClose = useRef(false);

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
      const ws = new WebSocket(`ws://localhost:8080/ws?token=${token}`);
      wsRef.current = ws;

      wsRef.current.onopen = () => {
        console.log("WebSocket connection opened");
        wsManualClose.current = false;
        wsRef.current?.send(JSON.stringify({ event: "get_todos" }));
      };
      wsRef.current.onmessage = (event) => {
        const message = JSON.parse(event.data);
        const data = JSON.parse(message.data);
        action(message.event, data);
      };
      wsRef.current.onclose = (e) => {
        console.log("WebSocket connection closed", wsManualClose.current, e);
        if (!wsManualClose.current) setTimeout(() => connect(), 3000);
      };
    };
    connect();

    return () => {
      if (wsRef.current) {
        wsManualClose.current = true;
        wsRef.current.close();
      }
    };
  }, []);

  return (
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
              } catch (e) {
                console.error(e);
              }
            }}
          >
            {todo.title}
          </span>
          {todo.completed && "✅"}
          <button
            onClick={async () => {
              try {
                await fetch(`http://localhost:8080/todos/${todo.id}`, {
                  method: "DELETE",
                  headers: {
                    Authorization: `Bearer ${token}`,
                  },
                });
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
      <WebSocketTodoList token={token} todos={todos} />
    </>
  );
};
