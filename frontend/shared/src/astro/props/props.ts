// Utilities for forwarding props to nested components.
//
// Example (forwarding props to a nested <input>):
//
// type InputProps = WithPrefix<HTMLAttributes<"input">, "input">;
//
// interface Props extends InputProps, HTMLAttributes<"div"> {}
//
// Usage:
//
// <Element
//   input:type="text"
//   input:placeholder="Name"
//   input={{ autocomplete: "name" }}
//   id="elementID"
// />
//
// Extraction:
//
// const { ...props } = Astro.props;
// const inputProps = ExtractPrefixed(props, "input");
//
// inputProps = {
//   type: "text",
//   placeholder: "Name",
//   autocomplete: "name",
// };
// props = {
//   id: "elementID",
// }

export type WithPrefix<Attrs, Prefix extends string> = {
  [K in keyof Attrs as `${Prefix}:${K & string}`]?: Attrs[K];
} & {
  [K in Prefix]?: Attrs;
};

export function ExtractPrefixed(
  props: Record<string, unknown>,
  prefix: string,
  cutPrefix: boolean = true,
) {
  const prefixed = (props[prefix] ?? {}) as Record<string, unknown>;
  delete props[prefix];
  for (const [key, value] of Object.entries(props)) {
    if (key.startsWith(`${prefix}:`)) {
      prefixed[cutPrefix ? key.slice(prefix.length + 1) : key] = value;
      delete props[key];
    }
  }
  return prefixed;
}
