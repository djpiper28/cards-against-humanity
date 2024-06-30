import { Meta } from "storybook-solidjs";
import { GameLobbyLoaded } from "./GameLobby";
import { CardPack, GameStateInLobby } from "../../gameLogicTypes";
import { LobbyLoadedProps } from "./gameLoadedProps";

export default {
  component: GameLobbyLoaded,
} as Meta;

const cardPacks: CardPack[] = [
  {
    id: "1",
    name: "Pack 1",
    whiteCards: 120,
    blackCards: 130,
  },
  {
    id: "2",
    name: "Pack 2",
    whiteCards: 130,
    blackCards: 140,
  },
  {
    id: "3",
    name: "Pack 3",
    whiteCards: 30,
    blackCards: 40,
  },
];

export const NotStartedNoCardsSelected: Meta<LobbyLoadedProps> = {
  args: {
    setSettings: () => console.log,
    setSelectedPackIds: () => console.log,
    setSelectedCardIds: () => console.log,
    setCommandError: () => console.log,
    setStateAsDirty: () => console.log,
    players: [
      {
        id: "1",
        name: "Player 1",
        connected: true,
        points: 0,
      },
      {
        id: "2",
        name: "Player 2",
        connected: false,
        points: 0,
      },
    ],
    commandError: "",
    dirtyState: false,
    cardPacks: cardPacks,
    state: {
      ownerId: "",
      settings: {
        gamePassword: "",
        maxPlayers: 5,
        cardPacks: [],
        maxRounds: 10,
        playingToPoints: 10,
      },
      creationTime: new Date(),
      gameState: GameStateInLobby,
    },
  },
  navigate: () => console.log,
};

export const NotStartedCardsSelected: Meta<LobbyLoadedProps> = {
  args: {
    setSettings: () => console.log,
    setSelectedPackIds: () => console.log,
    setSelectedCardIds: () => console.log,
    setCommandError: () => console.log,
    setStateAsDirty: () => console.log,
    players: [
      {
        id: "1",
        name: "Player 1",
        connected: true,
        points: 0,
      },
      {
        id: "2",
        name: "Player 2",
        connected: false,
        points: 0,
      },
    ],
    commandError: "",
    dirtyState: false,
    cardPacks: cardPacks,
    state: {
      ownerId: "",
      settings: {
        gamePassword: "",
        maxPlayers: 5,
        cardPacks: cardPacks.map((x) => x.id),
        maxRounds: 10,
        playingToPoints: 10,
      },
      creationTime: new Date(),
      gameState: GameStateInLobby,
    },
  },
  navigate: () => console.log,
};
