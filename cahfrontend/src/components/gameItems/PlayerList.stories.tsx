import { Meta } from "storybook-solidjs";
import PlayerList from "./PlayerList";

export default {
  component: PlayerList,
} as Meta;

export const Primary = {
  args: {
    players: [
      { name: "Player 1", connected: true },
      { name: "Player 2", connected: false },
      { name: "Player 3", connected: true },
      { name: "Player 4", connected: false },
    ],
  },
};
