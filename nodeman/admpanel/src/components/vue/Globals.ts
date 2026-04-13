import mitt from "mitt";

export enum Colors {
  Background = "Background",
  Foreground = "Foreground",
}

export const Events = {
  SetColor: "set-color",
} as const;

export interface SetColorDetail {
  color: (typeof Colors)[keyof typeof Colors];
  value: string;
}

type Events = {
  [Events.SetColor]: SetColorDetail;
};
const emitter = mitt<Events>();
export default emitter;
