import styles from "../style.module.css";
import {TodoItem1} from "../App";

interface InputProps {
    todoItem: TodoItem1;
    todoList: TodoItem1[];
    setTodoList: (xxx: TodoItem1[]) => void;
}

const Todo = ({todoItem, todoList, setTodoList}: InputProps) => {
    const deleteTodo = () => {
        setTodoList(todoList.filter(item => item.id !== todoItem.id));
    }
    return (
        <div>
            <div className={styles.todoitem}>
                <h3 className={styles.todoname}>{todoItem.name}</h3>
                <button onClick={deleteTodo} className={styles.deletebutton}>Done.</button>
            </div>
        </div>
    )
}

export default Todo;