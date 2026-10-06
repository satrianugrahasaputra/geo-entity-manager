import React from 'react';
import { Entity } from '../types/entity';
import { MapPin, Info, Edit, Trash2 } from 'lucide-react';

interface SidebarProps {
  entities: Entity[];
  onSelect: (entity: Entity) => void;
  selectedEntity: Entity | null;
}

const Sidebar: React.FC<SidebarProps> = ({ entities, onSelect, selectedEntity }) => {
  return (
    <div className="w-80 h-full bg-white border-r border-gray-200 flex flex-col z-20 shadow-lg">
      <div className="p-4 border-b border-gray-200">
        <h1 className="text-xl font-bold text-gray-800 flex items-center gap-2">
          <MapPin className="text-blue-500" />
          Geo Entities
        </h1>
        <p className="text-sm text-gray-500 mt-1">Daftar entitas geografis</p>
      </div>
      
      <div className="flex-1 overflow-y-auto p-2">
        {entities.length === 0 ? (
          <div className="p-4 text-center text-sm text-gray-500">
            Tidak ada data
          </div>
        ) : (
          <div className="space-y-2">
            {entities.map((entity) => (
              <div 
                key={entity.id}
                onClick={() => onSelect(entity)}
                className={`p-3 rounded-lg border cursor-pointer transition-all ${
                  selectedEntity?.id === entity.id 
                    ? 'border-blue-500 bg-blue-50 shadow-sm' 
                    : 'border-gray-200 hover:border-gray-300 hover:bg-gray-50'
                }`}
              >
                <div className="flex justify-between items-start">
                  <h3 className="font-semibold text-gray-800 text-sm truncate pr-2">{entity.name}</h3>
                  <span className={`px-1.5 py-0.5 rounded text-[10px] font-medium text-white ${
                    entity.status === 'active' ? 'bg-green-500' : 
                    entity.status === 'inactive' ? 'bg-gray-400' :
                    entity.status === 'maintenance' ? 'bg-yellow-500' : 'bg-red-500'
                  }`}>
                    {entity.status}
                  </span>
                </div>
                <p className="text-xs text-gray-500 mt-1 capitalize">{entity.type.replace('_', ' ')}</p>
              </div>
            ))}
          </div>
        )}
      </div>

      {selectedEntity && (
        <div className="border-t border-gray-200 p-4 bg-gray-50 flex flex-col h-1/3">
          <div className="flex justify-between items-center mb-3">
            <h3 className="font-bold text-gray-800 flex items-center gap-1.5">
              <Info size={16} /> Detail
            </h3>
            <div className="flex gap-2">
              <button className="text-blue-600 hover:text-blue-800" title="Edit">
                <Edit size={16} />
              </button>
              <button className="text-red-600 hover:text-red-800" title="Hapus">
                <Trash2 size={16} />
              </button>
            </div>
          </div>
          <div className="flex-1 overflow-y-auto text-sm space-y-2">
            <div>
              <span className="text-gray-500 block text-xs">Nama</span>
              <span className="font-medium">{selectedEntity.name}</span>
            </div>
            <div>
              <span className="text-gray-500 block text-xs">Deskripsi</span>
              <span className="text-gray-700">{selectedEntity.description || '-'}</span>
            </div>
            <div className="grid grid-cols-2 gap-2">
              <div>
                <span className="text-gray-500 block text-xs">Lat</span>
                <span>{selectedEntity.latitude.toFixed(5)}</span>
              </div>
              <div>
                <span className="text-gray-500 block text-xs">Lng</span>
                <span>{selectedEntity.longitude.toFixed(5)}</span>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default Sidebar;
