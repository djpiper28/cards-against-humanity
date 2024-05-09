import { Meta } from "storybook-solidjs";
import GameSettingsView from "./GameSettingsView";

export default {
  component: GameSettingsView,
} as Meta;

export const Primar = {
  args: {
    settings: {
      gamePassword: "password",
      maxPlayers: 3,
      playingToPoints: 10,
      maxRounds: 5,
      cardPacks: ["1", "2"],
    },
    packs: [
      { id: "1", name: "Pack 1" },
      { id: "2", name: "Pack 2" },
      { id: "3", name: "Pack 3" },
    ],
  },
};
