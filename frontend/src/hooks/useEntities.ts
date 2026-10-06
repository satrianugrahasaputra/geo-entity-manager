import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { getEntities, createEntity } from '../api/entities';
import { EntityPayload } from '../types/entity';

export const useEntities = () => {
  return useQuery({
    queryKey: ['entities'],
    queryFn: getEntities,
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
