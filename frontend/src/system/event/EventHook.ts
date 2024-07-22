import { useSystem } from '../index.ts';
import { Eventer } from './Eventer';

function initEventer() {
  return new Eventer();
}
function emitEvent(event: string, data?: any) {
  useSystem().emitEvent(event, data);
}

function mountEvent(event: string | string[], callback: (source: string, data: any) => void): void {
  useSystem().mountEvent(event, callback);
}
function redirectEvent(source: string, target: string) {
  mountEvent(source, (_: string, data: unknown) => {
    emitEvent(target, data);
  });
}
export { initEventer, emitEvent, mountEvent, redirectEvent };
