import { useState } from "react";
//import Counter from './Counter';
//import FunctionalCounter from "./FunctionalCounter";

export default function ConditionalComponent() {
    const [display,setDisplay] = useState(true);

    let output = !!display ? <h3>This is a conditional Component</h3> : <h3>Noting to see here</h3>;    

    return (<div>{output}</div>);

    /*
    if(!!display) {
        return (<div>
            <Counter></Counter>
            </div>);
    } else {
        return (<div>            
            <FunctionalCounter></FunctionalCounter>
            </div>);
    }
    */
}