const FunctionEvent = () => {
    const handleClick = () => {
        console.log('Button is clicked.')
    }
    return (
        <div>
            Functional Component
            <button onClick={handleClick}>Click here</button>
        </div>
    )
}

export default FunctionEvent;