import { cookieStorage } from "@solid-primitives/storage";
import {
  authenticationCookie,
  gameIdParamCookie,
  gamePasswordCookie,
  playerIdCookie,
} from "./gameState";

export default function clearGameCookies() {
  cookieStorage.removeItem(gamePasswordCookie);
  cookieStorage.removeItem(gameIdParamCookie);
  cookieStorage.removeItem(playerIdCookie);
  cookieStorage.removeItem(authenticationCookie);
}
