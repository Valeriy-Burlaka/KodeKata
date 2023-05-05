type Key = number;
type Value = string;
type Items = Map<Key, Value>;

export const NoValue = -1;

export class LRUCache {
  public capacity: number;
  public storage: {
    items: Items;
    itemsOrder: Array<Key>;
  };
  
  constructor (capacity: number) {
    this.capacity = capacity;
    this.storage = {
      items: new Map(),
      itemsOrder: [],
    };
  }
  
  public get(key: Key): Value | typeof NoValue {
    return this.storage.items.get(key) || NoValue;
  }
  
  private _add(key: Key, value: Value) {
    this.storage.items.set(key, value);
    this.storage.itemsOrder.push(key);
  }
  
  public put(key: Key, value: Value) {
    // The item doesn't exist in our cache
    if (this.get(key) < 0) {
      // We have capacity to add a new item - add it!
      if (this.storage.itemsOrder.length < this.capacity) {
        this._add(key, value);
      // We don't have capacity - remove the oldest elem from the cache and add a new item
      } else {
        const lastItemKey = this.storage.itemsOrder[0];
        this.storage.itemsOrder.shift();
        this.storage.items.delete(lastItemKey);
        this._add(key, value);
      }
    // The item exists in our cache, so we need to update the value and its order (it becomes last)
    } else {
      const thisKeyPosition = this.storage.itemsOrder.indexOf(key);
      const newItemsOrder = [
        ...this.storage.itemsOrder.slice(0, thisKeyPosition),
        ...this.storage.itemsOrder.slice(thisKeyPosition + 1),
      ];
      this.storage.itemsOrder = newItemsOrder;
      this._add(key, value);
    }

  }
}
