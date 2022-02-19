# JavaScript "Robot"

Create a Robot class that models a simple robot.
The robot should be able to take "tasks" as chain'able method invocation, and report when a task starts and when it ends.
Start assuming that each task takes 1s to execute. For example:

```javascript
const wallE = new Robot();
wallE.standUp().walk();
// outputs:
> started standing
> done standing
> start walking paces
> done walking paces
```

Bonus points if:

* Tasks can take arguments (e.g. wallE.standUp().walk(10)), 
* Tasks can have preconditions (e.g. can only walk if standing).
