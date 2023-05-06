import { NO_VALUE } from './constants';

export type Key = number;
export type Value = string;
export type Items = Map<Key, Value>;
export type NoValue = typeof NO_VALUE;

export interface Cache {
  get: (key: Key) => Value | NoValue;
  put: (key: Key, value: Value) => void;
  readonly itemsOrder: Key[];
}
