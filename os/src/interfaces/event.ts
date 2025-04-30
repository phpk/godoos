import mitt from 'mitt';

type Events = {
  refreshDesktop: undefined;
  setLanguages: string;
  recoverData: undefined;
};

export const eventBus = mitt<Events>();