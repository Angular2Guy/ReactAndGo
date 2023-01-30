import { useState } from "react";

const FunctionalCounter = () => {
    const [counter,setCounter] = useState(0);
    const [text,setText] = useState('Counter value: ');

    const changeValue = (amount: number) => {
        setCounter(counter + amount);
    }
    const increment = () => {
        changeValue(1);
    }
    const decrement = () => {
        changeValue(-1)
    }

    return (<div>
        <div>{text}{counter}</div>
        <div>
            <button onClick={increment}>Increment</button>
            <button onClick={decrement}>Increment</button>
        </div>
    </div>)
}

export default FunctionalCounter;