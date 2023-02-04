import React, { useState } from 'react';
import './App.css';
import Header from './components/Header';
import Form from './components/Form';
import TodoList from './components/TodoList';

export interface TodoItem1 {
  name: string;
  id: string;
}

function App() {
  const [todo, setTodo] = useState('');
  const [todoList, setTodoList] = useState([] as TodoItem1[]);
  return (
    <div className="App">
      <Header></Header>
      <Form todo={todo} setTodo={setTodo} todoList={todoList} setTodoList={setTodoList}></Form>
      <TodoList todoList={todoList} setTodoList={setTodoList}></TodoList>
    </div>
  );
}

export default App;
