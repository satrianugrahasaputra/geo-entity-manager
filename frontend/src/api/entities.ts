import client from './client';
import type { ListResponse } from '../types/api';
import type { Entity, EntityPayload } from '../types/entity';

export interface EntityFilters {
  type?: string;
  status?: string;
  search?: string;
}

export const getEntities = async (filters?: EntityFilters): Promise<ListResponse<Entity>> => {
  const { data } = await client.get<ListResponse<Entity>>('/entities', { params: filters });
  return data;
};

export const createEntity = async (payload: EntityPayload): Promise<Entity> => {
  const { data } = await client.post<Entity>('/entities', payload);
  return data;
};

export const updateEntity = async (id: string, payload: EntityPayload): Promise<Entity> => {
  const { data } = await client.put<Entity>(`/entities/${id}`, payload);
  return data;
};

export const deleteEntity = async (id: string): Promise<void> => {
  await client.delete(`/entities/${id}`);
};
