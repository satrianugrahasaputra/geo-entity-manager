import { z } from 'zod';

export const entityTypeSchema = z.enum(['vehicle', 'iot_device', 'facility', 'other']);
export const entityStatusSchema = z.enum(['active', 'inactive', 'maintenance', 'offline']);

// Skema untuk entitas utuh yang diterima dari backend
export const entitySchema = z.object({
  id: z.string().uuid(),
  name: z.string().min(3, 'Minimal 3 karakter').max(100, 'Maksimal 100 karakter'),
  type: entityTypeSchema,
  status: entityStatusSchema,
  description: z.string().max(500, 'Maksimal 500 karakter').default(''),
  latitude: z.number().min(-90, 'Harus lebih besar atau sama dengan -90').max(90, 'Harus lebih kecil atau sama dengan 90'),
  longitude: z.number().min(-180, 'Harus lebih besar atau sama dengan -180').max(180, 'Harus lebih kecil atau sama dengan 180'),
  created_at: z.string().datetime(),
  updated_at: z.string().datetime(),
});

// Skema payload untuk Create dan Update
export const entityPayloadSchema = z.object({
  name: z.string().min(3, 'Minimal 3 karakter').max(100, 'Maksimal 100 karakter'),
  type: entityTypeSchema,
  status: entityStatusSchema,
  description: z.string().max(500, 'Maksimal 500 karakter').optional(),
  latitude: z.number({
    required_error: 'Latitude tidak boleh kosong',
    invalid_type_error: 'Harus berupa angka',
  }).min(-90, 'Minimal -90').max(90, 'Maksimal 90'),
  longitude: z.number({
    required_error: 'Longitude tidak boleh kosong',
    invalid_type_error: 'Harus berupa angka',
  }).min(-180, 'Minimal -180').max(180, 'Maksimal 180'),
});
