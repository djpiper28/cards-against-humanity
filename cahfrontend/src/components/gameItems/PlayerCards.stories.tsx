import { Meta } from "storybook-solidjs";
import PlayerCards from "./PlayerCards";
import { gameState } from "../../gameState/gameState";

export default {
  component: PlayerCards,
} as Meta;

const cards = [
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
];

export const Primary: Meta = {
  args: {
    selectedCardIds: ["1", "3"],
    cards: cards,
    czarId: gameState.getPlayerId() + "NOT",
  },
};

export const Czar: Meta = {
  args: {
    selectedCardIds: [],
    cards: cards,
    czarId: gameState.getPlayerId(),
  },
};
