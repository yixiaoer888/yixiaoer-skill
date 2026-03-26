import { uploadResource } from './upload-resource';

const API_KEY = process.env.YIXIAOER_API_KEY;
const API_URL = process.env.YIXIAOER_API_URL || 'https://www.yixiaoer.cn/api';

/**
 * 平台配置注册中心
 * 为 15+ 平台扩展预留基础结构
 */
const PLATFORM_REGISTRY: Record<string, any> = {
  '百家号': (params: any) => ({
    title: params.title,
    coverType: 'single',
    covers: params.coverKey ? [{ key: params.coverKey, width: 1200, height: 800, size: 0 }] : [],
    category: [],
    declaration: params.declaration || 0,
    pubType: params.pubType ?? 1,
    articles: [{
      title: params.title,
      content: params.content,
      contentHtml: params.content,
      digest: params.title.substring(0, 50),
      isDraft: (params.pubType ?? 1) === 0
    }]
  }),
  '企鹅号': (params: any) => ({
    title: params.title,
    covers: params.coverKey ? [{ key: params.coverKey, width: 1200, height: 800, size: 0 }] : [],
    tags: params.tags || [],
    declaration: params.qiehao_declaration || params.declaration || 0,
    pubType: params.pubType ?? 1,
    articles: [{
      title: params.title,
      content: params.content,
      contentHtml: params.content,
      digest: params.title.substring(0, 50),
      isDraft: (params.pubType ?? 1) === 0
    }]
  }),
  '头条号': (params: any) => ({
    title: params.title,
    coverType: 'single',
    covers: params.coverKey ? [{ key: params.coverKey, width: 1200, height: 800, size: 0 }] : [],
    declaration: params.declaration || 0,
    isFirst: params.toutiao_is_first || false,
    advertisement: 0,
    pubType: params.pubType ?? 1,
    articles: [{
      title: params.title,
      content: params.content,
      contentHtml: params.content,
      digest: params.title.substring(0, 50),
      isDraft: (params.pubType ?? 1) === 0
    }]
  })
};

async function main() {
  const args = process.argv.slice(2);
  const title = args.find(a => a.startsWith('--title='))?.split('=')[1];
  const content = args.find(a => a.startsWith('--content='))?.split('=')[1];
  const platformsStr = args.find(a => a.startsWith('--platforms='))?.split('=')[1];
  const accountIdsStr = args.find(a => a.startsWith('--account_ids='))?.split('=')[1];
  const coverKeyArg = args.find(a => a.startsWith('--cover_key='))?.split('=')[1];
  const coverUrlArg = args.find(a => a.startsWith('--cover_url='))?.split('=')[1];
  const tagsArg = args.find(a => a.startsWith('--tags='))?.split('=')[1];
  const pubType = parseInt(args.find(a => a.startsWith('--pub_type='))?.split('=')[1] || '1');
  const declaration = parseInt(args.find(a => a.startsWith('--declaration='))?.split('=')[1] || '0');
  
  // 平台专有参数解析
  const toutiao_is_first = args.find(a => a.startsWith('--toutiao_is_first='))?.split('=')[1] === 'true';
  const qiehao_declaration = parseInt(args.find(a => a.startsWith('--qiehao_declaration='))?.split('=')[1] || '0');

  if (!title || !content || !platformsStr) {
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
    // 1. 处理封面资源
    let coverKey = coverKeyArg;
    if (!coverKey && coverUrlArg) {
      coverKey = await uploadResource(coverUrlArg);
    }

    // 2. 将内容云端存证 (Article Storage)
    const wrappedContent = `<html><body>${content}</body></html>`;
    const storageRes = await fetch(`${API_URL}/storages/articles`, {
      method: 'POST',
      headers: { 
        'Authorization': API_KEY,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        title: title,
        content: wrappedContent,
        contentHtml: wrappedContent
      })
    });

    if (!storageRes.ok) throw new Error(`Content storage failed: ${await storageRes.text()}`);
    const { data: { publishContentId } } = await storageRes.json();

    // 3. 构建统一的任务集 Body
    const platformForms: Record<string, any> = {};
    platforms.forEach(p => {
      if (PLATFORM_REGISTRY[p]) {
        platformForms[p] = PLATFORM_REGISTRY[p]({
          title, content: wrappedContent, coverKey, tags, pubType, declaration,
          toutiao_is_first, qiehao_declaration
        });
      }
    });

    const taskBody = {
      desc: title,
      platforms: platforms,
      publishType: 'article',
      publishChannel: 'cloud',
      isDraft: pubType === 0,
      coverKey: coverKey,
      publishArgs: {
        platformForms: platformForms,
        accountForms: targetAccountIds.map(accountId => {
          // 查找该账号所属平台（基于已有参数推断，或 API 层面自动匹配）
          // 这里简化处理：所有账号共享该内容的配置
          return {
            platformAccountId: accountId,
            publishContentId: publishContentId,
            coverKey: coverKey,
            cover: coverKey ? { key: coverKey, width: 1200, height: 800, size: 0 } : undefined,
            contentPublishForm: platformForms[platforms[0]] || {} // 兜底使用首个平台表单
          };
        })
      }
    };

    // 4. 原子提交任务
    const publishRes = await fetch(`${API_URL}/taskSets/v2`, {
      method: 'POST',
      headers: {
        'Authorization': API_KEY,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(taskBody)
    });

    if (!publishRes.ok) throw new Error(`Batch publishing failed: ${await publishRes.text()}`);
    const result = await publishRes.json();

    console.log(JSON.stringify(result, null, 2));

  } catch (error) {
    console.error(JSON.stringify({ 
      error: "Failed to publish article batch", 
      details: (error as Error).message 
    }));
    process.exit(1);
  }
}

main();
