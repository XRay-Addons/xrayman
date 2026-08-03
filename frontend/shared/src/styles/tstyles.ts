import { type SizeObserver } from "@xrayman/shared/runtime/dom/size-observer";
import { createFluidHandlers } from "./flex-width";

export function initTStyles(so: SizeObserver) {
  for (const [className, fn] of createFluidHandlers()) {
    so.add(className, fn);
  }
}
