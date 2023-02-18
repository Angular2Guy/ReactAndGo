//import React, { useState } from 'react';
import './App.scss';
import Header from './components/Header';
//import Form from './components/Form';
//import TodoList from './components/TodoList';
import LoginModal from './components/LoginModal';
import { RecoilRoot } from 'recoil';
import LocationModal from './components/LocationModal';
import TargetPriceModal from './components/TargetPriceModal';

export interface TodoItem1 {
  name: string;
  id: string;
}

function App() {
  //const [todo, setTodo] = useState('');
  //const [todoList, setTodoList] = useState([] as TodoItem1[]);
  return (
    <div className="App">
      <RecoilRoot>
      <Header/>
      <LoginModal/>
      <LocationModal/>
      <TargetPriceModal/>
      {/*
      <Form todo={todo} setTodo={setTodo} todoList={todoList} setTodoList={setTodoList}></Form>
      <TodoList todoList={todoList} setTodoList={setTodoList}></TodoList>
      */}
      </RecoilRoot>
    </div>
  );
}

export default App;
