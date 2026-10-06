import { useQuery } from '@tanstack/react-query';
import { getEntities } from '../api/entities';

export const useEntities = () => {
  return useQuery({
    queryKey: ['entities'],
    queryFn: getEntities,
  });
};
