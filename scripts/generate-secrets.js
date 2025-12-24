#!/usr/bin/env node

/**
 * Generate Secure Secrets for TelemetryFlow Core
 * 
 * Usage:
 *   node scripts/generate-secrets.js
 *   node scripts/generate-secrets.js --length 64
 *   node scripts/generate-secrets.js --format hex
 */

const crypto = require('crypto');

const args = process.argv.slice(2);
let length = 32;
let format = 'base64';

for (let i = 0; i < args.length; i++) {
  if (args[i] === '--length' && args[i + 1]) {
    length = parseInt(args[i + 1]);
    i++;
  } else if (args[i] === '--format' && args[i + 1]) {
    format = args[i + 1];
    i++;
  } else if (args[i] === '--help' || args[i] === '-h') {
    console.log(`
TelemetryFlow Core - Secure Secret Generator

Usage:
  node scripts/generate-secrets.js [options]

Options:
  --length <number>   Length in bytes (default: 32)
  --format <format>   Output format: base64, hex, base64url (default: base64)
  --help, -h          Show this help

Examples:
  node scripts/generate-secrets.js
  node scripts/generate-secrets.js --length 64
  node scripts/generate-secrets.js --format hex
`);
    process.exit(0);
  }
}

if (length < 32) {
  console.error('❌ Error: Length must be at least 32 bytes');
  process.exit(1);
}

const validFormats = ['base64', 'hex', 'base64url'];
if (!validFormats.includes(format)) {
  console.error(`❌ Error: Format must be one of: ${validFormats.join(', ')}`);
  process.exit(1);
}

function generateSecret(bytes, encoding) {
  const buffer = crypto.randomBytes(bytes);
  if (encoding === 'base64url') {
    return buffer.toString('base64')
      .replace(/\+/g, '-')
      .replace(/\//g, '_')
      .replace(/=/g, '');
  }
  return buffer.toString(encoding);
}

const jwtSecret = generateSecret(length, format);
const sessionSecret = generateSecret(length, format);

console.log('\n🔐 TelemetryFlow Core - Secret Generator');
console.log('=========================================');
console.log(`Length: ${length} bytes | Format: ${format}\n`);

console.log('Generated Secrets:');
console.log('------------------\n');
console.log('JWT_SECRET:');
console.log(`  ${jwtSecret}\n`);
console.log('SESSION_SECRET:');
console.log(`  ${sessionSecret}\n`);

console.log('.env Format:');
console.log('------------');
console.log(`JWT_SECRET=${jwtSecret}`);
console.log(`JWT_EXPIRES_IN=24h`);
console.log(`SESSION_SECRET=${sessionSecret}\n`);

console.log('Docker Example:');
console.log('---------------');
console.log(`docker run -d \\
  -e JWT_SECRET="${jwtSecret}" \\
  -e SESSION_SECRET="${sessionSecret}" \\
  telemetryflow-core:latest\n`);

console.log('Security Tips:');
console.log('--------------');
console.log('✓ Never commit secrets to git');
console.log('✓ Use different secrets per environment');
console.log('✓ Rotate secrets every 90 days');
console.log('✓ Store in secrets manager (AWS Secrets Manager, etc.)\n');

process.exit(0);
