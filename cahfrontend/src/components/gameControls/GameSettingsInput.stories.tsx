import { Meta } from "storybook-solidjs";
import GameSettingsInput from "./GameSettingsInput";

export default {
  component: GameSettingsInput,
} as Meta;

export const Primary = {
  args: {
    settings: {
      gamePassword: "password",
      maxPlayers: 3,
      playingToPoints: 10,
      maxRounds: 5,
    },
  },
};
