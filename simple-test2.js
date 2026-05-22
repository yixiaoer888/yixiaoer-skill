/**
 * simple-test2.js - 修复 $id 冲突问题
 * 用法: node simple-test2.js
 */

const fs   = require('fs');
const path = require('path');
const Ajv  = require('ajv');
const addFormats = require('ajv-formats');

const SCHEMA_DIR = path.join(__dirname, 'schemas', 'platforms');
const ajv = new Ajv({ allErrors: true, strict: false, verbose: true });
addFormats(ajv);

// 缓存已编译的 validate 函数
const validateCache = {};

function getValidator(schemaFile) {
  if (validateCache[schemaFile]) return validateCache[schemaFile];

  const schemaPath = path.join(SCHEMA_DIR, schemaFile);
  if (!fs.existsSync(schemaPath)) return null;

  const schema = JSON.parse(fs.readFileSync(schemaPath, 'utf8'));
  // 删掉 $id 避免重复编译冲突
  delete schema.$id;
  sanitize(schema);

  try {
    validateCache[schemaFile] = ajv.compile(schema);
    return validateCache[schemaFile];
  } catch (err) {
    console.error(`  ❌ ajv compile 失败 (${schemaFile}): ${err.message}`);
    return null;
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

function testOne(schemaFile, payload, expectValid, label) {
  const validate = getValidator(schemaFile);
  if (!validate) {
    console.error(`  ❌ 无法加载 schema: ${schemaFile}`);
    return false;
  }

  const valid = validate(payload);
  if (valid === expectValid) {
    console.log(`  ✅ ${label}  (${valid ? '通过' : '按预期拒绝'})`);
    return true;
  } else {
    console.error(`  ❌ ${label}: 期望 ${expectValid} 实际 ${valid}`);
    if (validate.errors) {
      validate.errors.forEach(e => {
        console.error(`     - ${e.instancePath || '/'} ${e.message}  ${JSON.stringify(e.params)}`);
      });
    }
    return false;
  }
}

// ── 主测试 ─────────────────────────────────────────────────────
console.log('=== 蚁小二 Schema 验证测试 (修复版) ===\n');

let passed = 0, failed = 0;

// 1. 抖音视频 - 正确（contentPublishForm 层级，不含 video 字段）
passed += testOne('douyin.video.schema.json', {
  formType: 'task',
  title: '2026年AI趋势',
  description: '深入探讨AI发展',
  tags: ['AI', '科技'],
  syncComment: false,
}, true, '抖音视频(正确)') ? 1 : 0;
failed += (passed > 0 && testOne.toString().includes('dummy')) ? 0 : 0; // placeholder

// 重新计算
// 直接手写
let p = 0, f = 0;

function T(sf, pl, ev, label) {
  const v = getValidator(sf);
  if (!v) { console.error(`  ❌ ${label}: 无法加载 schema`); f++; return; }
  const valid = v(pl);
  if (valid === ev) { console.log(`  ✅ ${label}  (${valid ? '通过' : '按预期拒绝'})`); p++; }
  else { console.error(`  ❌ ${label}: 期望 ${ev} 实际 ${valid}`); if(v.errors)v.errors.forEach(e=>console.error(`     - ${e.instancePath||'/'}: ${e.message}`)); f++; }
}

console.log('--- 视频类 ---');
// 抖音
T('douyin.video.schema.json', {formType:'task',title:'测试',description:'描述',tags:['A'],syncComment:false}, true, '抖音视频(正确)');
T('douyin.video.schema.json', {formType:'task',description:'缺title',tags:['A']}, false, '抖音视频(缺title)');
T('douyin.video.schema.json', {formType:'task',title:'测',description:'描',tags:['A'],extra:1}, false, '抖音视频(多余字段)');

// 小红书
T('xiaohongshu.video.schema.json', {formType:'task',title:'好物',description:'推荐',visibleType:0}, true, '小红书视频(正确)');
T('xiaohongshu.video.schema.json', {formType:'task',visibleType:99}, false, '小红书视频(错误visibleType)');

// 知乎视频
T('zhihu.video.schema.json', {formType:'task',description:'描述',category:[{yixiaoerId:'1',yixiaoerName:'科技',raw:{}}}, true, '知乎视频(正确)');

console.log('\n--- 文章类 ---');
// 知乎文章（需要先创建 schema，暂时跳过）
// T('zhihu.article.schema.json', {formType:'task',title:'标题',content:'<p>内容</p>'}, true, '知乎文章(正确)');

console.log(`\n=== 结果: ${p} 通过, ${f} 失败 ===`);
process.exit(f > 0 ? 1 : 0);
