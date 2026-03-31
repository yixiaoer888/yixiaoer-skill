import * as fs from 'fs';
import * as path from 'path';

/**
 * 蚁小二 开放 API 助手
 */

export const API_KEY = process.env.YIXIAOER_API_KEY;
export const API_URL = process.env.YIXIAOER_API_URL || 'https://www.yixiaoer.cn/api';

/**
 * 解析并获取 --payload 参数中的 JSON 对象
 */
export function getPayload<T = any>(): T {
  const args = process.argv.slice(2);
  const payloadArg = args.find((a: string) => a.startsWith('--payload='))?.split('=')[1];

  if (!payloadArg) {
    throw new Error("Missing required parameter: --payload");
  }

  try {
    return JSON.parse(payloadArg) as T;
  } catch (error) {
    throw new Error(`Invalid JSON in --payload: ${(error as Error).message}`);
  }
}

/**
 * 低级 API 调用封装
 */
export async function callApi(endpoint: string, options: RequestInit = {}) {
  if (!API_KEY) {
    throw new Error("Missing YIXIAOER_API_KEY environment variable");
  }

  const url = endpoint.startsWith('http') ? endpoint : `${API_URL}${endpoint.startsWith('/') ? '' : '/'}${endpoint}`;
  
  const headers: HeadersInit = {
    'Authorization': API_KEY,
    'Content-Type': 'application/json',
    ...(options.headers || {})
  };

  const response = await fetch(url, {
    ...options,
    headers
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(`HTTP ${response.status}: ${errorText}`);
  }

  const result = await response.json();
  return result;
}

/**
 * 助手函数: 通用资源上传
 */
export async function uploadResource(urlOrPath: string, bucket: string = 'cloud-publish'): Promise<string> {
  let buffer: ArrayBuffer;
  let fileName = 'file.jpg';

  // 1. 获取资源内容
  if (urlOrPath.startsWith('http')) {
    const res = await fetch(urlOrPath);
    if (!res.ok) throw new Error(`HTTP error downloading file during sync upload: ${res.status}`);
    buffer = await res.arrayBuffer();
    const urlObj = new URL(urlOrPath);
    fileName = urlObj.pathname.split('/').pop() || 'image.jpg';
    if (!fileName.includes('.')) fileName += '.jpg';
  } else {
    const absolutePath = path.isAbsolute(urlOrPath) ? urlOrPath : path.resolve(process.cwd(), urlOrPath);
    if (!fs.existsSync(absolutePath)) {
      throw new Error(`Local file not found: ${absolutePath}`);
    }
    const fileBuffer = fs.readFileSync(absolutePath);
    buffer = fileBuffer.buffer.slice(fileBuffer.byteOffset, fileBuffer.byteOffset + fileBuffer.byteLength);
    fileName = path.basename(absolutePath);
  }

  // 2. 获取预签名上传地址 (使用 callApi)
  const uploadInfo = await callApi(`/storages/${bucket}/upload-url?fileKey=${encodeURIComponent(fileName)}`);
  const data = uploadInfo.data || uploadInfo;
  const { serviceUrl, key } = data;

  if (!serviceUrl) {
    throw new Error(`Invalid upload info response: ${JSON.stringify(uploadInfo)}`);
  }

  // 3. 执行 PUT 上传
  const putRes = await fetch(serviceUrl, {
    method: 'PUT',
    body: buffer,
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' }
  });

  if (!putRes.ok) {
    throw new Error(`Failed to upload to OSS: ${await putRes.text()}`);
  }

  return key;
}

/**
 * 统一错误处理并输出到标准输出
 */
export function handleError(error: any, context: string) {
  console.error(JSON.stringify({ 
    error: `Failed to ${context}`, 
    details: error instanceof Error ? error.message : String(error) 
  }, null, 2));
  process.exit(1);
}
