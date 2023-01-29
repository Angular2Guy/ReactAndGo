import { Component } from "react";

interface PropsTypes {
    name: string
}

interface StateTypes {
    
}

export default class Resume extends Component<PropsTypes,StateTypes> {
    render() {
        const {name} = this.props;
        return <h1>This is a class component {name}</h1>
    }
}