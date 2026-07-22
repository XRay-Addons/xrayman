import { type Handler } from "@xrayman/shared/runtime/dom/size-observer";

const breakpoints = [576, 768, 992, 1200] as const;

const fluidWidths = {
  // continious version of col-10 col-md-6 col-lg-4 col-xl-2
  "fluid-sm": [86.7, 50, 33.3, 16.7],
  // continious version of col-12 col-md-8 col-lg-6 col-xl-4
  "fluid-md": [100, 66.7, 50, 33.3],
  // continious version of col-12 col-md-10 col-lg-8 col-xl-8
  "fluid-lg": [100, 83.3, 66.7, 66.7],
} as const;

function interpolate(width: number, values: readonly number[]) {
  // TODO: binary search
  if (width <= values[0]) {
    return values[0];
  }

  for (let i = 1; i < breakpoints.length; i++) {
    const x1 = breakpoints[i - 1];
    const x2 = breakpoints[i];

    if (width <= x2) {
      return values[i - 1] + ((values[i] - values[i - 1]) * (width - x1)) / (x2 - x1);
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
