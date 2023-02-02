import {Component} from "react";

export default class Form extends Component {
    state = {
        firstname: '',
        lastname: '',
    };
    handleChange = (event: React.FormEvent<HTMLInputElement>) => {
        this.setState({firstname:  event.currentTarget.value});        
    }
    handleLastNameChange = (event: React.FormEvent<HTMLInputElement>) => {
        this.setState({lastname: event.currentTarget.value});
    }
    handleSubmit = (event: React.FormEvent) => {
        event.preventDefault();
        console.log({
            fname: this.state.firstname,
            lname: this.state.lastname,
        });
    }
    render() {
        return <div>Form
            <form onSubmit={this.handleSubmit}>
                <input type="text" value={this.state.firstname} onChange={this.handleChange}/>
                <input type="text" value={this.state.lastname} onChange={this.handleLastNameChange}/>
                <button type="submit">Submit</button>
            </form>
        </div>
    }
}