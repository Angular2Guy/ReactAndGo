import {Component} from 'react';

export default class ClassEvent extends Component {
    handleClick() {
        console.log(this);
        console.log('Button clicked.');
    }

    render() {
        return <div>This is a class based component.
            <button onClick={() => this.handleClick()}>Click here</button>
        </div>
    }
}