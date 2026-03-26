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
    // 1. 处理封面
    let coverUrl = coverUrlArg;
    let coverKey = coverKeyArg;
    
    // 如果只有 URL 没有 Key，先上传获取 Key
    if (!coverKey && coverUrlArg) {
      coverKey = await uploadResource(coverUrlArg);
    }
    // 如果只有 Key 没有 URL，获取访问地址 (公众号发布需要 URL 作为素材源)
    if (!coverUrl && coverKey) {
      coverUrl = `${API_URL}/storages/assets/access-url?fileKey=${coverKey}`;
    }

    // 2. 构造公众号特定文章结构
    const article = {
      title,
      author: author,
      digest: digest || title.substring(0, 50),
      content: contentArg, // 后端会自动处理标签替换
      contentSourceUrl,
      thumbUrl: coverUrl || '', // 这里的 thumbUrl 通常是图片的访问地址
      needOpenComment: 1,
      onlyFansCanComment: 0
    };

    const taskBody = {
      platformAccountIds: targetAccountIds,
      publishType: 'article',
      sendIgnoreReprint: 0,
      sendAll: notifySubscribers ? 1 : 0,
      coverKey: coverKey || '',
      desc: title,
      articles: [article]
    };

    // 3. 发起发布
    const publishRes = await fetch(`${API_URL}/wx/publish`, {
      method: 'POST',
      headers: { 
        'Authorization': API_KEY, 
        'Content-Type': 'application/json' 
      },
      body: JSON.stringify(taskBody)
    });
    
    if (!publishRes.ok) {
      const errorText = await publishRes.text();
      throw new Error(`WeChat publishing failed (Terminal Status ${publishRes.status}): ${errorText}`);
    }
    
    const result = await publishRes.json();
    console.log(JSON.stringify(result, null, 2));

  } catch (error) {
    console.error(JSON.stringify({ error: "WeChat Engine Error", details: (error as Error).message }));
    process.exit(1);
  }
}

main();
