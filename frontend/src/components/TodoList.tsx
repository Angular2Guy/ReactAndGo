import Todo from "./Todo";
import {TodoItem1} from "../App";

interface InputProps {
    todoList: TodoItem1[];
    setTodoList: (xxx: TodoItem1[]) => void;
}

const TodoList = ({todoList, setTodoList}: InputProps) => {
    return (<div>{todoList.map(todoItem => (
        <Todo key={todoItem.id} todoItem={todoItem} todoList={todoList} setTodoList={setTodoList}></Todo>
    ))}</div>);
}

export default TodoList;