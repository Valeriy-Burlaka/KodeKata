import { LRUCache as LRUCacheNaive } from './naive-lru-cache';

import { NO_VALUE } from './constants';
import { type Cache } from './types';

/*
  This is not a good unit-test because testing `.get()` implementation depends on working `.put()` impl. and so on.
  However, the test does its job well as an integration test, assessing the functionality of the cache as a whole.
*/
describe('LRU Cache', () => {
  let cache: Cache;

  beforeEach(() => {
    cache = new LRUCacheNaive(1);
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
    cache = new LRUCacheNaive(3);
    cache.put(1, '1');
    cache.put(2, '2');
    cache.put(3, '3');
    expect(cache.storage.itemsOrder).toMatchObject([1, 2, 3]);

    cache.put(2, 'foo');
    expect(cache.get(2)).toEqual('foo');
    expect(cache.storage.itemsOrder).toMatchObject([1, 3, 2]);
  });
});
