import { cookieStorage } from "@solid-primitives/storage";
import { cookieOptions } from "../apiClient";

export default function clearGameCookies() {
  console.log("Deleting game cookies");
  cookieStorage.clear(cookieOptions);
}
