import { uploadResource } from './upload-resource.ts';

const API_KEY = process.env.YIXIAOER_API_KEY;
const API_URL = process.env.YIXIAOER_API_URL || 'https://www.yixiaoer.cn/api';

async function main() {
  const args = process.argv.slice(2);
  const title = args.find(a => a.startsWith('--title='))?.split('=')[1];
  const contentArg = args.find(a => a.startsWith('--content='))?.split('=')[1];
  const accountIdsStr = args.find(a => a.startsWith('--account_ids='))?.split('=')[1];
  const coverKeyArg = args.find(a => a.startsWith('--cover_key='))?.split('=')[1];
  const coverUrlArg = args.find(a => a.startsWith('--cover_url='))?.split('=')[1];
  const author = args.find(a => a.startsWith('--author='))?.split('=')[1] || '';
  const digest = args.find(a => a.startsWith('--digest='))?.split('=')[1] || '';
  const contentSourceUrl = args.find(a => a.startsWith('--content_source_url='))?.split('=')[1] || '';
  const isOriginal = args.find(a => a.startsWith('--original='))?.split('=')[1] === 'true';
  const notifySubscribers = args.find(a => a.startsWith('--notify='))?.split('=')[1] === 'false' ? 0 : 1;
  const pubType = parseInt(args.find(a => a.startsWith('--pub_type='))?.split('=')[1] || '1'); // 0: 草稿, 1: 公开

  if (!title || !contentArg || !accountIdsStr) {
    console.error(JSON.stringify({ error: "Missing required parameters: --title, --content, --account_ids" }));
    process.exit(1);
  }

  const targetAccountIds = accountIdsStr.split(',').map(id => id.trim());

  if (!API_KEY) {
    console.error(JSON.stringify({ error: "Missing YIXIAOER_API_KEY environment variable." }));
    process.exit(1);
  }

  try {
    const publishContentId = Array.from({ length: 24 }, () => Math.floor(Math.random() * 16).toString(16)).join('');

    let coverKey = coverKeyArg;
    if (!coverKey && coverUrlArg) {
      coverKey = await uploadResource(coverUrlArg);
    }

    const wrappedContent = `<html><body>${contentArg}</body></html>`;
    const storageRes = await fetch(`${API_URL}/storages/articles`, {
      method: 'POST',
      headers: { 'Authorization': API_KEY, 'Content-Type': 'application/json' },
      body: JSON.stringify({ publishContentId, title, content: wrappedContent, contentHtml: wrappedContent })
    });
    if (!storageRes.ok) throw new Error(`Storage failed: ${await storageRes.text()}`);

    // 适配最新 WxGongZhongHaoArticleForm DTO
    const wechatForm = {
      contentList: [{
        title,
        content: wrappedContent,
        digest: digest || title.substring(0, 120),
        cover: coverKey ? { key: coverKey, width: 1200, height: 800, size: 0 } : undefined,
        createType: isOriginal ? 1 : 0, // 0:不申明 1:申明原创
        authorName: author,
        quickRepost: 1,
        categories: [],
        contentSourceUrl,
        quickPrivateMessage: 1
      }],
      notifySubscribers: notifySubscribers,
      sex: 0,
      pubType
    };

    const taskBody = {
      desc: title,
      platforms: ['微信公众号'],
      publishType: 'weixin-gongzhonghao',
      publishChannel: 'cloud',
      isDraft: pubType === 0,
      coverKey,
      publishArgs: {
        platformForms: { '微信公众号': wechatForm },
        accountForms: targetAccountIds.map(accountId => ({
          platformAccountId: accountId,
          publishContentId,
          coverKey,
          contentPublishForm: wechatForm
        }))
      }
    };

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
