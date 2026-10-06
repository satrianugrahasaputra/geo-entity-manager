import React, { useEffect } from 'react';
import { useForm as useHookForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { entityPayloadSchema } from '../schemas/entity';
import type { Entity, EntityPayload } from '../types/entity';

interface EntityFormProps {
  initialLocation: { lat: number; lng: number } | null;
  entityToEdit?: Entity | null;
  onSubmit: (data: EntityPayload) => void;
  onCancel: () => void;
  isLoading: boolean;
}

const EntityForm: React.FC<EntityFormProps> = ({ initialLocation, entityToEdit, onSubmit, onCancel, isLoading }) => {
  const { register, handleSubmit, setValue, reset, formState: { errors } } = useHookForm<EntityPayload>({
    resolver: zodResolver(entityPayloadSchema),
    defaultValues: {
      type: 'vehicle',
      status: 'active',
      description: '',
    }
  });

  useEffect(() => {
    if (entityToEdit) {
      reset({
        name: entityToEdit.name,
        type: entityToEdit.type,
        status: entityToEdit.status,
        description: entityToEdit.description,
        latitude: entityToEdit.latitude,
        longitude: entityToEdit.longitude,
      });
    }
  }, [entityToEdit, reset]);

  // Automatically update lat/lng when user clicks on map
  useEffect(() => {
    if (initialLocation) {
      setValue('latitude', initialLocation.lat, { shouldValidate: true });
      setValue('longitude', initialLocation.lng, { shouldValidate: true });
    }
  }, [initialLocation, setValue]);

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col h-full bg-white text-sm">
      <div className="flex-1 overflow-y-auto p-4 space-y-4">
        
        <div>
          <label className="block text-gray-700 font-medium mb-1">Nama Entitas</label>
          <input 
            type="text" 
            {...register('name')}
            className={`w-full border rounded px-3 py-2 ${errors.name ? 'border-red-500' : 'border-gray-300'}`}
            placeholder="Contoh: Truk Logistik 1"
          />
          {errors.name && <span className="text-red-500 text-xs mt-1 block">{errors.name.message}</span>}
        </div>

        <div>
          <label className="block text-gray-700 font-medium mb-1">Tipe</label>
          <select 
            {...register('type')}
            className="w-full border border-gray-300 rounded px-3 py-2 bg-white"
          >
            <option value="vehicle">Kendaraan (Vehicle)</option>
            <option value="iot_device">Perangkat IoT</option>
            <option value="facility">Fasilitas</option>
            <option value="other">Lainnya</option>
          </select>
          {errors.type && <span className="text-red-500 text-xs mt-1 block">{errors.type.message}</span>}
        </div>

        <div>
          <label className="block text-gray-700 font-medium mb-1">Status</label>
          <select 
            {...register('status')}
            className="w-full border border-gray-300 rounded px-3 py-2 bg-white"
          >
            <option value="active">Aktif</option>
            <option value="inactive">Nonaktif</option>
            <option value="maintenance">Pemeliharaan</option>
            <option value="offline">Offline</option>
          </select>
          {errors.status && <span className="text-red-500 text-xs mt-1 block">{errors.status.message}</span>}
        </div>

        <div>
          <label className="block text-gray-700 font-medium mb-1">Lokasi</label>
          <div className="grid grid-cols-2 gap-2">
            <div>
              <input 
                type="number" step="any"
                {...register('latitude', { valueAsNumber: true })}
                className={`w-full border rounded px-3 py-2 ${errors.latitude ? 'border-red-500' : 'border-gray-300'}`}
                placeholder="Latitude"
                readOnly
              />
            </div>
            <div>
              <input 
                type="number" step="any"
                {...register('longitude', { valueAsNumber: true })}
                className={`w-full border rounded px-3 py-2 ${errors.longitude ? 'border-red-500' : 'border-gray-300'}`}
                placeholder="Longitude"
                readOnly
              />
            </div>
          </div>
          <p className="text-xs text-blue-600 mt-1">Klik di peta untuk memilih lokasi</p>
          {(errors.latitude || errors.longitude) && (
            <span className="text-red-500 text-xs mt-1 block">Lokasi wajib dipilih dari peta</span>
          )}
        </div>

        <div>
          <label className="block text-gray-700 font-medium mb-1">Deskripsi</label>
          <textarea 
            {...register('description')}
            className={`w-full border rounded px-3 py-2 ${errors.description ? 'border-red-500' : 'border-gray-300'}`}
            rows={3}
          ></textarea>
          {errors.description && <span className="text-red-500 text-xs mt-1 block">{errors.description.message}</span>}
        </div>
      </div>

      <div className="p-4 border-t border-gray-200 flex justify-end gap-2 bg-gray-50">
        <button 
          type="button" 
          onClick={onCancel}
          className="px-4 py-2 border border-gray-300 rounded text-gray-700 hover:bg-gray-100"
          disabled={isLoading}
        >
          Batal
        </button>
        <button 
          type="submit"
          className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50"
          disabled={isLoading}
        >
          {isLoading ? 'Menyimpan...' : 'Simpan'}
        </button>
      </div>
    </form>
  );
};

export default EntityForm;
