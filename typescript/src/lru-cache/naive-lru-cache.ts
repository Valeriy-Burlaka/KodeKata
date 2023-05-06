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
  public storage: {
    items: Items;
    itemsOrder: Key[];
  };
  
  constructor (capacity: number) {
    this.capacity = capacity;
    this.storage = {
      items: new Map(),
      itemsOrder: [],
    };
  }
  
  public get(key: Key): Value | NoValue {
    return this.storage.items.get(key) || NO_VALUE;
  }
  
  public put(key: Key, value: Value) {
    // The item doesn't exist in our cache
    if (this.get(key) < 0) {
      // We have capacity to just add a new item
      if (this.storage.itemsOrder.length < this.capacity) {
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

  private addItem(key: Key, value: Value) {
    this.storage.items.set(key, value);
    this.storage.itemsOrder.push(key);
  }

  private removeOldestItem() {
    const lastItemKey = this.storage.itemsOrder[0];
    this.storage.items.delete(lastItemKey);
    this.storage.itemsOrder.shift();
  }

  private refreshKey(key: Key) {
    const thisKeyPosition = this.storage.itemsOrder.indexOf(key);
    const newItemsOrder = [
      ...this.storage.itemsOrder.slice(0, thisKeyPosition),
      ...this.storage.itemsOrder.slice(thisKeyPosition + 1),
      key,
    ];
    this.storage.itemsOrder = newItemsOrder;
  }

  private updateItem(key: Key, value: Value) {
    this.storage.items.set(key, value);
    this.refreshKey(key);
  }
}
