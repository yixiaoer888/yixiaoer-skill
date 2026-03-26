import { uploadResource } from './upload-resource';

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
    digest: params.title.substring(0, 50),
    isDraft: (params.pubType ?? 1) === 0
  }]
});

/**
 * 全量文章引擎注册中心 (21+ 平台)
 */
const PLATFORM_REGISTRY: Record<string, any> = {
  '头条号': (params: any) => baseArticleForm(params),
  '百家号': (params: any) => {
    const form = baseArticleForm(params);
    return { ...form, coverType: 'single' };
  },
  '企鹅号': (params: any) => {
    const form = baseArticleForm(params);
    return { ...form, tags: params.tags || [] };
  },
  '搜狐号': (params: any) => baseArticleForm(params),
  '一点号': (params: any) => baseArticleForm(params),
  '大鱼号': (params: any) => baseArticleForm(params),
  '网易号': (params: any) => baseArticleForm(params),
  '知乎': (params: any) => {
    const form = baseArticleForm(params);
    return { ...form, topics: [] };
  },
  '爱奇艺': (params: any) => baseArticleForm(params),
  '新浪微博': (params: any) => baseArticleForm(params),
  '哔哩哔哩': (params: any) => {
    const form = baseArticleForm(params);
    return { ...form, category: params.category || [], tag: (params.tags || []).join(','), original: 1 };
  },
  '雪球号': (params: any) => baseArticleForm(params),
  '快传号': (params: any) => baseArticleForm(params),
  '豆瓣': (params: any) => baseArticleForm(params),
  'CSDN': (params: any) => {
    const form = baseArticleForm(params);
    return { ...form, categories: [], tags: params.tags || [], type: 'original' };
  },
  '车家号': (params: any) => baseArticleForm(params),
  '简书': (params: any) => baseArticleForm(params),
  'WiFi万能钥匙': (params: any) => {
    const form = baseArticleForm(params);
    return { ...form, category: [], desc: params.title };
  },
  'AcFun (A站)': (params: any) => {
    const form = baseArticleForm(params);
    return { ...form, category: [], tags: params.tags || [], desc: params.title, type: 0 };
  },
  '易车号': (params: any) => baseArticleForm(params),
  '抖音': (params: any) => {
    const form = baseArticleForm(params);
    return { ...form, description: params.title };
  }
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
  const pubType = parseInt(args.find(a => a.startsWith('--pub_type='))?.split('=')[1] || '1');
  const declaration = parseInt(args.find(a => a.startsWith('--declaration='))?.split('=')[1] || '0');

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
    // 1. 处理封面
    let coverKey = coverKeyArg;
    if (!coverKey && coverUrlArg) {
      coverKey = await uploadResource(coverUrlArg);
    }

    // 2. 正文存证
    const wrappedContent = `<html><body>${contentArg}</body></html>`;
    const storageRes = await fetch(`${API_URL}/storages/articles`, {
      method: 'POST',
      headers: { 'Authorization': API_KEY, 'Content-Type': 'application/json' },
      body: JSON.stringify({ title, content: wrappedContent, contentHtml: wrappedContent })
    });
    if (!storageRes.ok) throw new Error(`Storage failed: ${await storageRes.text()}`);
    const { data: { publishContentId } } = await storageRes.json();

    // 3. 构造任务集
    const platformForms: Record<string, any> = {};
    platforms.forEach(p => {
      const factory = PLATFORM_REGISTRY[p];
      if (factory) {
        platformForms[p] = factory({ title, content: wrappedContent, coverKey, tags, pubType, declaration });
      }
    });

    const taskBody = {
      desc: title,
      platforms,
      publishType: 'article',
      publishChannel: 'cloud',
      isDraft: pubType === 0,
      coverKey,
      publishArgs: {
        platformForms,
        accountForms: targetAccountIds.map(accountId => ({
          platformAccountId: accountId,
          publishContentId,
          coverKey,
          cover: coverKey ? { key: coverKey, width: 1200, height: 800, size: 0 } : undefined,
          contentPublishForm: platformForms[platforms[0]] || {}
        }))
      }
    };

    // 4. 发起发布
    const publishRes = await fetch(`${API_URL}/taskSets/v2`, {
      method: 'POST',
      headers: { 'Authorization': API_KEY, 'Content-Type': 'application/json' },
      body: JSON.stringify(taskBody)
    });
    if (!publishRes.ok) throw new Error(`Publishing failed: ${await publishRes.text()}`);
    const result = await publishRes.json();

    console.log(JSON.stringify(result, null, 2));

  } catch (error) {
    console.error(JSON.stringify({ error: "Article Engine Error", details: (error as Error).message }));
    process.exit(1);
  }
}

main();
