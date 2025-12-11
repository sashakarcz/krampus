import apiClient from './client';

export const getEvents = (params = {}) => {
  return apiClient.get('/api/events', { params });
};

export const getPrograms = () => {
  return apiClient.get('/api/programs');
};
