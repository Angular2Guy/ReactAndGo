import { useState } from "react";
import ListItem from "./ListItem";

export default function Todo() {
    const [todo, setTodo] = useState("");
    const [todoList,setTodoList] = useState([] as string[]);
    const handleChange = (event: React.FormEvent<HTMLInputElement>) => {
        setTodo(event.currentTarget.value);        
    }
    const handleSubmit = (event: React.FormEvent) => {
        event.preventDefault();
        let tempList = todoList;
        tempList.push(todo);
        setTodoList(tempList);    
        setTodo("");    
    }
    return (<div>
        <form onSubmit={handleSubmit}>
            <input value={todo} onChange={handleChange} type="text"></input>
            <button type="submit">Add</button>
        </form>
        {todoList.map((item, index) => (
            <ListItem key={index} name={item}></ListItem>        
        ))}
    </div>);
}