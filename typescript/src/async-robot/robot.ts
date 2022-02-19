
const sleep = (seconds: number) => {
  return new Promise((resolve, _reject) => {
    setTimeout(() => resolve(true), seconds * 1000);
  });
}

type TaskName = 'walking' | 'standing' | 'sitting';
interface Task {
  name: TaskName;
  task: () => Promise<void>;
}

export class Robot {
  queue: Array<Task> = [];
  idle = true;
  state: TaskName = 'sitting';

  startTask (task: Task) {
    if (this.idle) {
      this.idle = false;
      console.time(task.name);
      console.log(`Started ${task.name}`);
      task.task().then(() => this.finishTask(task));
    } else {
      this.queue.push(task);
    }
  }

  finishTask (task: Task) {
    console.log(`Done ${task.name}`);
    console.timeEnd(task.name);
    this.idle = true;
    if (this.queue.length) {
      const nextTask = this.queue.shift();
      nextTask && this.startTask(nextTask);
    }
  }

  sitDown (seconds = 1) {
    this.startTask({
      name: 'sitting',
      task: async () => {
        this.state = 'sitting';
        await sleep(seconds);
      }
    });

    return this;
  }

  standUp (seconds = 1) {
    this.startTask({
      name: 'standing',
      task: async () => {
        this.state = 'standing';
        await sleep(seconds);
      }
    });

    return this;
  }

  walk (seconds = 1) {
    this.startTask({
      name: 'walking',
      task: async () => {
        if (this.state === 'sitting') {
          console.log("Can't walk while sitting. Stand up first.");
          return
        }
        this.state = 'walking';
        await sleep(seconds);
      }
    });

    return this;
  }
}
