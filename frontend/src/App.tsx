import React, { useState } from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { BrowserRouter } from 'react-router-dom';
import Map from './components/Map';
import Sidebar from './components/Sidebar';
import { useEntities } from './hooks/useEntities';
import { Entity } from './types/entity';

const queryClient = new QueryClient();

const AppContent: React.FC = () => {
  const { data, isLoading, isError, error } = useEntities();
  const [selectedEntity, setSelectedEntity] = useState<Entity | null>(null);
  const [isAdding, setIsAdding] = useState(false);
  const [tempLocation, setTempLocation] = useState<{lat: number; lng: number} | null>(null);

  const entities = data?.data || [];

  // Reset temp location when adding is cancelled
  React.useEffect(() => {
    if (!isAdding) setTempLocation(null);
  }, [isAdding]);

  return (
    <div className="w-full h-screen flex overflow-hidden">
      {/* Left Sidebar */}
      <Sidebar 
        entities={entities} 
        onSelect={setSelectedEntity} 
        selectedEntity={selectedEntity} 
        isAdding={isAdding}
        setIsAdding={setIsAdding}
        tempLocation={tempLocation}
      />

      {/* Main Map Area */}
      <div className="flex-1 relative">
        <div className="absolute inset-0 z-0">
          <Map 
            entities={entities} 
            isAdding={isAdding}
            tempLocation={tempLocation}
            onLocationSelect={setTempLocation}
          />
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
      </div>
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
