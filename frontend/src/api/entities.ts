import client from './client';
import { ListResponse } from '../types/api';
import { Entity, EntityPayload } from '../types/entity';

export const getEntities = async (): Promise<ListResponse<Entity>> => {
  const { data } = await client.get<ListResponse<Entity>>('/entities');
  return data;
};

export const createEntity = async (payload: EntityPayload): Promise<Entity> => {
  const { data } = await client.post<Entity>('/entities', payload);
  return data;
};
