import { lazy } from "solid-js";
import type { RouteDefinition } from "@solidjs/router";

export const routes: RouteDefinition[] = [
  {
    path: "/",
    component: lazy(() => import("./pages/home")),
  },
  {
    path: "/create",
    component: lazy(() => import("./pages/create")),
  },
  {
    path: "/about",
    component: lazy(() => import("./pages/about")),
  },
  {
    path: "**",
    component: lazy(() => import("./errors/404")),
  },
];
