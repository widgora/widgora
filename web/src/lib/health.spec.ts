import { describe, expect, it } from 'vitest';
import { checkHealth } from './health';

describe('checkHealth', () => {
	it('reports ok for a 200 response', async () => {
		const fetchFn = async () => new Response('ok\n', { status: 200 });
		expect(await checkHealth(fetchFn as typeof fetch)).toBe('ok');
	});

	it('reports unreachable for an error status', async () => {
		const fetchFn = async () => new Response('', { status: 503 });
		expect(await checkHealth(fetchFn as typeof fetch)).toBe('unreachable');
	});

	it('reports unreachable when the request fails', async () => {
		const fetchFn = async () => {
			throw new TypeError('network down');
		};
		expect(await checkHealth(fetchFn as typeof fetch)).toBe('unreachable');
	});
});
