import { Meta } from "storybook-solidjs";
import CurrentRoundResults from "./CurrentRoundResults";

export default {
  component: CurrentRoundResults,
} as Meta;

export const InProgress = {
  args: {
    blackCard: {
      name: "What is Batman's guilty pleasure?",
      pack: "CAH",
      cardsToPlay: 1,
    },
    showCards: false,
    plays: [
      {
        winner: false,
        whiteCards: [
          {
            id: "1",
            name: "A windmill full of corpses",
            pack: "CAH",
          },
        ],
      },
      {
        winner: false,
        whiteCards: [
          {
            id: "2",
            name: "Dead babies",
            pack: "CAH",
          },
        ],
      },
      {
        winner: false,
        whiteCards: [
          {
            id: "3",
            name: "A super soaker full of cat pee",
            pack: "CAH",
          },
        ],
      },
    ],
  },
};

export const RoundEnd = {
  args: {
    blackCard: {
      name: "What is Batman's guilty pleasure?",
      pack: "CAH",
      cardsToPlay: 1,
    },
    showCards: true,
    plays: [
      {
        winner: false,
        whiteCards: [
          {
            id: "1",
            name: "A windmill full of corpses",
            pack: "CAH",
          },
        ],
      },
      {
        winner: true,
        whiteCards: [
          {
            id: "2",
            name: "Dead babies",
            pack: "CAH",
          },
        ],
      },
      {
        winner: false,
        whiteCards: [
          {
            id: "3",
            name: "A super soaker full of cat pee",
            pack: "CAH",
          },
        ],
      },
    ],
  },
};

export const RoundEndMultiCard = {
  args: {
    blackCard: {
      name: "_ + _ = profit",
      pack: "CAH",
      cardsToPlay: 2,
    },
    showCards: true,
    plays: [
      {
        winner: false,
        whiteCards: [
          {
            id: "1",
            name: "A windmill full of corpses",
            pack: "CAH",
          },
          {
            id: "11",
            name: "7 dead and 22 injured",
            pack: "CAH",
          },
        ],
      },
      {
        winner: true,
        whiteCards: [
          {
            id: "2",
            name: "Dead babies",
            pack: "CAH",
          },
          {
            id: "22",
            name: "A lifetime of sadness",
            pack: "CAH",
          },
        ],
      },
      {
        winner: false,
        whiteCards: [
          {
            id: "3",
            name: "A metric ton of bricks",
            pack: "CAH",
          },
          {
            id: "33",
            name: "A white van man",
            pack: "CAH",
          },
        ],
      },
    ],
  },
};
