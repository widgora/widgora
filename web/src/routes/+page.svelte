<script lang="ts">
	import { checkHealth, type Health } from '#lib/health.ts';

	let health = $state<Health | 'checking'>('checking');

	$effect(() => {
		checkHealth().then((result) => (health = result));
	});
</script>

<svelte:head>
	<title>Widgora</title>
</svelte:head>

<main>
	<h1>Widgora</h1>
	<p>Homelab inventory and documentation. Early development.</p>
	<p>Server: {health}</p>
</main>
