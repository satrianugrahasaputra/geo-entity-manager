import { z } from 'zod';
import { entitySchema, entityPayloadSchema, entityTypeSchema, entityStatusSchema } from '../schemas/entity';

export type EntityType = z.infer<typeof entityTypeSchema>;
export type EntityStatus = z.infer<typeof entityStatusSchema>;

export type Entity = z.infer<typeof entitySchema>;
export type EntityPayload = z.infer<typeof entityPayloadSchema>;
