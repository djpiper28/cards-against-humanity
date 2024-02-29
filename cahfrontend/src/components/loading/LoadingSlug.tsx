import { createSignal } from "solid-js";
import { makeTimer } from "@solid-primitives/timer";

export default function LoadingSlug() {
  const [loadingIndicator, setLoadingIndicator] = createSignal<string>("");
  makeTimer(
    () => {
      let indicator = loadingIndicator();
      indicator += ".";
      if (indicator.length > 3) {
        indicator = "";
      }

      setLoadingIndicator(indicator);
    },
    250,
    setInterval,
  );

  return (
    <span class="font-bold" aria-hidden>
      {loadingIndicator()}
    </span>
  );
}
