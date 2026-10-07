import { describe, expect, it } from 'vitest';
import { entityPayloadSchema } from './entity';

const valid = {
  name: 'Truk 01',
  type: 'vehicle',
  status: 'active',
  description: 'Armada utama',
  latitude: -6.2,
  longitude: 106.8,
};

describe('entityPayloadSchema', () => {
  it('menerima payload valid', () => {
    expect(entityPayloadSchema.safeParse(valid).success).toBe(true);
  });

  it('men-trim name sebelum validasi panjang', () => {
    const ok = entityPayloadSchema.safeParse({ ...valid, name: '  Pos A  ' });
    expect(ok.success).toBe(true);
    expect(ok.data?.name).toBe('Pos A');

    expect(entityPayloadSchema.safeParse({ ...valid, name: '  ab  ' }).success).toBe(false);
  });

  it.each([
    ['name terlalu pendek', { name: 'ab' }],
    ['name terlalu panjang', { name: 'a'.repeat(101) }],
    ['description > 500', { description: 'a'.repeat(501) }],
    ['type tidak dikenal', { type: 'drone' }],
    ['status tidak dikenal', { status: 'broken' }],
    ['latitude < -90', { latitude: -90.0001 }],
    ['latitude > 90', { latitude: 90.0001 }],
    ['longitude < -180', { longitude: -180.0001 }],
    ['longitude > 180', { longitude: 180.0001 }],
  ])('menolak %s', (_label, override) => {
    expect(entityPayloadSchema.safeParse({ ...valid, ...override }).success).toBe(false);
  });

  it.each([
    ['batas bawah', { name: 'abc', latitude: -90, longitude: -180, description: '' }],
    ['batas atas', { name: 'a'.repeat(100), latitude: 90, longitude: 180, description: 'a'.repeat(500) }],
  ])('menerima nilai %s', (_label, override) => {
    expect(entityPayloadSchema.safeParse({ ...valid, ...override }).success).toBe(true);
  });

  it('description opsional', () => {
    const { description: _omit, ...rest } = valid;
    expect(entityPayloadSchema.safeParse(rest).success).toBe(true);
  });
});
