import { uploadResource } from './upload-resource.ts';

const API_KEY = process.env.YIXIAOER_API_KEY;
const API_URL = process.env.YIXIAOER_API_URL || 'https://www.yixiaoer.cn/api';

/**
 * 基础文章表单构造器 (适用于 90% 的平台)
 */
const baseArticleForm = (params: any) => ({
  title: params.title,
  covers: params.coverKey ? [{ key: params.coverKey, width: 1200, height: 800, size: 0 }] : [],
  pubType: params.pubType ?? 1,
  declaration: params.declaration ?? 0,
  articles: [{
    title: params.title,
    content: params.content,
    contentHtml: params.content,
    digest: params.digest || params.title.substring(0, 50),
    isDraft: (params.pubType ?? 1) === 0
  }]
});

/**
 * 公众号专用表单构造器 (基于 WxGongZhongHaoArticleForm DTO)
 */
const wechatArticleForm = (params: any) => ({
  contentList: [{
    title: params.title,
    content: params.content,
    digest: params.digest || params.title.substring(0, 120),
    cover: params.coverKey ? { key: params.coverKey, width: 1200, height: 800, size: 0 } : undefined,
    createType: params.original ? 1 : 0, // 0:不申明 1:申明原创
    authorName: params.author || '',
    quickRepost: 1,
    categories: [],
    contentSourceUrl: params.contentSourceUrl || '',
    quickPrivateMessage: 1
  }],
  notifySubscribers: params.notify === false ? 0 : 1, // 0:不群发 1:群发
  sex: 0,
  pubType: params.pubType ?? 1
});

/**
 * 全量文章与公众号引擎注册中心
 */
const PLATFORM_REGISTRY: Record<string, any> = {
  '微信公众号': (params: any) => wechatArticleForm(params),
  '头条号': (params: any) => baseArticleForm(params),
  '百家号': (params: any) => ({ 
    ...baseArticleForm(params), 
    coverType: 'single', 
    category: params.categoryId ? [{ 
      yixiaoerId: params.categoryId, 
      yixiaoerName: params.categoryName || '默认' 
    }] : [] 
  }),
  '企鹅号': (params: any) => ({ ...baseArticleForm(params), tags: params.tags || [] }),
  '搜狐号': (params: any) => baseArticleForm(params),
  '一点号': (params: any) => baseArticleForm(params),
  '大鱼号': (params: any) => baseArticleForm(params),
  '网易号': (params: any) => baseArticleForm(params),
  '知乎': (params: any) => ({ ...baseArticleForm(params), topics: [] }),
  '爱奇艺': (params: any) => baseArticleForm(params),
  '新浪微博': (params: any) => baseArticleForm(params),
  '哔哩哔哩': (params: any) => ({ ...baseArticleForm(params), category: params.category || [], tag: (params.tags || []).join(','), original: 1 }),
  '雪球号': (params: any) => baseArticleForm(params),
  '快传号': (params: any) => baseArticleForm(params),
  '豆瓣': (params: any) => baseArticleForm(params),
  'CSDN': (params: any) => ({ ...baseArticleForm(params), categories: [], tags: params.tags || [], type: 'original' }),
  '车家号': (params: any) => baseArticleForm(params),
  '简书': (params: any) => baseArticleForm(params),
  'WiFi万能钥匙': (params: any) => ({ ...baseArticleForm(params), category: [], desc: params.title }),
  'AcFun (A站)': (params: any) => ({ ...baseArticleForm(params), category: [], tags: params.tags || [], desc: params.title, type: 0 }),
  '易车号': (params: any) => baseArticleForm(params),
  '抖音': (params: any) => ({ ...baseArticleForm(params), description: params.title })
};

async function main() {
  const args = process.argv.slice(2);
  const title = args.find(a => a.startsWith('--title='))?.split('=')[1];
  const contentArg = args.find(a => a.startsWith('--content='))?.split('=')[1];
  const platformsStr = args.find(a => a.startsWith('--platforms='))?.split('=')[1];
  const accountIdsStr = args.find(a => a.startsWith('--account_ids='))?.split('=')[1];
  const coverKeyArg = args.find(a => a.startsWith('--cover_key='))?.split('=')[1];
  const coverUrlArg = args.find(a => a.startsWith('--cover_url='))?.split('=')[1];
  const tagsArg = args.find(a => a.startsWith('--tags='))?.split('=')[1];
  const author = args.find(a => a.startsWith('--author='))?.split('=')[1];
  const digest = args.find(a => a.startsWith('--digest='))?.split('=')[1];
  const contentSourceUrl = args.find(a => a.startsWith('--content_source_url='))?.split('=')[1];
  const original = args.find(a => a.startsWith('--original='))?.split('=')[1] === 'true';
  const notify = args.find(a => a.startsWith('--notify='))?.split('=')[1] !== 'false';
  const pubType = parseInt(args.find(a => a.startsWith('--pub_type='))?.split('=')[1] || '1');
  const declaration = parseInt(args.find(a => a.startsWith('--declaration='))?.split('=')[1] || '0');
  const categoryId = args.find(a => a.startsWith('--category_id='))?.split('=')[1];
  const categoryName = args.find(a => a.startsWith('--category_name='))?.split('=')[1];

  if (!title || !contentArg || !platformsStr) {
    console.error(JSON.stringify({ error: "Missing required parameters: --title, --content, --platforms" }));
    process.exit(1);
  }

  const platforms = platformsStr.split(',').map(p => p.trim());
  const targetAccountIds = accountIdsStr ? accountIdsStr.split(',').map(id => id.trim()) : [];
  const tags = tagsArg ? tagsArg.split(',').map(t => t.trim()) : [];

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

    const platformForms: Record<string, any> = {};
    platforms.forEach(p => {
      const factory = PLATFORM_REGISTRY[p];
      if (factory) {
        platformForms[p] = factory({ 
          title, content: wrappedContent, coverKey, tags, pubType, 
          declaration, author, digest, contentSourceUrl, original, notify,
          categoryId, categoryName
        });
      }
    });

    // 公众号特有类型识别
    const isWechatOnly = platforms.length === 1 && platforms[0] === '微信公众号';
    const publishType = isWechatOnly ? 'weixin-gongzhonghao' : 'article';

    const taskBody = {
      desc: title,
      platforms,
      publishType,
      publishChannel: 'cloud',
      isDraft: pubType === 0,
      coverKey,
      publishArgs: {
        platformForms,
        accountForms: targetAccountIds.map(accountId => ({
          platformAccountId: accountId,
          publishContentId,
          coverKey,
          contentPublishForm: platformForms[platforms[0]] || {}
        }))
      }
    };

    const publishRes = await fetch(`${API_URL}/taskSets/v2`, {
      method: 'POST',
      headers: { 'Authorization': API_KEY, 'Content-Type': 'application/json' },
      body: JSON.stringify(taskBody)
    });
    if (!publishRes.ok) throw new Error(`Publishing failed: ${await publishRes.text()}`);
    const result = await publishRes.json();

    console.log(JSON.stringify(result, null, 2));

  } catch (error) {
    console.error(JSON.stringify({ error: "Article/WeChat Engine Error", details: (error as Error).message }));
    process.exit(1);
  }
}

main();
