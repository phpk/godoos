import { createORM } from '../orm/orm.ts';

interface Users {
  id?: number;
  username: string;
  password: string;
  nickname: string;
  email: string;
  phone: string;
  isMaster: boolean;
}

export const userDb = createORM<Users>('users');