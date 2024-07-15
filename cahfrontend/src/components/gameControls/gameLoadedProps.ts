import { Settings } from "./GameSettingsInput";
import { GameLogicCardPack } from "../../api";
import { GamePlayerList } from "../../gameState/gamePlayersList";
import { GameLobbyState } from "../../gameState/gameLobbyState";
import { RpcRoundInformationMsg } from "../../rpcTypes";
import { WhiteCard } from "../../gameLogicTypes";

export interface LobbyLoadedProps {
  setSettings: (settings: Settings) => void;
  setSelectedPackIds: (ids: string[]) => void;
  setSelectedCardIds: (ids: string[]) => void;
  setCommandError: (error: string) => void;
  setStateAsDirty: () => void;
  navigate: (url: string) => void;
  players: GamePlayerList;
  commandError: string;
  dirtyState: boolean;
  cardPacks: GameLogicCardPack[];
  state: GameLobbyState;
  roundState: RpcRoundInformationMsg;
  allPlays: WhiteCard[][];
}
