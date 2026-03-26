import { uploadResource } from './upload-resource.ts';

const API_KEY = process.env.YIXIAOER_API_KEY;
const API_URL = process.env.YIXIAOER_API_URL || 'https://www.yixiaoer.cn/api';

/**
 * 基础视频表单构造器
 */
const baseVideoForm = (params: any) => ({
  title: params.title || '',
  description: params.description || '',
  tags: params.tags || [],
  pubType: params.pubType ?? 1,
  declaration: params.declaration ?? 0
});

/**
 * 全量视频引擎注册中心 (30+ 平台)
 */
const PLATFORM_REGISTRY: Record<string, any> = {
  '抖音': (params: any) => ({ ...baseVideoForm(params), tagType: '位置', shopping_cart: [], visibleType: 0 }),
  '头条号': (params: any) => baseVideoForm(params),
  '哔哩哔哩': (params: any) => ({ ...baseVideoForm(params), category: params.category || [], original: 1 }),
  '哔哩哔哩-Open': (params: any) => ({ ...baseVideoForm(params), category: params.category || [] }),
  '百家号': (params: any) => baseVideoForm(params),
  '小红书': (params: any) => baseVideoForm(params),
  '快手': (params: any) => ({ ...baseVideoForm(params), visibleType: 0 }),
  '快手-Open': (params: any) => baseVideoForm(params),
  '新浪微博': (params: any) => baseVideoForm(params),
  '视频号': (params: any) => ({ ...baseVideoForm(params), category: params.category || [] }),
  '知乎': (params: any) => baseVideoForm(params),
  '企鹅号': (params: any) => baseVideoForm(params),
  '爱奇艺': (params: any) => baseVideoForm(params),
  '网易号': (params: any) => baseVideoForm(params),
  '一点号': (params: any) => baseVideoForm(params),
  '搜狐号': (params: any) => baseVideoForm(params),
  '腾讯微视': (params: any) => baseVideoForm(params),
  '搜狐视频': (params: any) => baseVideoForm(params),
  '皮皮虾': (params: any) => baseVideoForm(params),
  '腾讯视频': (params: any) => baseVideoForm(params),
  '多多视频': (params: any) => baseVideoForm(params),
  '美拍': (params: any) => baseVideoForm(params),
  'AcFun': (params: any) => baseVideoForm(params),
  '大鱼号': (params: any) => baseVideoForm(params),
  '车家号': (params: any) => baseVideoForm(params),
  '蜂网': (params: any) => baseVideoForm(params),
  '得物': (params: any) => baseVideoForm(params),
  '美柚': (params: any) => baseVideoForm(params),
  '小红书商家号': (params: any) => baseVideoForm(params),
  '易车号': (params: any) => baseVideoForm(params)
};

async function main() {
  const args = process.argv.slice(2);
  const title = args.find(a => a.startsWith('--title='))?.split('=')[1];
  const description = args.find(a => a.startsWith('--description='))?.split('=')[1];
  const platformsStr = args.find(a => a.startsWith('--platforms='))?.split('=')[1];
  const accountIdsStr = args.find(a => a.startsWith('--account_ids='))?.split('=')[1];
  const videoKeyArg = args.find(a => a.startsWith('--video_key='))?.split('=')[1];
  const videoUrlArg = args.find(a => a.startsWith('--video_url='))?.split('=')[1];
  const coverKeyArg = args.find(a => a.startsWith('--cover_key='))?.split('=')[1];
  const coverUrlArg = args.find(a => a.startsWith('--cover_url='))?.split('=')[1];
  const tagsArg = args.find(a => a.startsWith('--tags='))?.split('=')[1];
  const pubType = parseInt(args.find(a => a.startsWith('--pub_type='))?.split('=')[1] || '1');
  const declaration = parseInt(args.find(a => a.startsWith('--declaration='))?.split('=')[1] || '0');

  if (!title || (!videoKeyArg && !videoUrlArg) || !platformsStr) {
    console.error(JSON.stringify({ error: "Missing required parameters: --title, --video_key/url, --platforms" }));
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
    // 1. 处理视频资源
    let videoKey = videoKeyArg;
    if (!videoKey && videoUrlArg) {
      videoKey = await uploadResource(videoUrlArg);
    }

    // 2. 处理封面资源
    let coverKey = coverKeyArg;
    if (!coverKey && coverUrlArg) {
      coverKey = await uploadResource(coverUrlArg);
    }

    // 3. 构造任务集
    const platformForms: Record<string, any> = {};
    platforms.forEach(p => {
      const factory = PLATFORM_REGISTRY[p];
      if (factory) {
        platformForms[p] = factory({ title, description, tags, pubType, declaration });
      }
    });

    const taskBody = {
      desc: title,
      platforms,
      publishType: 'video',
      publishChannel: 'cloud',
      isDraft: pubType === 0,
      videoKey,
      coverKey,
      publishArgs: {
        platformForms,
        accountForms: targetAccountIds.map(accountId => ({
          platformAccountId: accountId,
          video: { key: videoKey, width: 1920, height: 1080, size: 0 },
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
    if (!publishRes.ok) throw new Error(`Video publishing failed: ${await publishRes.text()}`);
    const result = await publishRes.json();

    console.log(JSON.stringify(result, null, 2));

  } catch (error) {
    console.error(JSON.stringify({ error: "Video Engine Error", details: (error as Error).message }));
    process.exit(1);
  }
}

main();
