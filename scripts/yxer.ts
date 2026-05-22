import * as fs from 'fs';
import * as path from 'path';
import { callApi, uploadResource, handleError, API_KEY } from './api';
import { SchemaValidator } from './validator';

const validator = new SchemaValidator();
const SKILL_VERSION = '3.0.0';

const HELP_TEXT = `
Usage: yxer <command> [options]

Commands:
  publish <type> <platforms> <payload_path> [clientId]   Publish content
  accounts [platform] [--name <kw>] [--status <n>] [--json]  List accounts
  upload <file_path_or_url> [--bucket <b>]                   Upload resource
  validate <platform> <type> <payload_path>                   Validate payload
  categories <account_id> [--type video|article]             List categories
  locations <account_id> [--query <kw>] [--type <n>]       Search POI locations
  music <account_id> [--query <kw>]                         Search music
  goods <account_id> [--query <kw>]                         List goods
  collections <account_id> [--type video|article]             List collections
  challenges <account_id> [--query <kw>] [--type <t>]      List challenges
  records [--platform <p>] [--limit <n>] [--status <s>]  List publish records
  prepare <platform> <type>                                  Prepare publish data

Options:
  --json       Output as JSON (default: human-readable)
  --debug      Show detailed logs
  --help       Show this help
`;

// ─── helpers ───────────────────────────────────────────────────────────────────

function getFlag(args: string[], flag: string): string | undefined {
  const idx = args.indexOf(flag);
  if (idx !== -1 && idx + 1 < args.length) return args[idx + 1];
  const eq = args.find(a => a.startsWith(flag + '='));
  if (eq) return eq.split('=')[1];
  return undefined;
}

function hasFlag(args: string[], flag: string): boolean {
  return args.includes(flag) || args.some(a => a.startsWith(flag + '='));
}

function getPositionals(args: string[]): string[] {
  const positionals: string[] = [];
  for (let i = 0; i < args.length; i++) {
    const arg = args[i];
    if (arg.startsWith('--')) {
      if (!arg.includes('=') && i + 1 < args.length && !args[i + 1].startsWith('--')) {
        i++;
      }
      continue;
    }
    positionals.push(arg);
  }
  return positionals;
}

function emitSuccess(action: string, data: any) {
  console.log(JSON.stringify({ success: true, action, version: SKILL_VERSION, data }, null, 2));
}

function walk(value: any, visit: (value: any, path: string) => void, currentPath = '$') {
  visit(value, currentPath);
  if (Array.isArray(value)) {
    value.forEach((item, index) => walk(item, visit, `${currentPath}[${index}]`));
    return;
  }
  if (value && typeof value === 'object') {
    Object.entries(value).forEach(([key, child]) => walk(child, visit, `${currentPath}.${key}`));
  }
}

function looksLikeExternalUrl(value: any): boolean {
  return typeof value === 'string' && /^https?:\/\//i.test(value);
}

function requireUploadedResource(resource: any, pathLabel: string, requiredFields: string[], errors: string[]) {
  if (!resource || typeof resource !== 'object') {
    errors.push(`${pathLabel}: missing uploaded resource object`);
    return;
  }
  for (const field of ['key', ...requiredFields]) {
    if (resource[field] === undefined || resource[field] === null || resource[field] === '') {
      errors.push(`${pathLabel}: missing uploaded resource field "${field}"`);
    }
  }
  walk(resource, (value, valuePath) => {
    if (looksLikeExternalUrl(value)) {
      errors.push(`${pathLabel}${valuePath.slice(1)}: external URL is not allowed; run "yxer upload" and use the returned key`);
    }
  });
}

function assertRawObject(value: any, pathLabel: string, errors: string[]) {
  if (Array.isArray(value)) {
    value.forEach((item, index) => assertRawObject(item, `${pathLabel}[${index}]`, errors));
    return;
  }
  if (!value || typeof value !== 'object') return;
  const hasDynamicIdentity = value.yixiaoerId !== undefined || value.yixiaoerName !== undefined || value.id !== undefined || value.name !== undefined;
  if (hasDynamicIdentity && (value.raw === undefined || value.raw === null || typeof value.raw !== 'object')) {
    errors.push(`${pathLabel}: dynamic platform object must include complete "raw" data from a yxer query command`);
  }
}

function preflightPublishPayload(type: string, platforms: string[], rawPayload: any): { accountIds: string[]; errors: string[] } {
  const errors: string[] = [];
  const accountIds: string[] = [];

  if (!['video', 'image-text', 'article'].includes(type)) {
    errors.push(`publish type "${type}" is not supported; expected video, image-text, or article`);
  }
  if (platforms.length === 0 || platforms.some(p => !p)) {
    errors.push('at least one target platform is required');
  }
  if (!rawPayload || typeof rawPayload !== 'object') {
    errors.push('payload must be a JSON object');
    return { accountIds, errors };
  }
  if (!Array.isArray(rawPayload.accountForms) || rawPayload.accountForms.length === 0) {
    errors.push('payload.accountForms must be a non-empty array');
    return { accountIds, errors };
  }

  rawPayload.accountForms.forEach((form: any, index: number) => {
    const formPath = `accountForms[${index}]`;
    const accountId = form?.platformAccountId || form?.account_id;
    if (!accountId) {
      errors.push(`${formPath}: missing platformAccountId`);
    } else {
      accountIds.push(String(accountId));
    }

    const cpf = form?.contentPublishForm;
    if (!cpf || typeof cpf !== 'object') {
      errors.push(`${formPath}: missing contentPublishForm`);
    }

    if (type === 'video') {
      const video = form?.video || cpf?.video;
      requireUploadedResource(video, `${formPath}.video`, ['size', 'width', 'height'], errors);
      const cover = form?.cover || cpf?.cover;
      if (cover) requireUploadedResource(cover, `${formPath}.cover`, ['size', 'width', 'height'], errors);
    }

    if (type === 'image-text') {
      if (!Array.isArray(form?.images) || form.images.length === 0) {
        errors.push(`${formPath}.images: image-text publish requires at least one uploaded image`);
      } else {
        form.images.forEach((image: any, imageIndex: number) => {
          requireUploadedResource(image, `${formPath}.images[${imageIndex}]`, ['size', 'width', 'height'], errors);
        });
      }
      const cover = form?.cover || cpf?.cover;
      if (cover) requireUploadedResource(cover, `${formPath}.cover`, ['size', 'width', 'height'], errors);
    }

    if (type === 'article' && (!cpf?.content || typeof cpf.content !== 'string')) {
      errors.push(`${formPath}.contentPublishForm.content: article publish requires content`);
    }

    walk(form, (value, valuePath) => {
      if (looksLikeExternalUrl(value)) {
        errors.push(`${formPath}${valuePath.slice(1)}: external URL is not allowed in publish payload; upload resources first`);
      }
    });

    ['location', 'music', 'collection', 'collections', 'challenge', 'challenges', 'goods', 'group', 'miniapp'].forEach(field => {
      if (form?.[field] !== undefined) assertRawObject(form[field], `${formPath}.${field}`, errors);
      if (cpf?.[field] !== undefined) assertRawObject(cpf[field], `${formPath}.contentPublishForm.${field}`, errors);
    });
  });

  return { accountIds, errors };
}

async function assertAccountsOnline(platforms: string[], accountIds: string[]) {
  const wanted = new Set(accountIds);
  const found = new Map<string, any>();

  for (const platform of platforms) {
    const result: any = await callApi(`/v2/platform/accounts?platform=${encodeURIComponent(platform)}`, { method: 'GET' });
    const accounts = Array.isArray(result.data) ? result.data : (result.data ? [result.data] : []);
    for (const account of accounts) {
      const id = String(account.platformAccountId || account.id || '');
      if (wanted.has(id)) found.set(id, account);
    }
  }

  const errors: string[] = [];
  for (const accountId of wanted) {
    const account = found.get(accountId);
    if (!account) {
      errors.push(`account ${accountId}: not found in target platform account list`);
      continue;
    }
    if (account.status !== 1) {
      errors.push(`account ${accountId}: status=${account.status}; publish requires status=1`);
    }
  }

  if (errors.length > 0) {
    throw new Error(`Account preflight failed:\n${errors.map(e => `- ${e}`).join('\n')}`);
  }
}

function parseContentType(filePathOrUrl: string): string {
  const ext = path.extname(filePathOrUrl).toLowerCase();
  const map: Record<string, string> = {
    '.jpg': 'image/jpeg', '.jpeg': 'image/jpeg', '.png': 'image/png',
    '.webp': 'image/webp', '.gif': 'image/gif',
    '.mp4': 'video/mp4', '.mov': 'video/quicktime',
  };
  return map[ext] || 'application/octet-stream';
}

function getFileSize(filePath: string): number {
  try { return fs.statSync(filePath).size; } catch { return 0; }
}

function tryGetImageSize(filePath: string): { width?: number; height?: number } {
  try {
    const imgSize = require('image-size');
    const dims = imgSize(filePath);
    return { width: dims.width, height: dims.height };
  } catch { return {}; }
}

// ─── command handlers ─────────────────────────────────────────────────────────

async function handlePublish(args: string[]) {
  const [type, platformStr, payloadPath, clientId] = getPositionals(args);
  if (!type || !platformStr || !payloadPath) {
    throw new Error('Usage: yxer publish <type> <platforms> <payload_path> [clientId]');
  }
  const platforms = platformStr.split(',');
  const absPath = path.resolve(process.cwd(), payloadPath);
  if (!fs.existsSync(absPath)) throw new Error(`Payload file not found: ${absPath}`);

  const rawPayload: any = JSON.parse(fs.readFileSync(absPath, 'utf-8'));
  const preflight = preflightPublishPayload(type, platforms, rawPayload);
  if (preflight.errors.length > 0) {
    throw new Error(`Publish preflight failed:\n${preflight.errors.map(e => `- ${e}`).join('\n')}`);
  }

  console.log(`[yxer] Validating payload for ${platformStr}...`);
  for (const platform of platforms) {
    const { valid, errors } = await validator.validate(platform, type, rawPayload);
    if (!valid) {
      console.error(`[yxer] ❌ Validation failed for ${platform}:`);
      errors?.forEach(e => console.error(`  - ${e}`));
      process.exit(1);
    }
  }

  console.log(`[yxer] Checking account status for ${platformStr}...`);
  await assertAccountsOnline(platforms, preflight.accountIds);

  const payload: any = {
    action: 'publish',
    publishType: type,
    platforms,
    publishArgs: rawPayload,
  };
  if (clientId) {
    payload.publishChannel = 'local';
    payload.clientId = clientId;
  } else {
    payload.publishChannel = 'cloud';
  }

  console.log(`[yxer] Publishing ${type} to ${platformStr}...`);
  const result: any = await callApi('/taskSets/v2', {
    method: 'POST',
    body: JSON.stringify(payload),
  });

  if (hasFlag(args, '--json')) {
    emitSuccess('publish', result);
  } else {
    emitSuccess('publish', result);
  }
}

async function handleAccounts(args: string[]) {
  const platform = getPositionals(args)[0];
  const nameFilter = getFlag(args, '--name');
  const statusFilter = getFlag(args, '--status');
  const json = hasFlag(args, '--json');

  let url = '/v2/platform/accounts';
  const params = new URLSearchParams();
  if (platform) params.append('platform', platform);
  if (params.toString()) url += '?' + params.toString();

  let result: any = await callApi(url, { method: 'GET' });
  let accounts: any[] = Array.isArray(result.data) ? result.data : (result.data ? [result.data] : []);

  if (nameFilter) {
    accounts = accounts.filter((a: any) =>
      (a.name || a.nickname || a.remark || '').includes(nameFilter)
    );
  }
  if (statusFilter !== undefined) {
    const s = parseInt(statusFilter);
    accounts = accounts.filter((a: any) => a.status === s);
  }

  if (json) {
    emitSuccess('accounts', accounts);
  } else {
    console.log(`账号列表${platform ? ` (${platform})` : ''}:`);
    accounts.forEach((a: any, i: number) => {
      const icon = a.status === 1 ? '✅' : '❌';
      const name = a.name || a.nickname || a.remark || '未命名';
      console.log(`  ${i + 1}. ${name} (${a.platformAccountId || a.id || '?'}) ${icon}`);
    });
    if (accounts.length === 0) console.log('  (无在线账号)');
  }
}

async function handleUpload(args: string[]) {
  const filePathOrUrl = getPositionals(args)[0];
  const bucket = getFlag(args, '--bucket') || 'cloud-publish';

  if (!filePathOrUrl) throw new Error('Usage: yxer upload <file_path_or_url> [--bucket cloud-publish|material-library]');

  const contentType = parseContentType(filePathOrUrl);
  let size: number | undefined;

  if (!filePathOrUrl.startsWith('http')) {
    const absPath = path.resolve(process.cwd(), filePathOrUrl);
    size = getFileSize(absPath);
  }

  console.log(`[yxer] Uploading ${filePathOrUrl} (${contentType})...`);
  const key = await uploadResource(filePathOrUrl, bucket, contentType, size);

  const ext = path.extname(filePathOrUrl).toLowerCase();
  const isImage = contentType.startsWith('image/');
  const isVideo = contentType.startsWith('video/');

  const result: any = { success: true, key, contentType, bucket };

  if (isImage) {
    const dims = tryGetImageSize(filePathOrUrl.startsWith('http') ? '' : path.resolve(process.cwd(), filePathOrUrl));
    if (dims.width) result.width = dims.width;
    if (dims.height) result.height = dims.height;
  }

  if (isVideo) {
    // Duration: implement with ffprobe if needed
  }

  if (size) result.size = size;
  result.format = ext.replace('.', '');

  emitSuccess('upload', result);
}

async function handleValidate(args: string[]) {
  let platform: string, type: string, payloadPath: string;

  if (args.length >= 3) {
    [platform, type, payloadPath] = args;
  } else if (args.length === 2) {
    const [platformDotType, pPath] = args;
    const parts = platformDotType.split('.');
    if (parts.length === 2) {
      platform = parts[0];
      type = parts[1];
      payloadPath = pPath;
    } else {
      throw new Error('Usage: yxer validate <platform> <type> <payload_path>');
    }
  } else {
    throw new Error('Usage: yxer validate <platform> <type> <payload_path>');
  }

  const absPath = path.resolve(process.cwd(), payloadPath);
  if (!fs.existsSync(absPath)) throw new Error(`File not found: ${absPath}`);
  const data: any = JSON.parse(fs.readFileSync(absPath, 'utf-8'));

  const { valid, errors } = await validator.validate(platform, type, data);
  if (valid) {
    const preflight = data.accountForms ? preflightPublishPayload(type, [platform], data) : { errors: [] };
    if (preflight.errors.length > 0) {
      console.error(JSON.stringify({ success: false, action: 'validate', version: SKILL_VERSION, errorCode: 'YIXIAOER_USAGE_ERR', message: 'Publish preflight failed', details: preflight.errors }, null, 2));
      process.exit(1);
    }
    emitSuccess('validate', { platform, type, valid: true });
  } else {
    console.error(JSON.stringify({ success: false, action: 'validate', version: SKILL_VERSION, errorCode: 'YIXIAOER_USAGE_ERR', message: 'Schema validation failed', details: errors || [] }, null, 2));
    process.exit(1);
  }
}

async function handleCategories(args: string[]) {
  const accountId = getPositionals(args)[0];
  const type = getFlag(args, '--type') || 'video';
  if (!accountId) throw new Error('Usage: yxer categories <account_id> [--type video|article]');

  const url = `/platform-accounts/${accountId}/categories?publishType=${type}`;
  const result: any = await callApi(url, { method: 'GET' });
  emitSuccess('categories', result.data || result);
}

async function handleLocations(args: string[]) {
  const accountId = getPositionals(args)[0];
  const query = getFlag(args, '--query') || getFlag(args, '--keyword');
  const type = getFlag(args, '--type') || '1';
  if (!accountId) throw new Error('Usage: yxer locations <account_id> [--query 关键词] [--type 0|1|2|3]');

  const params = new URLSearchParams({ locationType: type });
  if (query) params.append('keyWord', query);
  const url = `/platform-accounts/${accountId}/location?${params.toString()}`;
  const result: any = await callApi(url, { method: 'GET' });
  emitSuccess('locations', result.data || result);
}

async function handleMusic(args: string[]) {
  const accountId = getPositionals(args)[0];
  const query = getFlag(args, '--query') || getFlag(args, '--keyword');
  if (!accountId) throw new Error('Usage: yxer music <account_id> [--query 关键词]');

  const params = new URLSearchParams();
  if (query) params.append('keyWord', query);
  const url = `/platform-accounts/${accountId}/music?${params.toString()}`;
  const result: any = await callApi(url, { method: 'GET' });
  emitSuccess('music', result.data || result);
}

async function handleGoods(args: string[]) {
  const accountId = getPositionals(args)[0];
  const query = getFlag(args, '--query') || getFlag(args, '--keyword');
  if (!accountId) throw new Error('Usage: yxer goods <account_id> [--query 关键词]');

  const params = new URLSearchParams();
  if (query) params.append('keyWord', query);
  const url = `/platform-accounts/${accountId}/goods?${params.toString()}`;
  const result: any = await callApi(url, { method: 'GET' });
  emitSuccess('goods', result.data || result);
}

async function handleCollections(args: string[]) {
  const accountId = getPositionals(args)[0];
  const type = getFlag(args, '--type') || 'video';
  if (!accountId) throw new Error('Usage: yxer collections <account_id> [--type video|article]');

  const url = `/platform-accounts/${accountId}/collections?publishType=${type}`;
  const result: any = await callApi(url, { method: 'GET' });
  emitSuccess('collections', result.data || result);
}

async function handleChallenges(args: string[]) {
  const accountId = getPositionals(args)[0];
  const query = getFlag(args, '--query') || getFlag(args, '--keyword');
  const type = getFlag(args, '--type') || 'video';
  if (!accountId) throw new Error('Usage: yxer challenges <account_id> [--query 关键词] [--type video]');

  const params = new URLSearchParams({ publishType: type });
  if (query) params.append('keyWord', query);
  const url = `/platform-accounts/${accountId}/challenges?${params.toString()}`;
  const result: any = await callApi(url, { method: 'GET' });
  emitSuccess('challenges', result.data || result);
}

async function handleRecords(args: string[]) {
  const platform = getFlag(args, '--platform');
  const limit = getFlag(args, '--limit') || '10';
  const status = getFlag(args, '--status');

  const params = new URLSearchParams({ size: limit });
  if (platform) params.append('platform', platform);
  if (status) params.append('status', status);
  const url = `/v2/taskSets?${params.toString()}`;
  const result: any = await callApi(url, { method: 'GET' });
  emitSuccess('records', result.data || result);
}

async function handlePrepare(args: string[]) {
  const positionals = getPositionals(args);
  const platform = positionals[0];
  const type = positionals[1] || 'video';
  if (!platform) throw new Error('Usage: yxer prepare <platform> <type>');

  const accountResult: any = await callApi(`/v2/platform/accounts?platform=${encodeURIComponent(platform)}`, { method: 'GET' });
  const allAccounts = accountResult.data || accountResult;
  const accounts = Array.isArray(allAccounts) ? allAccounts.filter((a: any) => a.status === 1) : [];

  let categories: any = null;
  if (accounts.length > 0 && (type === 'video' || type === 'article')) {
    try {
      const catResult: any = await callApi(`/platform-accounts/${(accounts[0] as any).platformAccountId || (accounts[0] as any).id}/categories?publishType=${type}`, { method: 'GET' });
      categories = catResult.data || catResult;
    } catch {
      // Some platforms don't support categories — ignore
    }
  }

  emitSuccess('prepare', {
    platform,
    type,
    accounts,
    categories,
    defaultFormType: 'task',
    workflow: `workflows/publish-${type}.md`,
    docsIndex: `docs/publish/${type}/index.md`,
    platformDoc: `docs/publish/${type}/${platform}.md`,
    schema: `schemas/platforms/${platform}.${type.replace(/-([a-z])/g, (_: any, c: string) => c.toUpperCase())}.schema.json`,
  });
}

// ─── main ─────────────────────────────────────────────────────────────────────

async function main() {
  const args = process.argv.slice(2);
  const command = args[0];

  if (!command || command === '--help' || command === '-h') {
    console.log(HELP_TEXT);
    return;
  }

  try {
    switch (command) {
      case 'publish':     await handlePublish(args.slice(1)); break;
      case 'accounts':    await handleAccounts(args.slice(1)); break;
      case 'upload':      await handleUpload(args.slice(1)); break;
      case 'validate':    await handleValidate(args.slice(1)); break;
      case 'categories':  await handleCategories(args.slice(1)); break;
      case 'locations':   await handleLocations(args.slice(1)); break;
      case 'music':       await handleMusic(args.slice(1)); break;
      case 'goods':       await handleGoods(args.slice(1)); break;
      case 'collections': await handleCollections(args.slice(1)); break;
      case 'challenges':  await handleChallenges(args.slice(1)); break;
      case 'records':     await handleRecords(args.slice(1)); break;
      case 'prepare':     await handlePrepare(args.slice(1)); break;
      default:
        console.error(`Unknown command: ${command}`);
        console.log(HELP_TEXT);
        process.exit(1);
    }
  } catch (error) {
    handleError(error, `running command: ${command}`);
  }
}

main();
