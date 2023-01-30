import React from 'react';
import './App.css';
//import Hello from './components/Hello';
//import Message from './components/Message';
//import Profile from './components/Profile';
//import Counter from './components/Counter';
//import Resume from './components/Resume';
//import FunctionEvent from './components/FunctionEvent';
//import ClassEvent from './components/ClassEvent';
//import FunctionalCounter from './components/FunctionalCounter';
import ConditionalComponent from './components/ConditionalComponent';

function App() {  
  return (
    <div className="App">
      <ConditionalComponent></ConditionalComponent>
     {/*
      <FunctionalCounter></FunctionalCounter>
      <ClassEvent/>
      <FunctionEvent></FunctionEvent>
      <Resume name='Max'></Resume>
      <Profile firstName="Max" lastName='Jones'>
       <p>This is a profile of a person.</p>
       </Profile>
      <Counter/>
     <Hello/>
     <Message messageCode='10' messageContent='This is a message from props'/>
  */}
    </div>
  );
}

export default App;
