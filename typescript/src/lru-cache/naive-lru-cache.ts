type Key = number;
type Value = string;
type NoValue = -1;

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
  
  public get(key: Key): Value | NoValue {
    return this.storage.elements.get(key) || -1;
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

let cache = new LRUCache(1);
console.log(cache.capacity);
console.log(cache.get(0) === -1)

cache.put(1, '2');
console.log(cache.get(1) === '2');

cache.put(2, "2");
console.log(cache.get(1) === -1);
console.log(cache.get(2) === "2");

cache = new LRUCache(3);
cache.put(1, "1");
cache.put(2, "2");
cache.put(3, "3");
console.log(cache.storage.elementsOrder.join(",") === "1,2,3");
// console.log(cache.storage.elementsOrder)

cache.put(2, "foo");
console.log(cache.get(2) === "foo");
// console.log(cache.storage.elementsOrder)
console.log(cache.storage.elementsOrder.join(",") === "1,3,2");