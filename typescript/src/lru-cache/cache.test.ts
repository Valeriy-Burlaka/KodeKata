import { LRUCache, NoValue } from './naive-lru-cache';

/*
  This is not a good unit-test because testing `.get()` implementation depends on working `.put()` impl. and so on.
  However, the test does its job well as an integration test, assessing the functionality of the cache as a whole.
*/
describe('LRU Cache', () => {
  let cache: LRUCache;

  beforeEach(() => {
    cache = new LRUCache(1);
  });

  it('returns no value when key is not in cache', () => {
    expect(cache.get(1)).toEqual(NoValue);
  });

  it('returns a value from cache', () => {
    cache.put(1, '2');
    expect(cache.get(1)).toEqual('2');
  });

  it('evicts the oldest key if cache is full', () => {
    cache.put(1, '2');
    expect(cache.get(1)).toEqual('2');
    
    cache.put(2, 'new');
    expect(cache.get(1)).toEqual(NoValue); // evicted
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
    expect(cache.storage.itemsOrder).toMatchObject([1, 2, 3]);

    cache.put(2, 'foo');
    expect(cache.get(2)).toEqual('foo');
    expect(cache.storage.itemsOrder).toMatchObject([1, 3, 2]);
  });
});
