import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faGithub } from "@fortawesome/free-brands-svg-icons";

export default function About() {
  return (
    <section class="flex flex-col gap-3 p-8">
      <h1 class="text-2xl font-bold">About - Cards Against Humanity</h1>

      <h2 class="text-xl">A party game for horrible people</h2>
      <p>
        This is a free online version of Cards Against Humanity that is
        un-moderated and completely open. It has all of the packs from the
        CahJson project, which is a free community project to collate all of the
        cards. Whilst using this service I do ask you not to ruin the game for
        strangers.
      </p>

      <h2 class="text-xl font-bold">Contributing</h2>
      <a
        href="https://github.com/MonarchDevelopment/mtg-search-engine"
        class="w-8 text-black hover:scale-125 flex flex-row gap-3"
      >
        <FontAwesomeIcon icon={faGithub} />
        <p class="text-blue-300">Github</p>
      </a>
    </section>
  );
}
