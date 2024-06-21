import {
  startTransition,
  useEffect,
  useOptimistic,
  useRef,
  useState,
} from "react";
import config from "../config";

type WebSocketTodoListProps = {
  token: string;
  todos: (Todo | null)[];
  setTodos: React.Dispatch<React.SetStateAction<(Todo | null)[]>>;
  isConnected: boolean;
  setIsConnected: React.Dispatch<React.SetStateAction<boolean>>;
  setOptimisticTodo: ({ action, todo }: { action: string; todo: Todo }) => void;
};

const WebSocketTodoList: React.FC<WebSocketTodoListProps> = ({
  token,
  todos,
  setTodos,
  isConnected,
  setIsConnected,
  setOptimisticTodo,
}) => {
  const wsRef = useRef<WebSocket | null>(null);
  const wsManualClose = useRef(false);

  const action = (event: string, data: any) => {
    switch (event) {
      case "list":
        setTodos(data);
        break;
      case "create":
        setTodos((prev) => [...prev.filter((t) => t!.id !== data.id), data]);
        break;
      case "update":
        setTodos((prev) => prev.map((t) => (t!.id === data.id ? data : t)));
        break;
      case "delete":
        setTodos((prev) => prev.filter((t) => t!.id !== data.id));
        break;
      default:
        break;
    }
  };

  const updateTodo = async (todo: Todo) => {
    if (!isConnected) return;
    startTransition(async () => {
      try {
        setOptimisticTodo({
          action: "update",
          todo: {
            id: todo.id,
            title: todo.title,
            completed: !todo.completed,
          },
        });
        const req = await fetch(`${config.API_URL}/todos/${todo.id}`, {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({ completed: !todo.completed }),
        });
        const res = await req.json();
        setTodos((prev) => prev.map((t) => (t!.id === res.id ? res : t)));
      } catch (e) {
        console.error(e);
      }
    });
  };

  const deleteTodo = async (todo: Todo) => {
    startTransition(async () => {
      try {
        await fetch(`${config.API_URL}/todos/${todo.id}`, {
          method: "DELETE",
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });
        setTodos((prev) => prev.filter((t) => t!.id !== todo.id));
      } catch (e) {
        console.error(e);
      }
    });
  };

  useEffect(() => {
    const connect = () => {
      if (wsRef.current) {
        wsRef.current.close();
      }
      wsRef.current = new WebSocket(`${config.WS_URL}/ws?token=${token}`);
      wsRef.current.onopen = () => {
        console.log("WebSocket connection opened");
        wsManualClose.current = false;
        setIsConnected(true);
        wsRef.current?.send(JSON.stringify({ event: "get_todos" }));
      };
      wsRef.current.onmessage = (event) => {
        const message = JSON.parse(event.data);
        const data = JSON.parse(message.data);
        action(message.event, data);
      };
      wsRef.current.onclose = (e) => {
        console.log("WebSocket connection closed", wsManualClose.current, e);
        if (!wsManualClose.current) {
          setIsConnected(false);
          console.log("reconnecting");
          setTimeout(() => connect(), 3000);
        } else {
          wsManualClose.current = false;
        }
      };
    };
    connect();

    return () => {
      wsManualClose.current = true;
      wsRef.current?.close();
    };
  }, []);

  return (
    <ul>
      {todos.map((todo: any) => (
        <li key={todo.id}>
          <span>{todo.title}</span>
          {todo.completed && "✅"}
          <button
            onClick={async () => {
              if ("ontouchend" in document) return;
              await updateTodo(todo);
            }}
            onTouchEnd={async () => {
              await updateTodo(todo);
            }}
            disabled={!isConnected}
          >
            {!todo.completed ? "complete" : "uncomplete"}
          </button>
          <button
            onClick={async () => {
              await deleteTodo(todo);
            }}
            disabled={!isConnected}
          >
            DELETE
          </button>
          <button
            onClick={() => {
              console.log(todo);
            }}
          >
            LOG
          </button>
          {todo.sending && "⏳"}
        </li>
      ))}
    </ul>
  );
};

type Todo = {
  id: string;
  title: string;
  completed: boolean;
  sending?: boolean;
};

type TodoAppProps = {
  token: string;
};

export const TodoApp: React.FC<TodoAppProps> = ({ token }) => {
  const [todos, setTodos] = useState<(Todo | null)[]>([]);
  const [isConnected, setIsConnected] = useState(false);
  const [optimisticTodos, setOptimisticTodo] = useOptimistic(
    todos,
    (state, { action, todo }: { action: string; todo: Todo }) => {
      const sendingTodo = {
        ...todo,
        sending: true,
      };
      switch (action) {
        case "create":
          return [...state, sendingTodo];
        case "update":
          return state.map((t) => (t!.id === todo.id ? sendingTodo : t));
        case "delete":
          return state.filter((t) => t!.id !== todo.id);
        default:
          return state;
      }
    }
  );

  return (
    <>
      <form
        action={async (formData: FormData) => {
          try {
            setOptimisticTodo({
              action: "create",
              todo: {
                id: "optimistic_todo",
                title: formData.get("title") as string,
                completed: false,
              },
            });
            await new Promise((res) => setTimeout(res, 1000));
            const req = await fetch(`${config.API_URL}/todos`, {
              method: "POST",
              headers: {
                "Content-Type": "application/json",
                Authorization: `Bearer ${token}`,
              },
              body: JSON.stringify({ title: formData.get("title") }),
            });
            const res = await req.json();
            setTodos((prev) => [...prev.filter((t) => t!.id !== res.id), res]);
          } catch (e) {
            console.error(e);
          }
          return false;
        }}
      >
        <input type="text" name="title" />
        <button type="submit" disabled={!isConnected}>
          Add
        </button>
      </form>
      <WebSocketTodoList
        token={token}
        todos={optimisticTodos}
        setTodos={setTodos}
        isConnected={isConnected}
        setIsConnected={setIsConnected}
        setOptimisticTodo={setOptimisticTodo}
      />
    </>
  );
};
