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

interface Task {
  name: string;
  task: () => Promise<void>;
}

class Robot {
  queue: Array<Task> = [];
  idle = true;

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

  standUp (seconds = 1) {
    this.startTask({
      name: 'standing',
      task: async () => {
        await sleep(seconds);
      }
    });

    return this;
  }

  walk (seconds = 1) {
    this.startTask({
      name: 'walking',
      task: async () => {
        await sleep(seconds);
      }
    });

    return this;
  }
}

const wallE = new Robot();
wallE.standUp().walk().walk(5).walk(2).standUp(2);
