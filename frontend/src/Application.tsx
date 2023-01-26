import React from "react";
import ReactDom from "react-dom/client";

function Application() {
    return <div>Application</div>;
}


const domContainer = document.querySelector('#application');
const root = ReactDom.createRoot(domContainer!);
root.render(<Application/>);