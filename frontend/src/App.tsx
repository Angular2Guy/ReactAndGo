/*
  - Copyright 2022 Sven Loesekann
    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/
import './App.scss';
import Header from './components/Header';
//import Form from './components/Form';
//import TodoList from './components/TodoList';
import LoginModal from './components/LoginModal';
import { RecoilRoot } from 'recoil';
import LocationModal from './components/LocationModal';
import TargetPriceModal from './components/TargetPriceModal';
import Main from './components/Main';

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
      <Main/>
      {/*
      <Form todo={todo} setTodo={setTodo} todoList={todoList} setTodoList={setTodoList}></Form>
      <TodoList todoList={todoList} setTodoList={setTodoList}></TodoList>
      */}
      </RecoilRoot>
    </div>
  );
}

export default App;
