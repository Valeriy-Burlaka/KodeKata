import { NO_VALUE } from './constants';
import {
  type Cache,
  type Key,
  type Value,
  type NoValue,
} from './types';


class LinkedNode {
  constructor(
    public key: Key,
    public value: Value,
    public next: LinkedNode | null = null,
    public prev: LinkedNode | null = null,
    ) {}
}

type Items = Map<Key, LinkedNode>;

export class LRUCache implements Cache {
  // The oldest element is at the start of the queue
  // (i.e., first-in-first-out - an element that hasn't been used for the longest time is first to be evicted from the cache).
  private head: LinkedNode | null = null;
  // The newest element goes to the end of the queue
  private tail: LinkedNode | null = null;
  private storage: Items;

  public capacity: number;
  
  constructor (capacity: number) {
    this.capacity = capacity;
    this.storage = new Map();
  }
  
  public get(key: Key): Value | NoValue {
    return this.storage.get(key)?.value || NO_VALUE;
  }

  public put(key: Key, value: Value) {
    const retrieved = this.storage.get(key);
    // The item isn't in the cache
    if (!retrieved) {
      const newItem = new LinkedNode(key, value);
      // We have free capacity to add a new item w/o removing the oldest item.
      if (this.storage.size < this.capacity) {
        this.addItem(newItem);
      } else {
        this.removeHead();
        this.addItem(newItem);
      }
    // The item is in cache already. We need to update its value and refresh its position
    } else {
      this.updateItem(retrieved, value);
    }
  }

  // It's not needed for this cache impl. but we still need it for testing, to assert correct order of items in the items queue
  public get itemsOrder(): Key[] {
    function* genKeySequence(head: LinkedNode | null) {
      let cursor: LinkedNode | null = head;
      while (cursor) {
        yield cursor.key;
        cursor = cursor.prev;
      }
    }

    return [...genKeySequence(this.head)];
  }

  private addItem(item: LinkedNode) {
    this.storage.set(item.key, item);
    this.addToTail(item);
  }

  private addToTail(item: LinkedNode) {
    // This is the first element in the cache. It becomes a "head" and "tail" of our queue simultaneously.
    // It has no links to other elements yet.
    if (!this.tail) {
      this.tail = item;
      this.head = item;
    // The newest element goes to the end of the queue
    } else {
      this.tail.prev = item;
      item.next = this.tail;
      item.prev = null;
      this.tail = item;
    }
  }

  private removeHead() {
    if (!this.head) return;

    this.storage.delete(this.head.key);
    // We've had only one item in the cache
    if (this.head.prev === null && this.head.next === null) {
      this.head = null;
      this.tail = null;
      return;
    }

    this.head = this.head.prev;
    this.head!.next = null;
  }

  private updateItem(item: LinkedNode, newValue: Value) {
    item.value = newValue;
    this.moveToTail(item);
    this.storage.set(item.key, item);
  }

  private moveToTail(item: LinkedNode) {
    // this is the only item in queue
    if (item.prev === null && item.next === null) {
      return
    // already at tail
    } else if (item.next !== null && item.prev === null) {
      return
    // current head
    } else if (item.prev !== null && item.next === null) {
      item.prev.next = null;
      this.head = item.prev;
      this.addToTail(item);
    // middle item
    } else {
      item.prev!.next = item.next;
      item.next!.prev = item.prev;
      this.addToTail(item);
    }
  }
}