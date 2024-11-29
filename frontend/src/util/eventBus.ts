let eventBus = {
  events: {} as { [key: string]: Function[]},
  on: function (eventName: string, callback: Function) {
    if (!this.events[eventName]) {
      this.events[eventName] = [];
    }
    this.events[eventName].push(callback);
  },
  emit: function (eventName:string, ...args: any[]) {
    if (this.events[eventName]) {
      this.events[eventName].forEach((callback) => {
        callback(...args);
      });
    }
  }
};
export default eventBus;