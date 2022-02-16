const defaultGreeting = 'world';

export function hello(greeting: string = defaultGreeting): void {
  console.log(`Hello ${greeting}! `)
}

hello();
hello('Earth');

