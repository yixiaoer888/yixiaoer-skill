import { uploadResource } from './upload-resource';

const API_KEY = process.env.YIXIAOER_API_KEY;
const API_URL = process.env.YIXIAOER_API_URL || 'https://www.yixiaoer.cn/api';

/**
 * 图文分发引擎平台注册中心
 * 支持：抖音、小红书、快手、新浪微博、视频号、百家号、头条号、知乎
 */
const PLATFORM_REGISTRY: Record<string, any> = {
  '抖音': (params: any) => ({
    title: params.title || '',
    description: params.content || '',
    tagType: '位置',
    shopping_cart: [],
    visibleType: 0,
    declaration: params.declaration || 0
  }),
  '小红书': (params: any) => ({
    title: params.title || '',
    description: params.content || '',
    visibleType: 0,
    declaration: params.declaration || 0
  }),
  '快手': (params: any) => ({
    description: params.content || '',
    visibleType: 0,
    declaration: params.declaration || 0
  }),
  '新浪微博': (params: any) => ({
    description: params.content || '',
    visibleType: 0,
    declaration: params.declaration || 0
  }),
  '视频号': (params: any) => ({
    title: params.title || '',
    description: params.content || '',
    pubType: params.pub_type ?? 1
  }),
  '百家号': (params: any) => ({
    description: params.content || '',
    declaration: params.declaration || 0
  }),
  '头条号': (params: any) => ({
    description: params.content || '',
    pubType: params.pub_type ?? 1,
    declaration: params.declaration || 0
  }),
  '知乎': (params: any) => ({
    title: params.title || '',
    description: params.content || ''
  })
};

async function main() {
  const args = process.argv.slice(2);
  const title = args.find(a => a.startsWith('--title='))?.split('=')[1];
  const content = args.find(a => a.startsWith('--content='))?.split('=')[1];
  const platformsStr = args.find(a => a.startsWith('--platforms='))?.split('=')[1];
  const accountIdsStr = args.find(a => a.startsWith('--account_ids='))?.split('=')[1];
  const imageKeysArg = args.find(a => a.startsWith('--image_keys='))?.split('=')[1];
  const imageUrlsArg = args.find(a => a.startsWith('--image_urls='))?.split('=')[1];
  const pubType = parseInt(args.find(a => a.startsWith('--pub_type='))?.split('=')[1] || '1');
  const declaration = parseInt(args.find(a => a.startsWith('--declaration='))?.split('=')[1] || '0');

  if (!content || !platformsStr) {
    console.error(JSON.stringify({ error: "Missing required parameters: --content, --platforms" }));
    process.exit(1);
  }

  const platforms = platformsStr.split(',').map(p => p.trim());
  const targetAccountIds = accountIdsStr ? accountIdsStr.split(',').map(id => id.trim()) : [];

  if (!API_KEY) {
    console.error(JSON.stringify({ error: "Missing YIXIAOER_API_KEY environment variable." }));
    process.exit(1);
  }

  try {
    // 1. 处理图片素材 (支持多个)
    let imageKeys: string[] = imageKeysArg ? imageKeysArg.split(',').map(k => k.trim()) : [];
    if (imageKeys.length === 0 && imageUrlsArg) {
      const urls = imageUrlsArg.split(',').map(u => u.trim());
      for (const url of urls) {
        const key = await uploadResource(url);
        imageKeys.push(key);
      }
    }

    const imageFormItems = imageKeys.map(key => ({
      key: key,
      width: 1200,
      height: 800,
      size: 0
    }));

    // 2. 构建任务集 Body
    const platformForms: Record<string, any> = {};
    platforms.forEach(p => {
      if (PLATFORM_REGISTRY[p]) {
        platformForms[p] = PLATFORM_REGISTRY[p]({
          title, content, imageKeys, pubType, declaration
        });
      }
    });

    const taskBody = {
      desc: title || content.substring(0, 30),
      platforms: platforms,
      publishType: 'image-text',
      publishChannel: 'cloud',
      isDraft: pubType === 0,
      coverKey: imageKeys[0], // 默认取第一张作为封面
      publishArgs: {
        platformForms: platformForms,
        accountForms: targetAccountIds.map(accountId => ({
          platformAccountId: accountId,
          images: imageFormItems,
          contentPublishForm: platformForms[platforms[0]] || {}
        }))
      }
    };

    // 3. 提交发布
    const publishRes = await fetch(`${API_URL}/taskSets/v2`, {
      method: 'POST',
      headers: {
        'Authorization': API_KEY,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(taskBody)
    });

    if (!publishRes.ok) throw new Error(`Post publishing failed: ${await publishRes.text()}`);
    const result = await publishRes.json();

    console.log(JSON.stringify(result, null, 2));

  } catch (error) {
    console.error(JSON.stringify({ 
      error: "Failed to publish post batch", 
      details: (error as Error).message 
    }));
    process.exit(1);
  }
}

main();
