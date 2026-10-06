import React from 'react';
import { MapContainer, TileLayer, Marker, Popup } from 'react-leaflet';
import 'leaflet/dist/leaflet.css';
import L from 'leaflet';
import icon from 'leaflet/dist/images/marker-icon.png';
import iconShadow from 'leaflet/dist/images/marker-shadow.png';
import { Entity } from '../types/entity';

// Fix Leaflet default marker icon issue in React
const DefaultIcon = L.icon({
  iconUrl: icon,
  shadowUrl: iconShadow,
  iconSize: [25, 41],
  iconAnchor: [12, 41],
});
L.Marker.prototype.options.icon = DefaultIcon;

interface MapProps {
  entities: Entity[];
}

const Map: React.FC<MapProps> = ({ entities }) => {
  return (
    <MapContainer 
      center={[-6.200000, 106.816666]} // Default center to Jakarta
      zoom={11} 
      className="w-full h-full z-0"
    >
      <TileLayer
        attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
      />
      {entities.map((entity) => (
        <Marker 
          key={entity.id} 
          position={[entity.latitude, entity.longitude]}
        >
          <Popup>
            <div className="flex flex-col gap-1">
              <strong className="text-sm font-semibold">{entity.name}</strong>
              <span className="text-xs text-gray-500 capitalize">{entity.type.replace('_', ' ')}</span>
              <span className="text-xs">
                Status: <span className={`px-1 rounded text-white ${
                  entity.status === 'active' ? 'bg-green-500' : 
                  entity.status === 'inactive' ? 'bg-gray-400' :
                  entity.status === 'maintenance' ? 'bg-yellow-500' : 'bg-red-500'
                }`}>{entity.status}</span>
              </span>
            </div>
          </Popup>
        </Marker>
      ))}
    </MapContainer>
  );
};

export default Map;
