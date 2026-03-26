import { uploadResource } from './upload-resource';

/**
 * 图文动态统一引擎骨架
 * 待以后接入：小红书、微博、朋友圈、动态、微头条等
 */
const API_KEY = process.env.YIXIAOER_API_KEY;
const API_URL = process.env.YIXIAOER_API_URL || 'https://www.yixiaoer.cn/api';

const POST_REGISTRY: Record<string, any> = {
  // 以后在这里注册平台差异配置
};

async function main() {
  console.log("Post Engine Initialized. Waiting for platform registration...");
  // 核心逻辑参考 publish-article.ts
}

main();
