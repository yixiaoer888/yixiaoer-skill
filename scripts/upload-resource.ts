import * as fs from 'fs';
import * as path from 'path';

const API_KEY = process.env.YIXIAOER_API_KEY;
const API_URL = process.env.YIXIAOER_API_URL || 'https://www.yixiaoer.cn/api';

async function main() {
  const args = process.argv.slice(2);
  const urlArg = args.find(a => a.startsWith('--url='))?.split('=')[1];
  const bucket = args.find(a => a.startsWith('--bucket='))?.split('=')[1] || 'cloud-publish';

  if (!urlArg) {
    console.error(JSON.stringify({ error: "Missing parameter: --url" }));
    process.exit(1);
  }

  if (!API_KEY) {
     console.error(JSON.stringify({ error: "Missing YIXIAOER_API_KEY." }));
     process.exit(1);
  }

  try {
    // 1. 获取资源内容
    let buffer: ArrayBuffer;
    let fileName = 'file.jpg';

    if (urlArg.startsWith('http')) {
        const res = await fetch(urlArg);
        if (!res.ok) throw new Error(`HTTP error downloading file: ${res.status}`);
        buffer = await res.arrayBuffer();
        const urlObj = new URL(urlArg);
        fileName = urlObj.pathname.split('/').pop() || 'image.jpg';
        if (!fileName.includes('.')) fileName += '.jpg';
    } else {
        // 读取本地文件流
        const absolutePath = path.isAbsolute(urlArg) ? urlArg : path.resolve(process.cwd(), urlArg);
        if (!fs.existsSync(absolutePath)) {
            throw new Error(`Local file not found: ${absolutePath}`);
        }
        const fileBuffer = fs.readFileSync(absolutePath);
        buffer = fileBuffer.buffer.slice(fileBuffer.byteOffset, fileBuffer.byteOffset + fileBuffer.byteLength);
        fileName = path.basename(absolutePath);
    }

    // 2. 调用 /api/storages/[bucket]/upload-url 获取预签名上传地址
    const url = new URL(`${API_URL}/storages/${bucket}/upload-url`);
    url.searchParams.append('fileKey', fileName);

    const uploadUrlRes = await fetch(url.toString(), {
      headers: { 'Authorization': API_KEY }
    });
    
    if (!uploadUrlRes.ok) throw new Error(`Failed to get upload URL: ${await uploadUrlRes.text()}`);
    const { data: { serviceUrl, key } } = await uploadUrlRes.json();

    // 3. 执行 PUT 上传
    const putRes = await fetch(serviceUrl, {
      method: 'PUT',
      body: buffer,
      headers: { 
        'Content-Type': 'application/x-www-form-urlencoded' // 模拟插件中的习惯
      }
    });

    if (!putRes.ok) throw new Error(`Failed to upload to OSS: ${await putRes.text()}`);

    // 4. 返回 Key
    console.log(JSON.stringify({ key, name: fileName }, null, 2));

  } catch (error) {
    console.error(JSON.stringify({ 
      error: "Failed to upload resource", 
      details: (error as Error).message 
    }));
    process.exit(1);
  }
}

main();
