import { Meta } from "storybook-solidjs";
import { GameLobbyLoaded } from "./GameLobby";
import {
  BlackCard,
  CardPack,
  GameStateInLobby,
  GameStateWhiteCardsBeingSelected,
  WhiteCard,
} from "../../gameLogicTypes";
import { LobbyLoadedProps } from "./gameLoadedProps";
import { P } from "@solid-devtools/debugger/dist/index-9809b963";

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
    navigate: console.log,
  },
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
    navigate: console.log,
  },
};

const yourHand: WhiteCard[] = [
  {
    id: 1,
    bodyText: "Slapping Nigel Farage over and over",
  },
  {
    id: 2,
    bodyText: "Deez nuts",
  },
  {
    id: 3,
    bodyText: "Margret Thatcher's disapproving face",
  },
  {
    id: 4,
    bodyText: "The Queen",
  },
  {
    id: 5,
    bodyText: "The British Empire",
  },
  {
    id: 6,
    bodyText: "Garteth Southgate not swapping out the penalty takers",
  },
  {
    id: 7,
    bodyText: "The NHS",
  },
];

const blackCard: BlackCard = {
  id: 123,
  bodyText: "What is the UK's national sport?",
  cardsToPlay: 1,
};

export const StartedNoCardsPlayed: Meta<LobbyLoadedProps> = {
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
      gameState: GameStateWhiteCardsBeingSelected,
    },
    roundState: {
      yourHand: yourHand,
      yourPlays: [],
      roundNumber: 1,
      currentCardCzarId: "2",
      blackCard: blackCard,
      totalPlays: 0,
    },
    navigate: console.log,
  },
};
