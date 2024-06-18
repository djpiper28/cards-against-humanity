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

        <h2 class="text-xl font-bold">How to Play</h2>
        <ul class="list-disc ps-5">
          <li>
            Go to the <a href={createGameUrl}>create game page</a> and select
            some card packs, and a username. Then click create game.
          </li>
          <li>Share the link to your friends.</li>
          <li>
            They all need to select a unique name, and enter the same password,
            or blank if you didn't set one.
          </li>
          <li>When your players have all joined, you can start the game.</li>
        </ul>

        <h2 class="text-xl font-bold">Contributing</h2>
        <span>
          The project is open source and you can contribute by going go the{" "}
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
