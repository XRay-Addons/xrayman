import { Settings } from "@/services/api/generated/types.gen";
import type { ConfigData } from "@xrayman/shared/astro/primitives/Config/Config";
import {
  SubscrTitle,
  UsersMessage,
  UserPage,
  TgPage,
  UpdateInterval,
  Routing,
  AppLinks,
  CustomHeaders,
} from "./SettingsFields";

export function toRawData(data: Settings): ConfigData {
  const rd = {
    [SubscrTitle]: data.SubscrTitle,
    [UsersMessage]: data.UsersMessage,
    [UserPage]: data.UserPage,
    [TgPage]: data.TgPage,
    [Routing]: data.Routing,
    [UpdateInterval]: String(data.UpdateInterval),

    [AppLinks]: data.AppLinks.map((app) => ({
      Name: app.Name,
      Platforms: app.Platforms,
      URL: app.URL,
    })),

    [CustomHeaders]: data.CustomHeaders.map((header) => ({
      Key: header.Key,
      Value: header.Value,
    })),
  };
  return rd;
}

export function fromRawData(data: ConfigData): Settings {
  return {
    SubscrTitle: getString(data, SubscrTitle) ?? "",
    UsersMessage: getString(data, UsersMessage) ?? "",
    UserPage: getString(data, UserPage) ?? "",
    TgPage: getString(data, TgPage) ?? "",
    Routing: getString(data, Routing) ?? "",
    UpdateInterval: parseInteger(getString(data, UpdateInterval) ?? "1"),

    AppLinks: getObjectArray(data, AppLinks, {
      Name: "",
      Platforms: "",
      URL: "",
    }),

    CustomHeaders: getObjectArray(data, CustomHeaders, {
      Key: "",
      Value: "",
    }),
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

function getObjectArray<T extends Record<string, string>>(
  data: ConfigData,
  key: string,
  defaults: T,
): T[] {
  const value = data[key];

  if (value === undefined) {
    return [];
  }

  if (!Array.isArray(value)) {
    throw new Error(`Config "${key}" must be an array`);
  }

  return value.map((item, index) => {
    if (typeof item !== "object" || item === null || Array.isArray(item)) {
      throw new Error(`Config "${key}[${index}]" must be an object`);
    }

    return {
      ...defaults,
      ...Object.fromEntries(Object.entries(item).filter(([, value]) => typeof value === "string")),
    };
  });
}

function parseInteger(val: string): number {
  const value = val.trim();

  if (!/^[+-]?\d+$/.test(value)) {
    throw new Error(`Invalid integer: "${val}"`);
  }

  return Number(value);
}
