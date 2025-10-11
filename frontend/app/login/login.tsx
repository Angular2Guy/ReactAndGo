import { useNavigate } from "react-router";
import Button from '@mui/material/Button';
import * as React from 'react';

export function Login() {
  const navigate = useNavigate();
  
  const navToApp = () => {
    navigate("/app/app");
  }

  return (
    <div>
      <h1>Welcome to the App!</h1>
      <p>This is a simple React Router application.</p>      
      <Button variant="contained" color="primary" onClick={navToApp}>
        Login
      </Button>
    </div>
  );
}
