type Key = number;
type Value = string;
export const NoValue = -1;

export class LRUCache {
  public capacity: number;
  public storage: {
    elements: Map<Key, Value>;
    elementsOrder: Array<Key>;
  };
  
  constructor (capacity: number) {
    this.capacity = capacity;
    this.storage = {
      elements: new Map(),
      elementsOrder: [],
    };
  }
  
  public get(key: Key): Value | typeof NoValue {
    return this.storage.elements.get(key) || NoValue;
  }
  
  private _add(key: Key, value: Value) {
    this.storage.elements.set(key, value);
    this.storage.elementsOrder.push(key);
  }
  
  public put(key: Key, value: Value) {
    // The element doesn't exist in our cache
    if (this.get(key) < 0) {
      // We have capacity to add a new element - add it!
      if (this.storage.elementsOrder.length < this.capacity) {
        this._add(key, value);
      // We don't have capacity - remove the oldest elem from the cache and add a new element
      } else {
        const lastElementKey = this.storage.elementsOrder[0];
        this.storage.elementsOrder.shift();
        this.storage.elements.delete(lastElementKey);
        this._add(key, value);
      }
    // The element exists in our cache, so we need to update the value and its order (it becomes last)
    } else {
      const thisKeyPosition = this.storage.elementsOrder.indexOf(key);
      // console.log({ thisKeyPosition, currentOrder: this.storage.elementsOrder})
      const newElementsOrder = [
        ...this.storage.elementsOrder.slice(0, thisKeyPosition),
        ...this.storage.elementsOrder.slice(thisKeyPosition + 1),
      ];
      // console.log({ newElementsOrder })
      this.storage.elementsOrder = newElementsOrder;
      this._add(key, value);
    }

  }
}
