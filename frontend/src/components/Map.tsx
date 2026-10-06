import React from 'react';
import { MapContainer, TileLayer, Marker, Popup, useMapEvents } from 'react-leaflet';
import 'leaflet/dist/leaflet.css';
import L from 'leaflet';
import icon from 'leaflet/dist/images/marker-icon.png';
import iconShadow from 'leaflet/dist/images/marker-shadow.png';
import type { Entity } from '../types/entity';

// Fix Leaflet default marker icon issue in React
const DefaultIcon = L.icon({
  iconUrl: icon,
  shadowUrl: iconShadow,
  iconSize: [25, 41],
  iconAnchor: [12, 41],
});
L.Marker.prototype.options.icon = DefaultIcon;

// Custom icon for temporary marker
const TempIcon = L.icon({
  iconUrl: 'https://raw.githubusercontent.com/pointhi/leaflet-color-markers/master/img/marker-icon-2x-red.png',
  shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/0.7.7/images/marker-shadow.png',
  iconSize: [25, 41],
  iconAnchor: [12, 41],
});

interface MapProps {
  entities: Entity[];
  isInteractive?: boolean;
  tempLocation?: { lat: number; lng: number } | null;
  onLocationSelect?: (latlng: { lat: number; lng: number }) => void;
}

const MapEvents: React.FC<{ isInteractive?: boolean; onLocationSelect?: (latlng: { lat: number; lng: number }) => void }> = ({ isInteractive, onLocationSelect }) => {
  useMapEvents({
    click(e) {
      if (isInteractive && onLocationSelect) {
        onLocationSelect({ lat: e.latlng.lat, lng: e.latlng.lng });
      }
    }
  });
  return null;
};

const Map: React.FC<MapProps> = ({ entities, isInteractive, tempLocation, onLocationSelect }) => {
  return (
    <MapContainer 
      center={[-6.200000, 106.816666]} // Default center to Jakarta
      zoom={11} 
      className={`w-full h-full z-0 ${isInteractive ? 'cursor-crosshair' : ''}`}
    >
      <TileLayer
        attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
      />
      
      <MapEvents isInteractive={isInteractive} onLocationSelect={onLocationSelect} />

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

      {isInteractive && tempLocation && (
        <Marker 
          position={[tempLocation.lat, tempLocation.lng]} 
          icon={TempIcon} 
          opacity={0.9}
          draggable={true}
          eventHandlers={{
            dragend: (e) => {
              const marker = e.target;
              const position = marker.getLatLng();
              if (onLocationSelect) {
                onLocationSelect({ lat: position.lat, lng: position.lng });
              }
            }
          }}
        >
          <Popup>Lokasi (Geser marker atau klik peta)</Popup>
        </Marker>
      )}
    </MapContainer>
  );
};

export default Map;
