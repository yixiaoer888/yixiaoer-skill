import * as fs from 'fs';
import * as path from 'path';

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
  const title = args.find(a => a.startsWith('--title='))?.split('=')[1];
  const content = args.find(a => a.startsWith('--content='))?.split('=')[1];
  const coverKeyArg = args.find(a => a.startsWith('--cover_key='))?.split('=')[1];
  const coverUrlArg = args.find(a => a.startsWith('--cover_url='))?.split('=')[1];
  const tagsArg = args.find(a => a.startsWith('--tags='))?.split('=')[1];
  const declaration = parseInt(args.find(a => a.startsWith('--declaration='))?.split('=')[1] || '0');
  const pubType = parseInt(args.find(a => a.startsWith('--pubType='))?.split('=')[1] || '1');
  const accountIdsArg = args.find(a => a.startsWith('--account_ids='))?.split('=')[1];

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

    // 2. 内容存证
    const publishContentId = `art_${Date.now()}_${Math.floor(Math.random() * 1000)}`;
    const wrappedContent = `<html><body>${content}</body></html>`;
    const saveRes = await fetch(`${API_URL}/storages/articles`, {
      method: 'POST',
      headers: {
        'Authorization': API_KEY,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ publishContentId, content: wrappedContent })
    });
    
    if (!saveRes.ok) throw new Error(`Failed to save article content: ${await saveRes.text()}`);

    // 3. 获取企鹅号账号
    const accountsRes = await fetch(`${API_URL}/platform-accounts?platform=企鹅号&page=1&size=50`, {
      headers: { 'Authorization': API_KEY }
    });
    const accountsData = await accountsRes.json();
    const allAccounts = (accountsData.data?.data) || (accountsData.data) || accountsData || [];
    
    let targetAccountIds: string[] = [];
    if (accountIdsArg) {
      targetAccountIds = accountIdsArg.split(',');
    } else {
      if (allAccounts.length === 0) throw new Error('No available Qiehao accounts found.');
      targetAccountIds = [allAccounts[0].id];
    }

    const tags = tagsArg ? tagsArg.split(',').map(t => t.trim()) : [];

    // 4. 发起发布任务 (V2 接口)
    const taskBody = {
      desc: title,
      platforms: ['企鹅号'],
      publishType: 'article',
      publishChannel: 'cloud',
      coverKey: coverKey,
      publishArgs: {
        platformForms: {
          '企鹅号': {
            formType: 'task',
            title: title,
            tags: tags,
            declaration: declaration,
            pubType: pubType,
            covers: coverKey ? [{ key: coverKey, width: 1200, height: 800, size: 0 }] : []
          }
        },
        accountForms: targetAccountIds.map(accountId => ({
          platformAccountId: accountId,
          coverKey: coverKey,
          publishContentId: publishContentId,
          cover: coverKey ? { key: coverKey, width: 1200, height: 800, size: 0 } : undefined,
          contentPublishForm: {
            title: title,
            content: wrappedContent,
            contentHtml: wrappedContent,
            digest: title.substring(0, 50),
            covers: coverKey ? [{ key: coverKey, width: 1200, height: 800, size: 0 }] : [],
            tags: tags,
            declaration: declaration,
            pubType: pubType,
            articles: [
              {
                title: title,
                content: wrappedContent,
                contentHtml: wrappedContent,
                digest: title.substring(0, 50),
                covers: coverKey ? [{ key: coverKey, width: 1200, height: 800, size: 0 }] : [],
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
      error: "Failed to publish Qiehao article", 
      details: (error as Error).message 
    }));
    process.exit(1);
  }
}

main();
