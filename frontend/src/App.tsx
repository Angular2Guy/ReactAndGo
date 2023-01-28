import React from 'react';
import './App.css';
import Hello from './components/Hello';
import Message from './components/Message';
import Profile from './components/Profile';

function App() {
  return (
    <div className="App">
     <Hello/>
     <Message messageCode='10' messageContent='This is a message from props'/>
     <Profile firstName="Max" lastName='Jones'>
      <p>This is a profile of a person.</p>
      </Profile>
    </div>
  );
}

export default App;
