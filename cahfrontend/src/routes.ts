import type { RouteDefinition } from "@solidjs/router";
import { lazy } from "solid-js";
import Create from "./pages/create";
import Home from "./pages/home";
import Join from "./pages/join";
import PlayerJoin from "./pages/playerJoin";
import GameJoinErrorPage from "./pages/gameJoinError";

export const indexUrl = "/";
export const aboutUrl = "/about";

const gameUrlSegment = "/game";
export const joinGameUrl = `${gameUrlSegment}/join`;
export const createGameUrl = `${gameUrlSegment}/create`;
export const playerJoinUrl = `${gameUrlSegment}/playerJoin`;
export const gmaeJoinErrorUrl = `${gameUrlSegment}/joinError`;

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
    path: gmaeJoinErrorUrl,
    component: GameJoinErrorPage,
  },
  {
    path: "**",
    component: lazy(() => import("./errors/404")),
  },
];
