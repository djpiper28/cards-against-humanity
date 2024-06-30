import { Settings } from "./GameSettingsInput";
import { GameLogicCardPack } from "../../api";
import { GamePlayerList } from "../../gameState/gamePlayersList";
import { GameLobbyState } from "../../gameState/gameLobbyState";

export interface LobbyLoadedProps {
  setSettings: (settings: Settings) => void;
  setSelectedPackIds: (ids: string[]) => void;
  players: GamePlayerList;
  commandError: string;
  setCommandError: (error: string) => void;
  dirtyState: boolean;
  cardPacks: GameLogicCardPack[];
  state: GameLobbyState;
  setStateAsDirty: () => void;
}
