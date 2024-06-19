import { Meta } from "storybook-solidjs";
import PlayerCards from "./PlayerCards";

export default {
  component: PlayerCards,
} as Meta;

export const Primary = {
  args: {
    selectedCardIds: ["1", "3"],
    cards: [
      {
        name: "Brexit",
        pack: "CAH",
        played: true,
        id: "1",
      },
      {
        name: "White van man",
        pack: "CAH",
        played: false,
        id: "2",
      },
      {
        name: "Slapping Nigel Farage over and over",
        pack: "CAH",
        played: false,
        id: "3",
      },
      {
        name: "Dead babies",
        pack: "CAH",
        played: false,
        id: "4",
      },
      {
        name: "An endless stream of diarrhea",
        pack: "CAH",
        played: false,
        id: "5",
      },
      {
        name: "A good sniff",
        pack: "CAH",
        played: false,
        id: "6",
      },
      {
        name: "Margret Thatcher",
        pack: "CAH",
        played: false,
        id: "7",
      },
    ],
  },
};
