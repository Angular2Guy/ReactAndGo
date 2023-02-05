import styles from "../style.module.scss";
import {TodoItem1} from "../App";
import React from "react";

interface InputProps {
    todoItem: TodoItem1;
    todoList: TodoItem1[];
    setTodoList: (xxx: TodoItem1[]) => void;
}

const Todo = ({todoItem, todoList, setTodoList}: InputProps) => {
    const deleteTodo = () => {
        setTodoList(todoList.filter(item => item.id !== todoItem.id));
    }    
    //<> starts a fragment and </> ends a fragment to avoid an enclosing <div> pair.
    return (
        <>
            <div className={styles.todoitem}>
                <h3 className={styles.todoname}>{todoItem.name}</h3>
                <button onClick={deleteTodo} className={styles.deletebutton}>Done.</button>
            </div>
        </>
    )
}
//React.memo(...) caches components to prevent rendering if the parent changes
export default React.memo(Todo);