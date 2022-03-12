import { RainCollector } from './index';

describe('RainCollector class', () => {

  let collector: RainCollector | null;

  beforeEach(() => {
    collector = new RainCollector();
    collector.addRainRate(10, 40, 200);
    collector.addRainRate(20, 70, 300);
    collector.addRainRate(50, 60, 100);
  });

  afterEach(() => {
    collector = null;
  })

  test('can get a rain rate', () => {
    expect(collector?.getRainRate(0)).toEqual(0);
    expect(collector?.getRainRate(10)).toEqual(200);
    expect(collector?.getRainRate(15)).toEqual(200);
    expect(collector?.getRainRate(30)).toEqual(500);
    expect(collector?.getRainRate(40)).toEqual(500);
    expect(collector?.getRainRate(55)).toEqual(400);
    expect(collector?.getRainRate(80)).toEqual(0);
  });

  test('can get a rain accumulation', () => {
    // expect(collector?.getRainAccumulation(, 60)).toEqual(12000);
    expect(collector?.getRainAccumulation(30, 60)).toEqual(12600);
    expect(collector?.getRainAccumulation(25, 35)).toEqual(5500);
  });
});