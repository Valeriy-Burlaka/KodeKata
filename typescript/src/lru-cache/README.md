# LRU Cache

A Least Recently Used (LRU) Cache is a cache eviction algorithm that organizes elements in order of use.
In other words, an element that hasn't been used for the longest time will be evicted from the cache when cache capacity is exceeded.

Implement the `LRUCache` class with positive size capacity that has `get` and `put` methods:

```typescript
class LRUCache {

  get(key) {
    // Returns the value of the key if the key is in cache, otherwise return -1
  }

  set(key, value) {
    // 1. If the key doesn't exist and we have a free cache capacity, then add the key-value pair to the cache.
    //    If we don't have a free capacity, evict the least recently used key and then add the key-value pair to the cache.
    // 2. If the key exists, update the value of the key
  }
}
```
