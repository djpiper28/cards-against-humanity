import RoundedWhite from "../components/containers/RoundedWhite";
import Header from "../components/typography/Header";
import { createGameUrl } from "../routes";

export default function About() {
  return (
    <section class="flex flex-col gap-3">
      <RoundedWhite>
        <Header text="About - A Party Game For Horrible People" />

        <p>
          This is a free online version of Cards Against Humanity that is
          un-moderated and completely open. It has all of the packs from the
          CahJson project, which is a free community project to collate all of
          the cards. Whilst using this service I do ask you not to ruin the game
          for strangers.
        </p>

        <h2 class="text-xl font-bold">How To Play</h2>
        <ul class="list-disc ps-5">
          <li>
            Go to the{" "}
            <a class="text-blue-500 underline" href={createGameUrl}>
              create game page
            </a>{" "}
            and select some card packs, and a username. Then click create game.
          </li>
          <li>Share the link to your friends.</li>
          <li>
            They all need to select a unique name, and enter the same password,
            or blank if you didn't set one.
          </li>
          <li>When your players have all joined, you can start the game.</li>
        </ul>

        <h2 class="text-xl font-bold">Game Safety</h2>
        <span>
          To help reduce the amount of trolling, it is advised to{" "}
          <span class="font-bold">set a password</span> on your game, and play
          with friends. If you are playing with strangers, you can kick them
          from the game. No game chat is currently provided for your safety.
          <p class="opacity-0"> (and because I am lazy)</p>
          <span>
            Reasonable efforts have been made to stop hackers and cheaters,
            however the hackers gonna hack... If you find a bug or a way to
            exploit the system, please do tell me so I can fix it.
          </span>
        </span>

        <h2 class="text-xl font-bold">Reporting Issues And Contributing</h2>
        <span>
          The project is open source and you can report issues or contribute by
          going go the{" "}
          <a
            href="https://github.com/djpiper28/cards-against-humanity"
            class="w-8 text-blue-600"
          >
            Github
          </a>{" "}
          page.
        </span>
      </RoundedWhite>
    </section>
  );
}
