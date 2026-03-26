import { uploadResource } from './upload-resource';

/**
 * 微信公众号专用发布引擎
 */
const API_KEY = process.env.YIXIAOER_API_KEY;
const API_URL = process.env.YIXIAOER_API_URL || 'https://www.yixiaoer.cn/api';

async function main() {
  console.log("WeChat Engine Initialized. Dedicated to WeChat Official Accounts.");
}

main();
