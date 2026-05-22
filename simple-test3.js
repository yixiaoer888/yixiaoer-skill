/**
 * simple-test3.js - 最简化的 schema 验证测试
 * 用法: node simple-test3.js
 */

const fs   = require('fs');
const path = require('path');
const Ajv  = require('ajv');
const addFormats = require('ajv-formats');

const SCHEMA_DIR = path.join(__dirname, 'schemas', 'platforms');
const ajv = new Ajv({ allErrors: true, strict: false, verbose: true });
addFormats(ajv);

const cache = {};

function getV(schemaFile) {
  if (cache[schemaFile]) return cache[schemaFile];
  const sp = path.join(SCHEMA_DIR, schemaFile);
  if (!fs.existsSync(sp)) { console.error('  missing schema:', schemaFile); return null; }
  const schema = JSON.parse(fs.readFileSync(sp, 'utf8'));
  delete schema.$id;
  sanitize(schema);
  try { cache[schemaFile] = ajv.compile(schema); return cache[schemaFile]; }
  catch (e) { console.error('  compile fail:', schemaFile, e.message); return null; }
}

function sanitize(obj) {
  if (typeof obj === 'object' && obj !== null && !Array.isArray(obj)) {
    delete obj['elys'];
    for (const v of Object.values(obj)) sanitize(v);
  } else if (Array.isArray(obj)) obj.forEach(sanitize);
}

let p = 0, f = 0;

function T(sf, data, expect, label) {
  const v = getV(sf);
  if (!v) { console.error('  ❌', label, '(no validator)'); f++; return; }
  const ok = v(data);
  if (ok === expect) { console.log('  ✅', label); p++; }
  else {
    console.error('  ❌', label, ' expect:', expect, ' got:', ok);
    if (v.errors) v.errors.forEach(e => console.error('     -', e.instancePath || '/', e.message));
    f++;
  }
}

console.log('=== Schema 验证测试 ===\n');

// ── 抖音视频 ─────────────────────────────
T('douyin.video.schema.json',
  { formType:'task', title:'测试', description:'描述', tags:['AI'], syncComment:false },
  true, '抖音视频(正确)');

T('douyin.video.schema.json',
  { formType:'task', description:'缺title', tags:['A'] },
  false, '抖音视频(缺title)');

T('douyin.video.schema.json',
  { formType:'task', title:'测', description:'描', tags:['A'], extra:1 },
  false, '抖音视频(多余字段)');

// ── 小红书视频 ─────────────────────────
T('xiaohongshu.video.schema.json',
  { formType:'task', title:'好物', description:'推荐', visibleType:0 },
  true, '小红书视频(正确)');

T('xiaohongshu.video.schema.json',
  { formType:'task', visibleType:99 },
  false, '小红书视频(错误visibleType)');

// ── 知乎视频 ──────────────────────────
const zhihuPayload = {
  formType: 'task',
  description: '描述',
  category: [{ yixiaoerId:'1', yixiaoerName:'科技', raw:{} }],
  createType: 1,
  pubType: 1
};
T('zhihu.video.schema.json', zhihuPayload, true, '知乎视频(正确)');

// ── 知乎文章 (需要先创建 schema，暂时注释) ──
// T('zhihu.article.schema.json', { formType:'task', title:'标题', content:'<p>内容</p>' }, true, '知乎文章(正确)');

console.log('\n=== 结果:', p, '通过,', f, '失败 ===');
process.exit(f > 0 ? 1 : 0);
