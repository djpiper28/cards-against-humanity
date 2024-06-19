import type { RouteDefinition } from "@solidjs/router";
import { lazy } from "solid-js";
import Create from "./pages/create";
import Home from "./pages/home";
import Join from "./pages/join";
import PlayerJoin from "./pages/playerJoin";
import GameJoinErrorPage from "./pages/gameError";

export const indexUrl = "/";
export const aboutUrl = "/about";

const gameUrlSegment = "/game";
export const joinGameUrl = `${gameUrlSegment}/join`;
export const createGameUrl = `${gameUrlSegment}/create`;
export const playerJoinUrl = `${gameUrlSegment}/playerJoin`;
export const gameErrorUrl = `${gameUrlSegment}/error`;

/*
 * The routes required for normal gameplay are pre-loaded by default, error states, and other
 * pages are lazy-loaded to reduce the bundle size.
 **/
export const routes: RouteDefinition[] = [
  {
    path: indexUrl,
    component: Home,
  },
  {
    path: createGameUrl,
    component: Create,
  },
  {
    path: aboutUrl,
    component: lazy(() => import("./pages/about")),
  },
  {
    path: joinGameUrl,
    component: Join,
  },
  {
    path: playerJoinUrl,
    component: PlayerJoin,
  },
  {
    path: gameErrorUrl,
    component: GameJoinErrorPage,
  },
  {
    path: "**",
    component: lazy(() => import("./errors/404")),
  },
];
