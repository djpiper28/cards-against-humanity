import type { Component } from "solid-js";
import { Link, useRoutes } from "@solidjs/router";
import { aboutUrl, indexUrl, routes } from "./routes";

const App: Component = () => {
  const Route = useRoutes(routes);

  return (
    <div class="flex flex-col min-h-screen bg-gray-50">
      <nav class="flex flex-row gap-2 bg-gray-200 text-gray-900 px-4 items-center flex-wrap">
        <ul class="flex flex-row flex-grow justify-between items-center py-2 gap-2">
          <li>
            <Link href={indexUrl} class="no-underline hover:underline">
              <span class="font-bold text-3xl">Cards Aginst Humanity</span>
            </Link>
          </li>
          <li>
            <Link href={aboutUrl} class="no-underline hover:underline">
              About
            </Link>
          </li>
        </ul>
      </nav>

      <main class="p-2 md:p-3 flex flex-col w-full gap-5 flex-grow">
        <Route />
      </main>
    </div>
  );
};

export default App;
