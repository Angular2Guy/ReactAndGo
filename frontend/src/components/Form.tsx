import styles from "../style.module.css";
import {TodoItem1} from "../App";
import { nanoid } from 'nanoid';


interface InputProps {
    todo: string;
    setTodo: (xxx: string) => void;
    todoList: TodoItem1[];
    setTodoList: (xxx: TodoItem1[]) => void;
}

const Form = ({todo, setTodo, todoList, setTodoList}: InputProps) => {
    const handleChange = (event: React.FormEvent<HTMLInputElement>) => {
        setTodo(event.currentTarget.value as string);
    }
    const handleSubmit = (event: React.FormEvent) => {
        event.preventDefault();
        setTodoList([...todoList,{name: todo, id: nanoid()}]);
        setTodo('');
    }
    
    return (<div onSubmit={handleSubmit} className={styles.todoform}>
        <form>
            <input value={todo} onChange={handleChange} className={styles.todoinput} placeholder="Add Todo Item"></input>
            <button type="submit" className={styles.todobutton}>Add</button>
        </form>
    </div>)
}

export default Form;