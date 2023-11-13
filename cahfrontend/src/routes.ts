import { lazy } from "solid-js";
import type { RouteDefinition } from "@solidjs/router";

import Home from "./pages/home";
import About from "./pages/about";
import Create from "./pages/create";

export const routes: RouteDefinition[] = [
  {
    path: "/",
    component: Home,
  },
  {
    path: "/create",
    component: Create,
  },
  {
    path: "/about",
    component: About,
  },
  {
    path: "**",
    component: lazy(() => import("./errors/404")),
  },
];
