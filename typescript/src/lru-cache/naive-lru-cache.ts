import { NO_VALUE } from './constants';
import {
  type Cache,
  type Key,
  type Value,
  type Items,
  type NoValue,
} from './types';

export class LRUCache implements Cache {
  public capacity: number;
  public storage: Items;
  private _itemsOrder: Key[];
  
  constructor (capacity: number) {
    this.capacity = capacity;
    this.storage = new Map();
    this._itemsOrder = [];
  }
  
  public get(key: Key): Value | NoValue {
    return this.storage.get(key) || NO_VALUE;
  }
  
  public put(key: Key, value: Value) {
    // The item doesn't exist in our cache
    if (this.get(key) === NO_VALUE) {
      // We have capacity to just add a new item
      if (this._itemsOrder.length < this.capacity) {
        this.addItem(key, value);
      // We don't have capacity - remove the oldest elem from the cache and add a new item
      } else {
        this.removeOldestItem();
        this.addItem(key, value);
      }
    // The item exists in our cache, so we need to update the value and its order (it becomes last)
    } else {
      this.updateItem(key, value);
    }
  }

  public get itemsOrder(): Key[] {
    return this._itemsOrder;
  }

  private addItem(key: Key, value: Value) {
    this.storage.set(key, value);
    this._itemsOrder.push(key);
  }

  private removeOldestItem() {
    const lastItemKey = this._itemsOrder[0];
    this.storage.delete(lastItemKey);
    this._itemsOrder.shift();
  }

  private refreshKey(key: Key) {
    const thisKeyPosition = this._itemsOrder.indexOf(key);
    const newItemsOrder = [
      ...this._itemsOrder.slice(0, thisKeyPosition),
      ...this._itemsOrder.slice(thisKeyPosition + 1),
      key,
    ];
    this._itemsOrder = newItemsOrder;
  }

  private updateItem(key: Key, value: Value) {
    this.storage.set(key, value);
    this.refreshKey(key);
  }
}
