import * as fs from 'fs';
import * as path from 'path';
import { uploadResource } from './upload-resource';

const API_KEY = process.env.YIXIAOER_API_KEY;
const API_URL = process.env.YIXIAOER_API_URL || 'https://www.yixiaoer.cn/api';

async function main() {
  const args = process.argv.slice(2);
  const title = args.find(a => a.startsWith('--title='))?.split('=')[1];
  const content = args.find(a => a.startsWith('--content='))?.split('=')[1];
  const coverKeyArg = args.find(a => a.startsWith('--cover_key='))?.split('=')[1];
  const coverUrlArg = args.find(a => a.startsWith('--cover_url='))?.split('=')[1];
  const declaration = parseInt(args.find(a => a.startsWith('--declaration='))?.split('=')[1] || '2');

  if (!title || !content) {
    console.error(JSON.stringify({ error: "Missing required parameters: --title, --content" }));
    process.exit(1);
  }

  if (!API_KEY) {
     console.error(JSON.stringify({ error: "Missing YIXIAOER_API_KEY." }));
     process.exit(1);
  }

  try {
    // 1. 处理封面 (优先使用预先上传好的 Key)
    let coverKey = coverKeyArg;
    if (!coverKey && coverUrlArg) {
      // 若没有提供 Key 但有 URL，则回退到内部同步转存
      coverKey = await uploadResource(coverUrlArg, 'material-library');
    }

    // 2. 内容存证 (获取 publishContentId)
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

    // 3. 获取第一个百家号账号
    const accountsRes = await fetch(`${API_URL}/platform-accounts?platform=百家号&page=1&size=10`, {
      headers: { 'Authorization': API_KEY }
    });
    const accountsData = await accountsRes.json();
    const accounts = (accountsData.data?.data) || (accountsData.data) || accountsData;
    if (!accounts || accounts.length === 0) throw new Error('No available Baijiahao accounts found.');
    
    const accountId = accounts[0].id;

    // 4. 发起发布任务 (V2 接口)
    const taskBody = {
      desc: title,
      platforms: ['百家号'],
      publishType: 'article',
      publishChannel: 'cloud',
      coverKey: coverKey,
      publishArgs: {
        platformForms: {
          '百家号': {
            formType: 'task',
            coverType: coverKey ? 'single' : 'none',
            covers: coverKey ? [{ key: coverKey, width: 1200, height: 800, size: 0 }] : [],
            declaration: declaration
          }
        },
        accountForms: [
          {
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
              category: [{ id: '100', text: '其他', raw: null }],
              declaration: declaration,
              type: 1, 
              pubType: 1, 
              articles: [
                {
                  title: title,
                  content: wrappedContent,
                  contentHtml: wrappedContent,
                  digest: title.substring(0, 50),
                  covers: coverKey ? [{ key: coverKey, width: 1200, height: 800, size: 0 }] : [],
                  isDraft: false
                }
              ]
            }
          }
        ]
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
      error: "Failed to publish Baijiahao article", 
      details: (error as Error).message 
    }));
    process.exit(1);
  }
}

main();
