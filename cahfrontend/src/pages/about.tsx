export default function About() {
  return (
    <section class="flex flex-col gap-3">
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
    </section>
  );
}
