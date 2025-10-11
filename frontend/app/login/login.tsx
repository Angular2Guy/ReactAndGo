import { NavLink } from "react-router";

export function Login() {
  return (
    <div>
      <h1>Welcome to the App!</h1>
      <p>This is a simple React Router application.</p>
      <NavLink to="/app/app">Application</NavLink>
    </div>
  );
}