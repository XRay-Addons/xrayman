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
  [Colors.BG]: "#feffa3",
  [Colors.Card]: "#d4ffea",
  [Colors.Title]: "#ffd4e5",
  [Colors.Button]: "#eecbff",
  [Colors.Input]: "#dbdcff",
  [Colors.Table]: "#d4ffea",
});

function foo(x) {
  return "a";
}
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
