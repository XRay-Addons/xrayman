import { DynamicConfig } from "@/services/api/generated/types.gen";
import type { ConfigData } from "@xrayman/shared/astro/primitives/Config/Config";
import { SubscrTitle, UsersMessage, UserPage, TgPage, UpdateInterval } from "./ConfigFields";

export function toRawData(data: DynamicConfig): ConfigData {
  return {
    [SubscrTitle]: data.SubscrTitle,
    [UsersMessage]: data.UsersMessage,
    [UserPage]: data.UserPage,
    [TgPage]: data.TgPage,
    [UpdateInterval]: String(data.UpdateInterval),
  };
}

export function fromRawData(data: ConfigData): DynamicConfig {
  return {
    SubscrTitle: getString(data, SubscrTitle) ?? "",
    UsersMessage: getString(data, UsersMessage) ?? "",
    UserPage: getString(data, UserPage) ?? "",
    TgPage: getString(data, TgPage) ?? "",
    UpdateInterval: parseInteger(getString(data, UpdateInterval) ?? "1"),
  };
}

function getString(data: ConfigData, key: string): string | undefined {
  const value = data[key];

  if (value === undefined) {
    return undefined;
  }

  if (typeof value !== "string") {
    throw new Error(`Config "${key}" must be a string`);
  }

  return value;
}

function parseInteger(val: string): number {
  const value = val.trim();

  if (!/^[+-]?\d+$/.test(value)) {
    throw new Error(`Invalid integer: "${val}"`);
  }

  return Number(value);
}
