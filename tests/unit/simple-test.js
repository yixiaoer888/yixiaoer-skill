/**
 * simple-test.js - 最简 schema 验证测试
 * 用法: node simple-test.js
 */

const fs   = require('fs');
const path = require('path');
const Ajv  = require('ajv');
const addFormats = require('ajv-formats');

const SCHEMA_DIR = path.join(__dirname, '..', '..', 'schemas', 'platforms');

const ajv = new Ajv({ allErrors: true, strict: false, verbose: true });
addFormats(ajv);

let passed = 0, failed = 0;

function testOne(schemaFile, payload, expectValid, label) {
  const schemaPath = path.join(SCHEMA_DIR, schemaFile);
  if (!fs.existsSync(schemaPath)) {
    console.error(`  ❌ 缺少 schema: ${schemaFile}`);
    failed++;
    return;
  }

  const schema = JSON.parse(fs.readFileSync(schemaPath, 'utf8'));
  // 去除 ajv 不支持的 keyword
  sanitize(schema);

  let validate;
  try {
    validate = ajv.compile(schema);
  } catch (err) {
    console.error(`  ❌ ${label}: ajv compile 失败: ${err.message}`);
    failed++;
    return;
  }

  const valid = validate(payload);
  if (valid === expectValid) {
    console.log(`  ✅ ${label}  (${valid ? '通过' : '按预期拒绝'})`);
    passed++;
  } else {
    console.error(`  ❌ ${label}: 期望 ${expectValid} 实际 ${valid}`);
    if (validate.errors) validate.errors.forEach(e => console.error(`     - ${e.instancePath || '/'} ${e.message}`));
    failed++;
  }
}

function sanitize(obj) {
  if (typeof obj === 'object' && obj !== null && !Array.isArray(obj)) {
    delete obj['elys'];
    for (const v of Object.values(obj)) sanitize(v);
  } else if (Array.isArray(obj)) {
    obj.forEach(sanitize);
  }
}

// ── 测试用例 ────────────────────────────────────────────────────────
console.log('=== 蚁小二 Schema 验证测试 ===\n');

// 1. 抖音视频 - 正确
testOne('douyin.video.schema.json', {
  formType: 'task',
  title: '2026年AI趋势',
  description: '深入探讨AI发展',
  tags: ['AI', '科技']
}, true, '抖音视频(正确)');

// 2. 抖音视频 - 缺 title（应该失败）
testOne('douyin.video.schema.json', {
  formType: 'task',
  description: '缺少标题',
  tags: ['测试']
  // 缺少必填字段 title
}, false, '抖音视频(缺title)');

// 3. 抖音视频 - 多余字段（应该失败，additionalProperties:false）
testOne('douyin.video.schema.json', {
  formType: 'task',
  title: '测试',
  description: '描述',
  tags: ['测试'],
  extraField: 123  // 多余字段
}, false, '抖音视频(多余字段)');

// 4. 小红书视频 - 正确
testOne('xiaohongshu.video.schema.json', {
  formType: 'task',
  title: '好物分享',
  description: '推荐好物',
  visibleType: 0
}, true, '小红书视频(正确)');

// 5. 知乎文章 - 正确
testOne('zhihu.article.schema.json', {
  formType: 'task',
  title: '如何评价AI发展',
  content: '<p>文章内容</p>',
  tags: ['AI']
}, true, '知乎文章(正确)');

// 6. 知乎文章 - 缺 content（应该失败）
testOne('zhihu.article.schema.json', {
  formType: 'task',
  title: '标题'
  // 缺少必填字段 content
}, false, '知乎文章(缺content)');

// 7. 抖音文章 - content 仅空白（schema 会过，业务层应拒绝）
testOne('douyin.article.schema.json', {
  formType: 'task',
  title: '抖音文章标题',
  content: '   ',
  covers: [{ key: 'cover_key', size: 1024, width: 800, height: 600 }],
  visibleType: 0
}, true, '抖音文章(schema层允许空白content)');

console.log(`\n=== 结果: ${passed} 通过, ${failed} 失败 ===`);
process.exit(failed > 0 ? 1 : 0);
