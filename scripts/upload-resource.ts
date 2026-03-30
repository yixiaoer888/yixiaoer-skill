import * as fs from 'fs';
import * as path from 'path';

/**
 * 通用资源上传辅助函数
 * 将本地文件或远程 URL 上传到蚁小二 OSS
 */
export async function uploadResource(urlOrPath: string, bucket: string = 'cloud-publish') {
  const API_KEY = process.env.YIXIAOER_API_KEY;
  const API_URL = process.env.YIXIAOER_API_URL || 'https://www.yixiaoer.cn/api';

  if (!API_KEY) {
    throw new Error("Missing YIXIAOER_API_KEY environment variable.");
  }

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

  // 2. 调用 /api/storages/[bucket]/upload-url 获取预签名上传地址
  const url = `${API_URL}/storages/${bucket}/upload-url?fileKey=${encodeURIComponent(fileName)}`;
  const uploadInfoRes = await fetch(url, {
    headers: { 'Authorization': API_KEY }
  });
  
  if (!uploadInfoRes.ok) {
    throw new Error(`Failed to get upload info: ${await uploadInfoRes.text()}`);
  }

  const uploadInfo = await uploadInfoRes.json();
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
 * 主执行入口
 * 仅支持通过 --payload 传入参数。
 * 调用方式: node scripts/upload-resource.ts --payload='{"url":"https://...","bucket":"cloud-publish"}'
 */
async function main() {
  const args = process.argv.slice(2);
  const payloadArg = args.find(a => a.startsWith('--payload='))?.split('=')[1];

  if (!payloadArg) {
    console.error(JSON.stringify({ error: "Missing required parameter: --payload" }));
    process.exit(1);
  }

  try {
    const payload = JSON.parse(payloadArg);
    const urlArg = payload.url;
    const bucket = payload.bucket || 'cloud-publish';

    if (!urlArg) {
      throw new Error("Missing required field: url in payload");
    }

    const key = await uploadResource(urlArg, bucket);
    console.log(JSON.stringify({ 
        key, 
        name: urlArg.startsWith('http') ? new URL(urlArg).pathname.split('/').pop() : urlArg.split('/').pop() 
    }, null, 2));

  } catch (error) {
    console.error(JSON.stringify({ 
      error: "Failed to upload resource", 
      details: (error as Error).message 
    }));
    process.exit(1);
  }
}

// 简单的判断是否为直接执行
if (process.argv[1]?.replace(/\\/g, '/').endsWith('scripts/upload-resource.ts') || 
    process.argv[1]?.replace(/\\/g, '/').endsWith('scripts/upload-resource')) {
  main();
}

export {};
