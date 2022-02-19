/* Create a Robot class that models a simple robot.
The robot should be able to take "tasks" as chain'able method invocation, and report when a task starts and when it ends.
Start assuming that each task takes 1s to execute. For example:

const wallE = new Robot();
wallE.standUp().walk();
> started standing
> done standing
> start walking paces
> done walking paces

Bonus points if:

* Tasks can take arguments (e.g. wallE.standUp().walk(10)), 
* Tasks can have preconditions (e.g. can only walk if standing).

*/

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

class Robot {
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

const wallE = new Robot();
wallE.walk(5);
wallE.standUp(0).walk(2).sitDown(0).walk(5);
