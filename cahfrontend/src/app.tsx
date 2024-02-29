import type { Component } from "solid-js";
import { Link, useRoutes } from "@solidjs/router";
import { routes } from "./routes";

const App: Component = () => {
  const Route = useRoutes(routes);

  return (
    <div class="flex flex-col min-h-screen bg-gray-50">
      <nav class="flex flex-row gap-2 bg-gray-200 text-gray-900 px-4 items-center flex-wrap">
        <span class="font-bold text-2xl">Cards Aginst Humanity</span>
        <ul class="flex items-center">
          <li class="py-2 px-4">
            <Link href="/" class="no-underline hover:underline">
              Home
            </Link>
          </li>
          <li class="py-2 px-4">
            <Link href="/about" class="no-underline hover:underline">
              About
            </Link>
          </li>
        </ul>
      </nav>

      <main class="p-2 md:p-5 flex flex-col w-full gap-5 flex-grow">
        <Route />
      </main>
    </div>
  );
};

export default App;
