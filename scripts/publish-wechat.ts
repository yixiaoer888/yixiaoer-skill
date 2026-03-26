import { uploadResource } from './upload-resource.ts';

const API_KEY = process.env.YIXIAOER_API_KEY;
const API_URL = process.env.YIXIAOER_API_URL || 'https://www.yixiaoer.cn/api';

async function main() {
  const args = process.argv.slice(2);
  const title = args.find(a => a.startsWith('--title='))?.split('=')[1];
  const contentArg = args.find(a => a.startsWith('--content='))?.split('=')[1];
  const platformsStr = args.find(a => a.startsWith('--platforms='))?.split('=')[1] || '微信公众号';
  const accountIdsStr = args.find(a => a.startsWith('--account_ids='))?.split('=')[1];
  const coverKeyArg = args.find(a => a.startsWith('--cover_key='))?.split('=')[1];
  const coverUrlArg = args.find(a => a.startsWith('--cover_url='))?.split('=')[1];
  const author = args.find(a => a.startsWith('--author='))?.split('=')[1] || '';
  const digest = args.find(a => a.startsWith('--digest='))?.split('=')[1] || '';
  const contentSourceUrl = args.find(a => a.startsWith('--content_source_url='))?.split('=')[1] || '';
  const isOriginal = args.find(a => a.startsWith('--original='))?.split('=')[1] === 'true' ? 1 : 0;
  const notifySubscribers = args.find(a => a.startsWith('--notify='))?.split('=')[1] === 'false' ? 0 : 1;
  const pubType = parseInt(args.find(a => a.startsWith('--pub_type='))?.split('=')[1] || '1'); // 0: 草稿, 1: 公开

  if (!title || !contentArg || !accountIdsStr) {
    console.error(JSON.stringify({ error: "Missing required parameters: --title, --content, --account_ids" }));
    process.exit(1);
  }

  const platforms = platformsStr.split(',').map(p => p.trim());
  const targetAccountIds = accountIdsStr.split(',').map(id => id.trim());

  if (!API_KEY) {
    console.error(JSON.stringify({ error: "Missing YIXIAOER_API_KEY environment variable." }));
    process.exit(1);
  }

  try {
    // 0. 生成客户端内容 ID
    const publishContentId = Array.from({ length: 24 }, () => Math.floor(Math.random() * 16).toString(16)).join('');

    // 1. 处理封面
    let coverKey = coverKeyArg;
    if (!coverKey && coverUrlArg) {
      coverKey = await uploadResource(coverUrlArg);
    }

    // 2. 正文存证 (公众号也需要)
    const wrappedContent = `<html><body>${contentArg}</body></html>`;
    const storageRes = await fetch(`${API_URL}/storages/articles`, {
      method: 'POST',
      headers: { 'Authorization': API_KEY, 'Content-Type': 'application/json' },
      body: JSON.stringify({ publishContentId, title, content: wrappedContent, contentHtml: wrappedContent })
    });
    if (!storageRes.ok) throw new Error(`Storage failed: ${await storageRes.text()}`);

    // 3. 构造公众号特定表单
    const wechatForm = {
      articles: [{
        title,
        content: wrappedContent,
        authorName: author,
        cover: coverKey ? { key: coverKey, width: 1200, height: 800, size: 0 } : undefined,
        digest: digest || title.substring(0, 50),
        type: isOriginal,
        categories: [],
        quickRepost: 1,
        contentSourceUrl,
        quickPrivateMessage: 1,
        videoCardCount: 0
      }],
      notifySubscribers,
      sex: 0,
      pubType
    };

    const platformForms: Record<string, any> = {
      '微信公众号': wechatForm
    };

    const taskBody = {
      desc: title,
      platforms: ['微信公众号'],
      publishType: 'weixin-gongzhonghao',
      publishChannel: 'cloud',
      isDraft: pubType === 0,
      coverKey,
      publishArgs: {
        platformForms,
        accountForms: targetAccountIds.map(accountId => ({
          platformAccountId: accountId,
          publishContentId,
          coverKey,
          contentPublishForm: wechatForm
        }))
      }
    };

    // 4. 发起发布
    const publishRes = await fetch(`${API_URL}/taskSets/v2`, {
      method: 'POST',
      headers: { 'Authorization': API_KEY, 'Content-Type': 'application/json' },
      body: JSON.stringify(taskBody)
    });
    if (!publishRes.ok) throw new Error(`WeChat publishing failed: ${await publishRes.text()}`);
    const result = await publishRes.json();

    console.log(JSON.stringify(result, null, 2));

  } catch (error) {
    console.error(JSON.stringify({ error: "WeChat Engine Error", details: (error as Error).message }));
    process.exit(1);
  }
}

main();
