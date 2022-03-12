export class RainCollector {
  private rates: { [key: number]: number[] } = {};

  public addRainRate(startTime: number, endTime: number, rate: number): void {
    for (let t = startTime; t <= endTime; t++) {
      if (!this.rates[t]) {
        this.rates[t] = [rate];
      } else {
        this.rates[t].push(rate);
      }
    }
  }

  public getRainRate(time: number): number {
    if (!this.rates[time]) {
      return 0;
    } else {
      return this.rates[time].reduce((a, b) => a + b, 0);
    }
  }

  public getRainAccumulation(startTime: number, endTime: number): number {
    let result = 0;
    for (let t = startTime; t <= endTime; t++) {
      // console.log(`Rain rate at ${t}: ${this.getRainRate(t)}`);
      result += this.getRainRate(t);
      // console.log(`Accumulated result: ${result}`)
    }

    return result;
  }
}