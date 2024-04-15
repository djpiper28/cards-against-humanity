import { Meta } from "@storybook/react";
import CardsSelector from "./CardsSelector";

export default {
  component: CardsSelector,
} as Meta;

export const Primary = {
  args: {
    selectedPackIds: ["1", "2"],
    setSelectedPackIds: (packs: string[]) => console.log(packs),
    cards: [
      {
        id: "1",
        name: "Test Pack 1",
        whiteCards: 10,
        blackCards: 20,
      },
      {
        id: "2",
        name: "Test Pack 2",
        whiteCards: 15,
        blackCards: 25,
      },
      {
        id: "3",
        name: "Test Pack 3",
        whiteCards: 5,
        blackCards: 15,
      },
      {
        Id: "4",
        name: "CAH: Test Pack 4",
        whiteCards: 10,
        blackCards: 0,
      },
    ],
  },
};
