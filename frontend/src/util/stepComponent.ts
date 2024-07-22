const modules: any = import.meta.glob("@/components/*/*.vue");
const getComponentPath = (name: string) => {
    let model: any;
    for (let p in modules) {
      const n = modules[p].name.split("/").pop().split(".")[0];
      if (n === name) {
        model = p;
        break;
      }
    }
    return model;
};
export const stepComponent = (name: string) => {
    let path = getComponentPath(name);
    if(!path) {
        path = getComponentPath("NotFound")
    }
    return defineAsyncComponent(modules[path]);
};