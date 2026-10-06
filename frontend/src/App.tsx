import React from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { BrowserRouter } from 'react-router-dom';
import Map from './components/Map';
import { useEntities } from './hooks/useEntities';

const queryClient = new QueryClient();

const AppContent: React.FC = () => {
  const { data, isLoading, isError, error } = useEntities();

  return (
    <div className="relative w-full h-screen flex">
      {/* Map takes full background space */}
      <div className="absolute inset-0 z-0">
        <Map entities={data?.data || []} />
      </div>

      {/* Overlays for loading and error */}
      {isLoading && (
        <div className="absolute top-4 left-1/2 -translate-x-1/2 z-10 bg-white px-4 py-2 rounded-lg shadow font-medium text-sm text-gray-700">
          Memuat data...
        </div>
      )}
      
      {isError && (
        <div className="absolute top-4 left-1/2 -translate-x-1/2 z-10 bg-red-100 border border-red-400 text-red-700 px-4 py-2 rounded-lg shadow font-medium text-sm">
          Gagal memuat data: {(error as any)?.message || 'Terjadi kesalahan'}
        </div>
      )}

      {/* Empty state (when loaded but no data) */}
      {!isLoading && !isError && (!data?.data || data.data.length === 0) && (
        <div className="absolute top-4 left-1/2 -translate-x-1/2 z-10 bg-white px-4 py-2 rounded-lg shadow font-medium text-sm text-gray-700">
          Tidak ada entitas yang ditemukan.
        </div>
      )}
    </div>
  );
};

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <AppContent />
      </BrowserRouter>
    </QueryClientProvider>
  );
}

export default App;
