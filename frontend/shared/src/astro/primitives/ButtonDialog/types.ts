import type { HTMLAttributes } from "astro/types";
import type { WithPrefix } from "@xrayman/shared/astro/props/props";
import type { ComponentProps } from "astro/types";
import Button from "@xrayman/shared/astro/primitives/Button.astro";

export type BtnProps = ComponentProps<typeof Button>;
export type PageProps = WithPrefix<HTMLAttributes<"div">, "page">;

export interface ButtonDialogProps extends BtnProps, PageProps {}
