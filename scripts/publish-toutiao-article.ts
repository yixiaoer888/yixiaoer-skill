import * as fs from 'fs';
import * as path from 'path';

export {};

const API_KEY = process.env.YIXIAOER_API_KEY;
const API_URL = process.env.YIXIAOER_API_URL || 'https://www.yixiaoer.cn/api';

async function uploadResource(urlOrPath: string, bucket: string = 'cloud-publish') {
  let buffer: ArrayBuffer;
  let fileName = 'file.jpg';

  if (urlOrPath.startsWith('http')) {
    const res = await fetch(urlOrPath);
    if (!res.ok) throw new Error(`HTTP error downloading file during sync upload: ${res.status}`);
    buffer = await res.arrayBuffer();
    const urlObj = new URL(urlOrPath);
    fileName = urlObj.pathname.split('/').pop() || 'image.jpg';
    if (!fileName.includes('.')) fileName += '.jpg';
  } else {
    // 读取本地文件流
    const absolutePath = path.isAbsolute(urlOrPath) ? urlOrPath : path.resolve(process.cwd(), urlOrPath);
    if (!fs.existsSync(absolutePath)) {
        throw new Error(`Local file not found: ${absolutePath}`);
    }
    const fileBuffer = fs.readFileSync(absolutePath);
    buffer = fileBuffer.buffer.slice(fileBuffer.byteOffset, fileBuffer.byteOffset + fileBuffer.byteLength);
    fileName = path.basename(absolutePath);
  }

  const uploadInfoRes = await fetch(`${API_URL}/storages/${bucket}/upload-url?fileKey=${fileName}`, {
    headers: { 'Authorization': API_KEY! }
  });
  const uploadInfo = await uploadInfoRes.json();
  const { serviceUrl, key } = uploadInfo.data || uploadInfo;

  await fetch(serviceUrl, {
    method: 'PUT',
    body: buffer,
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' }
  });

  return key;
}

async function main() {
  const args = process.argv.slice(2);
  const title = args.find((a: string) => a.startsWith('--title='))?.split('=')[1];
  const content = args.find((a: string) => a.startsWith('--content='))?.split('=')[1];
  const coverKeyArg = args.find((a: string) => a.startsWith('--cover_key='))?.split('=')[1];
  const coverUrlArg = args.find((a: string) => a.startsWith('--cover_url='))?.split('=')[1];
  
  const isFirst = args.find((a: string) => a.startsWith('--is_first='))?.split('=')[1] === 'true';
  const advertisement = parseInt(args.find((a: string) => a.startsWith('--advertisement='))?.split('=')[1] || '2');
  const declaration = parseInt(args.find((a: string) => a.startsWith('--declaration='))?.split('=')[1] || '0');
  const pubType = parseInt(args.find((a: string) => a.startsWith('--pub_type='))?.split('=')[1] || '1');
  const accountIdsArg = args.find((a: string) => a.startsWith('--account_ids='))?.split('=')[1];

  if (!title || !content) {
    console.error(JSON.stringify({ error: "Missing required parameters: --title, --content" }));
    process.exit(1);
  }

  if (!API_KEY) {
     console.error(JSON.stringify({ error: "Missing YIXIAOER_API_KEY." }));
     process.exit(1);
  }

  try {
    // 1. 处理封面
    let coverKey = coverKeyArg;
    if (!coverKey && coverUrlArg) {
      coverKey = await uploadResource(coverUrlArg, 'material-library');
    }

    // 2. 内容存证 (获取 publishContentId)
    const publishContentId = `art_toutiao_${Date.now()}_${Math.floor(Math.random() * 1000)}`;
    const wrappedContent = content.includes('<html>') ? content : `<html><body>${content}</body></html>`;
    const saveRes = await fetch(`${API_URL}/storages/articles`, {
      method: 'POST',
      headers: {
        'Authorization': API_KEY,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ publishContentId, content: wrappedContent })
    });
    
    if (!saveRes.ok) throw new Error(`Failed to save article content: ${await saveRes.text()}`);

    // 3. 获取账号
    let targetAccountIds = accountIdsArg ? accountIdsArg.split(',') : [];
    
    if (targetAccountIds.length === 0) {
      const accountsRes = await fetch(`${API_URL}/platform-accounts?platform=头条号&page=1&size=10`, {
        headers: { 'Authorization': API_KEY }
      });
      const accountsData = await accountsRes.json();
      const accounts = accountsData.data?.data || accountsData.data || accountsData;
      if (!accounts || accounts.length === 0) throw new Error('No available Toutiao accounts found.');
      targetAccountIds = [accounts[0].id];
    }

    // 4. 发起发布任务 (V2 接口)
    const covers = coverKey ? [{ key: coverKey, width: 1200, height: 800, size: 0 }] : [];
    
    const taskBody = {
      desc: title,
      platforms: ['头条号'],
      publishType: 'article',
      publishChannel: 'cloud',
      coverKey: coverKey,
      publishArgs: {
        platformForms: {
          '头条号': {
            formType: 'task',
            title: title,
            isFirst: isFirst,
            advertisement: advertisement,
            declaration: declaration,
            pubType: pubType,
            covers: covers
          }
        },
        accountForms: targetAccountIds.map((id: string) => ({
          platformAccountId: id,
          coverKey: coverKey,
          publishContentId: publishContentId,
          contentPublishForm: {
            title: title,
            content: wrappedContent,
            contentHtml: wrappedContent,
            digest: title.substring(0, 50),
            covers: covers,
            isFirst: isFirst,
            advertisement: advertisement,
            declaration: declaration,
            pubType: pubType,
            articles: [
              {
                title: title,
                content: wrappedContent,
                contentHtml: wrappedContent,
                digest: title.substring(0, 50),
                covers: covers,
                isDraft: pubType === 0
              }
            ]
          }
        }))
      }
    };

    const pubRes = await fetch(`${API_URL}/taskSets/v2`, {
      method: 'POST',
      headers: {
        'Authorization': API_KEY,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(taskBody)
    });

    const pubResult = await pubRes.json();
    console.log(JSON.stringify(pubResult.data || pubResult, null, 2));

  } catch (error) {
    console.error(JSON.stringify({ 
      error: "Failed to publish Toutiao article", 
      details: (error as Error).message 
    }));
    process.exit(1);
  }
}

main();
