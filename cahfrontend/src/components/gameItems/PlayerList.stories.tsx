import { Meta } from "storybook-solidjs";
import PlayerList from "./PlayerList";

export default {
  component: PlayerList,
} as Meta;

export const Primary: Meta = {
  args: {
    players: [
      {
        name: "Player 1",
        connected: true,
        hasPlayed: true,
        points: 0,
        id: "1",
      },
      {
        name: "Player 2",
        connected: false,
        hasPlayed: false,
        points: 2,
        id: "2",
      },
      {
        name: "Player 3",
        connected: true,
        hasPlayed: false,
        points: 10,
        id: "3",
      },
      {
        name: "Player 4",
        connected: false,
        hasPlayed: true,
        points: 5,
        id: "4",
      },
    ],
    czarId: "2",
  },
};
