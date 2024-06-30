import {
  GameLobbyState,
  GamePlayerList,
  RpcRoundInformationMsg,
} from "../rpc/rpc-game";
import { Settings } from "./GameSettingsInput";
import { GameLogicCardPack } from "../../api";

export interface LobbyLoadedProps {
  setSettings: (settings: Settings) => void;
  setSelectedPackIds: (ids: string[]) => void;
  players: GamePlayerList;
  commandError: string;
  setCommandError: (error: string) => void;
  dirtyState: boolean;
  cardPacks: GameLogicCardPack[];
  state: GameLobbyState;
  roundState: RpcRoundInformationMsg;
  setStateAsDirty: () => void;
}
