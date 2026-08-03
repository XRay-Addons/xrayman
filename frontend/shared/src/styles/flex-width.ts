import { type Handler } from "@xrayman/shared/runtime/dom/size-observer";

// x1.33 x1.29 x1.21
const breakpoints = [400, 600, 800, 1200, 1800] as const;

const fluidWidths = {
  "fluid-sm": [100.0, 80.0, 60.0, 30.0, 20.0],
  "fluid-md": [100.0, 100.0, 80.0, 50.0, 40.0],
  "fluid-lg": [100.0, 100.0, 80.0, 80.0, 80.0],
} as const;

function interpolate(width: number, values: readonly number[]) {
  // TODO: binary search
  if (width <= values[0]) {
    return values[0];
  }

  for (let i = 1; i < breakpoints.length; i++) {
    const x1 = breakpoints[i - 1];
    const val1 = x1 * values[i - 1];
    const x2 = breakpoints[i];
    const val2 = x2 * values[i];

    if (width <= x2) {
      return (val1 + ((val2 - val1) * (width - x1)) / (x2 - x1)) / width;
    }
  }

  return values.at(-1)!;
}

type FluidHandler = readonly [className: string, handler: Handler];

export function createFluidHandlers(): FluidHandler[] {
  return Object.entries(fluidWidths).map(([className, values]) => [
    className,
    (el: HTMLElement) => {
      el.style.setProperty("--fluid-width", String(interpolate(window.innerWidth, values)));
    },
  ]);
}
