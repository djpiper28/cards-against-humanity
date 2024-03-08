import { Meta } from "@storybook/react";
import GameSettingsInput from "./GameSettingsInput";

export default {
  component: GameSettingsInput,
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
