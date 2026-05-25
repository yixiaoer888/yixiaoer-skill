/**
 * run-test.js - 用 Node.js 直接测试 schema 验证
 * 用法: node run-test.js
 */

const fs   = require('fs');
const path = require('path');
const Ajv  = require('ajv');
const addFormats = require('ajv-formats');

const SCHEMA_DIR = path.join(__dirname, '..', '..', 'schemas', 'platforms');
const TEST_DIR  = path.join(__dirname, '..', 'fixtures', 'payloads');

const ajv = new Ajv({ allErrors: true, strict: false, verbose: true });
addFormats(ajv);

// ── 测试用例 ────────────────────────────────────────────────────────
const tests = [
  // [schema文件名, 测试文件名, 期望valid]
  ['douyin.video.schema.json',           'douyin-video-valid.json',   true],
  ['douyin.video.schema.json',           'douyin-video-invalid.json', false],
  ['xiaohongshu.video.schema.json',     'xiaohongshu-video-valid.json', true],
  ['zhihu.article.schema.json',         'zhihu-article-valid.json',  true],
  // 额外测试：additionalProperties 拒绝多余字段
  ['douyin.video.schema.json',          'douyin-extra-field.json',   false],
];

// ── 生成测试文件（如果不存在）─────────────────────────────────────
function ensureTestFiles() {
  if (!fs.existsSync(TEST_DIR)) fs.mkdirSync(TEST_DIR, { recursive: true });

  const files = {
    'douyin-video-valid.json': {
      formType: 'task', title: '2026年AI趋势', description: '深入探讨AI发展',
      tags: ['AI','科技']
    },
    'douyin-video-invalid.json': {
      formType: 'task', description: '缺少标题的视频', tags: ['测试']
      // 缺少必填字段 title
    },
    'xiaohongshu-video-valid.json': {
      formType: 'task', title: '好物分享', description: '推荐好物 #好物', visibleType: 0
    },
    'zhihu-article-valid.json': {
      formType: 'task', title: '如何评价AI发展', content: '<p>文章内容</p>', tags: ['AI']
    },
    'douyin-extra-field.json': {
      formType: 'task', title: '有多余字段', description: '测试',
      tags: ['测试'], extraField: 123  // 多余字段，应该被 additionalProperties:false 拒绝
    },
  };

  for (const [name, content] of Object.entries(files)) {
    const fp = path.join(TEST_DIR, name);
    if (!fs.existsSync(fp)) {
      fs.writeFileSync(fp, JSON.stringify(content, null, 2), 'utf8');
      console.log('  ✓ 生成:', name);
    }
  }
}

// ── 主测试逻辑 ────────────────────────────────────────────────────
function run() {
  console.log('=== 蚁小二 Schema 验证测试 ===\n');

  ensureTestFiles();

  let passed = 0, failed = 0;

  for (const [schemaFile, testFile, expectValid] of tests) {
    const schemaPath = path.join(SCHEMA_DIR, schemaFile);
    const testPath   = path.join(TEST_DIR, testFile);

    if (!fs.existsSync(schemaPath)) {
      console.error(`  ❌ 缺少 schema: ${schemaFile}`);
      failed++;
      continue;
    }
    if (!fs.existsSync(testPath)) {
      console.error(`  ❌ 缺少测试文件: ${testFile}`);
      failed++;
      continue;
    }

    const schema = JSON.parse(fs.readFileSync(schemaPath, 'utf8'));
    const data   = JSON.parse(fs.readFileSync(testPath,   'utf8'));

    // 去除 ajv 不支持的 keyword（如有）
    sanitizeSchema(schema);

    let validate;
    try {
      validate = ajv.compile(schema);
    } catch (err) {
      console.error(`  ❌ ${schemaFile}: ajv compile 失败: ${err.message}`);
      failed++;
      continue;
    }

    const valid = validate(data);
    const label = `${schemaFile.replace('.schema.json','')} ← ${testFile}`;

    if (valid === expectValid) {
      console.log(`  ✅ ${label}  (${valid ? '通过' : '按预期拒绝'})`);
      passed++;
    } else {
      console.error(`  ❌ ${label}  期望 ${expectValid} 实际 ${valid}`);
      if (validate.errors) {
        validate.errors.forEach(e => {
          console.error(`     - ${e.instancePath || '/'} ${e.message}  (${JSON.stringify(e.params)})`);
        });
      }
      failed++;
    }
  }

  console.log(`\n=== 结果: ${passed} 通过, ${failed} 失败 ===`);
  process.exit(failed > 0 ? 1 : 0);
}

function sanitizeSchema(obj) {
  if (typeof obj === 'object' && obj !== null && !Array.isArray(obj)) {
    delete obj['elys'];  // 去除不支持的 keyword
    for (const v of Object.values(obj)) sanitizeSchema(v);
  } else if (Array.isArray(obj)) {
    obj.forEach(sanitizeSchema);
  }
}

run();
