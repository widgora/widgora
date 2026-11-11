export type Health = 'ok' | 'unreachable';

export async function checkHealth(fetchFn: typeof fetch = fetch): Promise<Health> {
	try {
		const res = await fetchFn('/healthz', { cache: 'no-store' });
		return res.ok ? 'ok' : 'unreachable';
	} catch {
		return 'unreachable';
	}
}
