import { Robot } from './robot';

describe('Robot class', () => {
  test('can create robot', () => {
    const robot = new Robot();
    expect(robot).toHaveProperty('walk');
  });
});
