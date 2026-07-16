//import type { HTMLAttributes, HTMLTag } from "astro/types";

export type WithPrefix<Attrs, Prefix extends string> = {
  [K in keyof Attrs as `${Prefix}${K & string}`]?: Attrs[K];
};

/*export type WithPrefix<Tag extends HTMLTag, Prefix extends string> = {
  [K in keyof HTMLAttributes<Tag> as K extends string
    ? string extends K
      ? never
      : `${Prefix}${K}`
    : never]?: HTMLAttributes<Tag>[K];
};*/

/*type KnownKeys<T> = {
  [K in keyof T]: string extends K ? never : K;
}[keyof T];

export type WithPrefix<Tag extends HTMLTag, Prefix extends string> = {
  [K in KnownKeys<HTMLAttributes<Tag>> as K extends string
    ? `${Prefix}${K}`
    : never]: HTMLAttributes<Tag>[K];
};*/

/*export type WithPrefix<Tag extends HTMLTag, Prefix extends string> = {
  [K in keyof HTMLAttributes<Tag> as K extends string
    ? `${Prefix}${K}`
    : never]: HTMLAttributes<Tag>[K];
};*/

/*export type WithPrefix<Tag extends string, Prefix extends string> = {
  [K in keyof HTMLAttributes<Tag> as K extends string
    ? string extends K
      ? never
      : `input:${K}`
    : never]?: HTMLAttributes<Tag>[K];
};*/

/*export type WithPrefix<Tag extends HTMLTag, Prefix extends string> = {
  [K in keyof HTMLAttributes<Tag> as K extends string
    ? string extends K
      ? never
      : `${Prefix}${K}`
    : never]?: HTMLAttributes<Tag>[K];
};*/

export function ExtractPrefixed(props: Record<string, any>, prefix: string) {
  const prefixed: Record<string, any> = {};
  for (const [key, value] of Object.entries(props)) {
    if (key.startsWith(prefix)) {
      prefixed[key.slice(prefix.length)] = value;
      delete props[key];
    }
  }
  return prefixed;
}
