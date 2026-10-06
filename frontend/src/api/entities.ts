import client from './client';
import { ListResponse } from '../types/api';
import { Entity } from '../types/entity';

export const getEntities = async (): Promise<ListResponse<Entity>> => {
  const { data } = await client.get<ListResponse<Entity>>('/entities');
  return data;
};
