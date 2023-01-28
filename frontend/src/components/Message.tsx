import {Component} from 'react';

interface PropsTypes {
    messageContent: string;
    messageCode: string;
}

interface Result {

}

export default class Message extends Component<PropsTypes,Result> {    
    render() {
        return <h1>Message: {this.props.messageContent} Code: {this.props.messageCode}</h1>
    }
}