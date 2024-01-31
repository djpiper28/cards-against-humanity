import type { RouteDefinition } from "@solidjs/router";
import { lazy } from "solid-js";

import Create from "./pages/create";
import Home from "./pages/home";
import Join from "./pages/join";

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
    component: lazy(() => import("./pages/about")),
  },
  {
    path: "/join",
    component: Join,
  },
  {
    path: "**",
    component: lazy(() => import("./errors/404")),
  },
];
