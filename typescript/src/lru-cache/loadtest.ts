// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-nocheck

import Benchmark from 'benchmark';

import { LRUCache } from './naive-lru-cache';
// import { LRUCache } from './lru-cache';

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
  
  new Array(size).fill(0).forEach((_, index) => {
    cache.put(index, getItemValue(index));
    testKeys.push(index);
  });

  return { cache, testKeys };
}

function testUpdateExistingItemsInCache(cache: LRUCache, testKeys: number[]) {
  new Array(10000).fill(0).forEach(() => {
    const randomKey = randomPick(testKeys);
    cache.put(randomKey, Number(randomKey).toString(16));
  });
}

function testPutNewItemToTheFullCache(cache: LRUCache, testKeys: number[]) {
  new Array(10000).fill(0).forEach((_, index) => {
    const newKey = testKeys.length + index + 1;
    cache.put(newKey, Number(newKey).toString(16));
  });
}

async function runBenchmark(cacheSize, testResults) {
  return new Promise((resolve) => {
    const suite = new Benchmark.Suite();
    suite.add(
      'updating existing item in the cache',
      {
        fn: function() { testUpdateExistingItemsInCache(this.cache, this.testKeys) },
        onStart: function() {
          const { cache, testKeys } = setupCache(cacheSize);

          this.cache = cache;
          this.testKeys = testKeys;
        },
      })
      .add(
        'put a new item in the cache (evict the oldest item)',
        {
          fn: function() { testPutNewItemToTheFullCache(this.cache, this.testKeys) },
          onStart: function() {
            const { cache, testKeys } = setupCache(cacheSize);
          
            this.cache = cache;
            this.testKeys = testKeys;
          },
        }
      )
      .on('complete', function() {
        this.forEach(function(test) {
          if (testResults[test.name]) {
            testResults[test.name][cacheSize] = test.stats.mean * 1000;
          } else {
            testResults[test.name] = {[cacheSize]: test.stats.mean * 1000 }
          }
        });

        resolve();
      })
      .run({ async: true });
  })
}

(async function () {
  const testResults = {};

  const cacheSizes = Array.from({ length: 15 }, (_, i) => (i + 1) * 1000);
  for (const CACHE_SIZE of cacheSizes) {
    console.log('Testing cache size: ', CACHE_SIZE);
     await runBenchmark(CACHE_SIZE, testResults);

  }

  console.log(testResults);
})();
