// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-nocheck

import Benchmark from 'benchmark';

import { LRUCache } from './naive-lru-cache';

const CACHE_SIZE = 1000;

function randomPick(arr: Array<any>) {
  const randomIndex = Math.floor(Math.random() * arr.length);

  return arr[randomIndex];
}

function getItemValue(itemKey: number): string {
  return Number(itemKey).toString(16);
}

function setupCache(size: number): { cache: LRUCache, testKeys: number[] } {
  const cache = new LRUCache(size);
  const testKeys: number[] = [];
  
  new Array(CACHE_SIZE).fill(0).forEach((_, index) => {
    cache.put(index, getItemValue(index));
    testKeys.push(index);
  });

  return { cache, testKeys };
}

function testUpdateExistingItemsInCache(cache: LRUCache, testKeys: number[]) {
  // console.log('Running the test')
  new Array(10000).fill(0).forEach(() => {
    const randomKey = randomPick(testKeys);
    cache.put(randomKey, Number(randomKey).toString(16));
  });
}

function testPutNewItemToTheFullCache(cache: LRUCache) {
  new Array(10000).fill(0).forEach((_, index) => {
    const newKey = CACHE_SIZE + index + 1;
    cache.put(newKey, Number(newKey).toString(16));
  });
}

const suite = new Benchmark.Suite();
suite.add(
  'updating existing item in the cache',
  {
    fn: function() { testUpdateExistingItemsInCache(this.cache, this.testKeys) },
    onStart: function() {
      // console.log('Running "before"')
      const { cache, testKeys } = setupCache(CACHE_SIZE);
      
      this.cache = cache;
      this.testKeys = testKeys;
    },
  })
  .add(
    'put a new item in the cache (evict the oldest item)',
    {
      fn: function() { testPutNewItemToTheFullCache(this.cache) },
      onStart: function() {
        // console.log('Running "before"')
        const { cache } = setupCache(CACHE_SIZE);
        
        this.cache = cache;
      },
    }
  )
  .on('complete', function() {
    console.log('Finished running tests:');
    this.forEach(function(test) {
      console.log(`TEST: "${test.name}": Average time: ${test.stats.mean * 1000} milliseconds`);
    });
  })
  .run({ async: true });
