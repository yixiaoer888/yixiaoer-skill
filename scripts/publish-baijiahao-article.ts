/**
 * Publish Baijiahao Article (publish-baijiahao-article.ts)
 * 
 * 实现该脚本时使用的底层接口依赖：
 * 1. GET /api/platform-accounts
 * 2. GET /api/storages/[bucket]/upload-url
 * 3. POST /api/storages/articles (内容存证)
 * 4. POST /api/taskSets/v2 (任务派发)
 * 
 * 免外部包依赖运行环境：Node.js 18+
 */

const API_KEY = process.env.YIXIAOER_API_KEY;
const API_URL = process.env.YIXIAOER_API_URL || 'https://www.yixiaoer.cn/api';

async function uploadResource(urlOrPath: string, bucket: string = 'cloud-publish') {
  // 1. 获取文件的内容 (Buffer/ArrayBuffer)
  let buffer: ArrayBuffer;
  let fileName = 'file.jpg';

  if (urlOrPath.startsWith('http')) {
    const res = await fetch(urlOrPath);
    buffer = await res.arrayBuffer();
    const urlObj = new URL(urlOrPath);
    fileName = urlObj.pathname.split('/').pop() || 'image.jpg';
    if (!fileName.includes('.')) fileName += '.jpg';
  } else {
    // 简化处理：假设脚本主要处理 URL 形式的公开资源进行转存
    throw new Error('Local file path is not supported in this lightweight script yet.');
  }

  // 2. 获取上传预签名 URL
  const uploadInfoRes = await fetch(`${API_URL}/storages/${bucket}/upload-url?fileKey=${fileName}`, {
    headers: { 'Authorization': API_KEY! }
  });
  const uploadInfo = await uploadInfoRes.json();
  const { serviceUrl, key } = uploadInfo.data || uploadInfo;

  // 3. 执行上传 (PUT)
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
  const coverUrl = args.find(a => a.startsWith('--cover_url='))?.split('=')[1];
  const declaration = parseInt(args.find(a => a.startsWith('--declaration='))?.split('=')[1] || '2'); // 2为原创

  if (!title || !content) {
    console.error(JSON.stringify({ error: "Missing required parameters: --title, --content" }));
    process.exit(1);
  }

  if (!API_KEY) {
     console.error(JSON.stringify({ error: "Missing YIXIAOER_API_KEY. Please check your .env or environment settings." }));
     process.exit(1);
  }

  try {
    // 1. 处理封面上传（可选）
    const coverKey = coverUrl ? await uploadResource(coverUrl, 'material-library') : undefined;

    // 2. 文章内容存证获取 publishContentId
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

    // 3. 获取第一个可用的百家号账号
    const accountsRes = await fetch(`${API_URL}/platform-accounts?platform=百家号&page=1&size=10`, {
      headers: { 'Authorization': API_KEY }
    });
    const { data: { data: accounts } } = await accountsRes.json();
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
              type: 1, // 图文类型
              pubType: 1, // 公开发布
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
