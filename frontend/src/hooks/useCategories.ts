import { useQuery } from '@tanstack/react-query';
import apiService from '../services/api';

export const useCategories = () => {
  return useQuery({
    queryKey: ['categories'],
    queryFn: () => apiService.getMarketplaceCategories(),
    retry: 2,
    staleTime: 5 * 60 * 1000, // 5 minutes
    cacheTime: 10 * 60 * 1000, // 10 minutes
  });
};

export const usePublicCategories = () => {
  return useQuery({
    queryKey: ['publicCategories'],
    queryFn: () => apiService.getPublicCategories(),
    retry: 2,
    staleTime: 5 * 60 * 1000, // 5 minutes
    cacheTime: 10 * 60 * 1000, // 10 minutes
  });
};

export const useAgentCategories = () => {
  return useQuery({
    queryKey: ['agentCategories'],
    queryFn: () => apiService.getAgentCategories(),
    retry: 2,
    staleTime: 5 * 60 * 1000, // 5 minutes
    cacheTime: 10 * 60 * 1000, // 10 minutes
  });
};
