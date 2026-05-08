import { atom } from "nanostores";
import { type User } from "@/services/api/generated/types.gen";
import { type ApiReason } from "@xrayman/shared/services/api/api-reason";
import { State } from "./state";

export type StateSignal =
  | { status: typeof State.Idle }
  | { status: typeof State.LoggedIn; data: User }
  | { status: typeof State.LoggedOut }
  | { status: typeof State.ServerError; reason: ApiReason };

export const stateSignal = atom<StateSignal>({ status: State.Idle });
