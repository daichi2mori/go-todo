import { useEffect, useState } from "react";

type Todo = {
  id: number;
  content: string;
  completed: boolean;
  isEditing?: boolean;
};

function App() {
  const [todos, setTodos] = useState<Todo[]>([]);
  const [newTodo, setNewTodo] = useState("");
  const [editContent, setEditContent] = useState("");

  useEffect(() => {
    fetchTodo();
  }, []);

  const fetchTodo = async () => {
    const res = await fetch("http://localhost:8080/todo");
    const todos = await res.json();
    setTodos(todos);
  };

  const addTodo = async () => {
    await fetch("http://localhost:8080/todo", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ content: newTodo, completed: false }),
    });

    setNewTodo("");
    fetchTodo();
  };

  const startEditing = (id: number, content: string) => {
    setTodos(todos.map((todo) => (todo.id === id ? { ...todo, isEditing: true } : todo)));
    setEditContent(content);
  };

  const saveEdit = async (id: number) => {
    await fetch("http://localhost:8080/todo", {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ id: id, content: editContent }),
    });

    setEditContent("");
    fetchTodo();
  };

  const deleteTodo = async (id: number) => {
    await fetch(`http://localhost:8080/todo/${id}`, {
      method: "DELETE",
    });

    fetchTodo();
  };

  return (
    <div>
      <h1>Todo App</h1>
      <input type="text" value={newTodo} onChange={(e) => setNewTodo(e.target.value)} />
      <button onClick={addTodo}>Add Todo</button>
      <ul>
        {todos.map((todo) => (
          <li key={todo.id}>
            {todo.isEditing ? (
              <div>
                <input
                  type="text"
                  value={editContent}
                  onChange={(e) => setEditContent(e.target.value)}
                />
                <button onClick={() => saveEdit(todo.id)}>Save</button>
              </div>
            ) : (
              <div>
                <span>{todo.content}</span>
                <button onClick={() => startEditing(todo.id, todo.content)}>Edit Todo</button>
                <button onClick={() => deleteTodo(todo.id)}>Delete Todo</button>
              </div>
            )}
          </li>
        ))}
      </ul>
    </div>
  );
}

export default App;
