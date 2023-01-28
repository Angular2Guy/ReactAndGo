/*
export default function Hello() {
    return <h1>Hello World</h1>
}
*/


const name = 'Max';
const displayMessage = (name: string) => {
    return `${name} needs help.`
}
const Hello = () => <h1>The message is: {displayMessage(name)}</h1>;

export default Hello;