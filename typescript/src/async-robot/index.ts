import { Robot } from './robot';


const wallE = new Robot();
wallE.walk(5);
wallE.standUp(0).walk(2).sitDown(0).walk(5);
