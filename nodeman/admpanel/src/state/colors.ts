import { ref } from "vue";

export enum Colors {
  BG = "BG",
  Card = "Card",
  Title = "Title",
  Button = "Button",
  Input = "Input",
  Table = "Table",
}

export const colors = ref({
  [Colors.BG]: "undefined",
  [Colors.Card]: "undefined",
  [Colors.Title]: "undefined",
  [Colors.Button]: "undefined",
  [Colors.Input]: "undefined",
  [Colors.Table]: "undefined",
});

export interface SetColorPayload {
  color: (typeof Colors)[keyof typeof Colors];
  value: string;
}

export function setColor(payload: SetColorPayload) {
  colors.value[payload.color] = payload.value;
}

export function setColors(payload: SetColorPayload[]) {
  const edited = colors.value;
  payload.forEach(({ color, value }) => {
    edited[color] = value;
  });
  colors.value = edited;
}
