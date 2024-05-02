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
    },
  },
};
