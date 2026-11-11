// Fails when a production dependency of the web UI has a license that is not
// in ../.licenses-allowed. SPDX expressions with OR pass if any alternative is
// allowed; anything containing AND needs every part allowed.
import { execFileSync } from 'node:child_process';
import { readFileSync } from 'node:fs';

const allowed = new Set(
	readFileSync(new URL('../../.licenses-allowed', import.meta.url), 'utf8')
		.split('\n')
		.map((line) => line.trim())
		.filter((line) => line && !line.startsWith('#'))
);

function isAllowed(expression) {
	const ids = expression
		.replace(/[()]/g, ' ')
		.split(/\s+(?:AND|OR)\s+|\s+/)
		.filter(Boolean);
	if (/\sAND\s/.test(expression)) return ids.every((id) => allowed.has(id));
	return ids.some((id) => allowed.has(id));
}

const out = execFileSync('pnpm', ['licenses', 'list', '--prod', '--json'], { encoding: 'utf8' });
const byLicense = out.trim() ? JSON.parse(out) : {};

let count = 0;
const rejected = [];
for (const [license, packages] of Object.entries(byLicense)) {
	for (const pkg of packages) {
		count++;
		if (!isAllowed(license)) rejected.push(`${pkg.name}@${pkg.versions.join(',')}: ${license}`);
	}
}

if (rejected.length > 0) {
	console.error('Licenses not on the allowlist:');
	for (const line of rejected) console.error(`  ${line}`);
	process.exit(1);
}
console.log(`licenses ok: ${count} production packages`);
