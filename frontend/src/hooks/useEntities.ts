import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { getEntities, createEntity, updateEntity, deleteEntity } from '../api/entities';
import type { EntityFilters } from '../api/entities';
import type { EntityPayload } from '../types/entity';

export const useEntities = (filters?: EntityFilters) => {
  return useQuery({
    queryKey: ['entities', filters],
    queryFn: () => getEntities(filters),
  });
};

export const useCreateEntity = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: EntityPayload) => createEntity(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['entities'] });
    },
  });
};

export const useUpdateEntity = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, payload }: { id: string; payload: EntityPayload }) => updateEntity(id, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['entities'] });
    },
  });
};

export const useDeleteEntity = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => deleteEntity(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['entities'] });
    },
  });
};
