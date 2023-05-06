// import { LRUCache } from './naive-lru-cache';
import { LRUCache } from './lru-cache';

import { NO_VALUE } from './constants';
import { type Cache } from './types';

/*
  This is not a good unit-test because testing `.get()` implementation depends on working `.put()` impl. and so on.
  However, the test does its job well as an integration test, assessing the functionality of the cache as a whole.
*/
describe('LRU Cache', () => {
  let cache: Cache;

  beforeEach(() => {
    cache = new LRUCache(1);
  });

  it('returns no value when key is not in cache', () => {
    expect(cache.get(1)).toEqual(NO_VALUE);
  });

  it('returns a value from cache', () => {
    cache.put(1, '2');
    expect(cache.get(1)).toEqual('2');
  });

  it('evicts the oldest key if cache is full', () => {
    cache.put(1, '2');
    expect(cache.get(1)).toEqual('2');
    
    cache.put(2, 'new');
    expect(cache.get(1)).toEqual(NO_VALUE); // evicted
    expect(cache.get(2)).toEqual('new');
  });
  
  it('updates the key value', () => {
    cache.put(1, '2');
    expect(cache.get(1)).toEqual('2');

    cache.put(1, '3');
    expect(cache.get(1)).toEqual('3');
  });

  // When a key is updated, it becomes the most recently used key
  it('updates the key value and refreshes the key in queue', () => {
    cache = new LRUCache(3);
    cache.put(1, '1');
    cache.put(2, '2');
    cache.put(3, '3');
    expect(cache.itemsOrder).toMatchObject([1, 2, 3]);
    // Test updating head item
    cache.put(1, 'bar');
    expect(cache.get(1)).toEqual('bar');
    expect(cache.itemsOrder).toMatchObject([2, 3, 1]);
    // Test updating middle item
    cache.put(3, 'foo');
    expect(cache.get(3)).toEqual('foo');
    expect(cache.itemsOrder).toMatchObject([2, 1, 3]);
    // Test updating an item that is already at tail
    cache.put(3, 'foo');
    expect(cache.get(3)).toEqual('foo');
    expect(cache.itemsOrder).toMatchObject([2, 1, 3]);    
  });
});
