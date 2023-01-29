import { Component } from "react"

interface PropsTypes {
 
}

interface StateTypes {
    counter: number;
}

export default class Counter extends Component<PropsTypes, StateTypes> {
    constructor(props: PropsTypes) {
        super(props);
        this.state = {
            counter: 0
        };
      }

      updateCounter(amount: number) {
        this.setState({
            counter: this.state.counter + amount
        })
      }

    render() {
        return (<div><h3>Count value is: {this.state.counter}</h3>
               <button onClick={() => this.updateCounter(1)}>Increment</button>
               <button onClick={() => this.updateCounter(-1)}>Decrement</button>
               </div>)
    }
}