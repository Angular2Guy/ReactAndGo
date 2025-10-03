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