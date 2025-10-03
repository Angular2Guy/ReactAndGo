import { type RouteConfig, index, route } from "@react-router/dev/routes";

export default [
    index("./Login.tsx"),
    route("app", "./App.tsx")
] satisfies RouteConfig;