/**
 * 最简验证测试 - 直接调用 SchemaValidator，不依赖 api.ts
 */

import * as fs from 'fs';
import * as path from 'path';
import { SchemaValidator } from '../../scripts/validator';

const validator = new SchemaValidator();

function readJsonFile(filePath: string) {
  const content = fs.readFileSync(filePath, 'utf-8').replace(/^\uFEFF/, '');
  return JSON.parse(content);
}

const testCases = [
  // [platform, type, payloadFile, expectValid]
  ['douyin', 'video',     '../fixtures/payloads/douyin-video-valid.json',   true],
  ['douyin', 'video',     '../fixtures/payloads/douyin-video-invalid.json', false],
  ['douyin', 'video',     '../fixtures/payloads/douyin-video-url-invalid.json', false],
  ['douyin', 'video',     '../fixtures/payloads/douyin-video-location-missing-raw.json', false],
  ['Xiaohongshu', 'video', '../fixtures/payloads/xiaohongshu-video-valid.json', true],
  ['zhihu', 'article',   '../fixtures/payloads/zhihu-article-valid.json',  true],
  ['douyin', 'article',  '../fixtures/payloads/douyin-article-blank-content.json', false],
];

function localPreflight(type: string, data: any): string[] {
  const errors: string[] = [];
  const walk = (value: any, visit: (value: any, path: string) => void, currentPath = '$') => {
    visit(value, currentPath);
    if (Array.isArray(value)) {
      value.forEach((item, index) => walk(item, visit, `${currentPath}[${index}]`));
      return;
    }
    if (value && typeof value === 'object') {
      Object.entries(value).forEach(([key, child]) => walk(child, visit, `${currentPath}.${key}`));
    }
  };
  const isUrl = (value: any) => typeof value === 'string' && /^https?:\/\//i.test(value);
  const requireResource = (resource: any, label: string, fields: string[]) => {
    if (!resource || typeof resource !== 'object') {
      errors.push(`${label}: missing resource`);
      return;
    }
    for (const field of ['key', ...fields]) {
      if (resource[field] === undefined || resource[field] === null || resource[field] === '') {
        errors.push(`${label}: missing ${field}`);
      }
    }
    walk(resource, value => {
      if (isUrl(value)) errors.push(`${label}: contains external URL`);
    });
  };
  const requireRaw = (value: any, label: string) => {
    if (Array.isArray(value)) {
      value.forEach((item, index) => requireRaw(item, `${label}[${index}]`));
      return;
    }
    if (!value || typeof value !== 'object') return;
    const hasIdentity = value.yixiaoerId !== undefined || value.yixiaoerName !== undefined || value.id !== undefined || value.name !== undefined;
    if (hasIdentity && (!value.raw || typeof value.raw !== 'object')) {
      errors.push(`${label}: missing raw`);
    }
  };

  if (!data.accountForms) return errors;
  data.accountForms.forEach((form: any, index: number) => {
    const cpf = form.contentPublishForm || {};
    if (type === 'video') {
      requireResource(form.video || cpf.video, `accountForms[${index}].video`, ['size', 'width', 'height']);
    }
    if (type === 'image-text') {
      if (!Array.isArray(form.images) || form.images.length === 0) {
        errors.push(`accountForms[${index}].images: required`);
      } else {
        form.images.forEach((image: any, imageIndex: number) => {
          requireResource(image, `accountForms[${index}].images[${imageIndex}]`, ['size', 'width', 'height']);
        });
      }
    }
    if (type === 'article' && (typeof cpf.content !== 'string' || cpf.content.trim().length === 0)) {
      errors.push(`accountForms[${index}].contentPublishForm.content: required and must not be blank`);
    }
    walk(form, value => {
      if (isUrl(value)) errors.push(`accountForms[${index}]: contains external URL`);
    });
    ['location', 'music', 'collection', 'collections', 'challenge', 'challenges', 'goods', 'group', 'miniapp'].forEach(field => {
      if (form[field] !== undefined) requireRaw(form[field], `accountForms[${index}].${field}`);
      if (cpf[field] !== undefined) requireRaw(cpf[field], `accountForms[${index}].contentPublishForm.${field}`);
    });
  });
  return errors;
}

async function runTests() {
  let passed = 0;
  let failed = 0;

  for (const [platform, type, file, expectValid] of testCases) {
    const absPath = path.resolve(__dirname, file);
    if (!fs.existsSync(absPath)) {
      console.error(`❌ 文件不存在: ${absPath}`);
      failed++;
      continue;
    }

    const data = readJsonFile(absPath);
    const { valid, errors } = await validator.validate(platform, type, data);
    const preflightErrors = localPreflight(type as string, data);
    const finalValid = valid && preflightErrors.length === 0;

    const platformType = `${platform}(${type})`;
    if (finalValid === expectValid) {
      console.log(`✅ ${platformType}: ${finalValid ? '通过' : '按预期拒绝'}`);
      passed++;
    } else {
      console.error(`❌ ${platformType}: 期望 ${expectValid} 但得到 ${finalValid}`);
      if (errors) errors.forEach((e: string) => console.error(`   - ${e}`));
      preflightErrors.forEach((e: string) => console.error(`   - ${e}`));
      failed++;
    }
  }

  console.log(`\n结果: ${passed} 通过, ${failed} 失败`);
  process.exit(failed > 0 ? 1 : 0);
}

runTests().catch(err => {
  console.error('测试运行失败:', err);
  process.exit(1);
});
