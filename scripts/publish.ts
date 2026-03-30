import { uploadResource } from './upload-resource.ts';

const API_KEY = process.env.YIXIAOER_API_KEY;
const API_URL = process.env.YIXIAOER_API_URL || 'https://www.yixiaoer.cn/api';

/**
 * 通用发布引擎
 * 支持：文章(article)、视频(video)、图文/动态(image-text)、微信公众号(weixin-gongzhonghao)
 */
async function main() {
  const args = process.argv.slice(2);
  const params: Record<string, string> = {};
  
  // 1. 基础解析
  args.forEach(arg => {
    const match = arg.match(/^--([^=]+)=(.*)$/);
    if (match) {
      params[match[1]] = match[2];
    }
  });

  const type = params.type;
  const platformsStr = params.platforms;
  const accountIdsStr = params.account_ids;
  const title = params.title;
  const content = params.content || params.description; // 兼容 content 或 description
  
  if (!type || !platformsStr || !accountIdsStr) {
    console.error(JSON.stringify({ 
      error: "Missing core parameters: --type, --platforms, --account_ids" 
    }));
    process.exit(1);
  }

  if (!API_KEY) {
    console.error(JSON.stringify({ error: "Missing YIXIAOER_API_KEY environment variable." }));
    process.exit(1);
  }

  const platforms = platformsStr.split(',').map(p => p.trim());
  const accountIds = accountIdsStr.split(',').map(id => id.trim());

  try {
    // 2. 资源上传处理
    let videoKey = params.video_key;
    if (!videoKey && params.video_url) {
      videoKey = await uploadResource(params.video_url);
    }

    let coverKey = params.cover_key;
    if (!coverKey && params.cover_url) {
      coverKey = await uploadResource(params.cover_url);
    }

    let imageKeys: string[] = params.image_keys ? params.image_keys.split(',').map(k => k.trim()) : [];
    if (imageKeys.length === 0 && params.image_urls) {
      const urls = params.image_urls.split(',').map(u => u.trim());
      for (const url of urls) {
        const key = await uploadResource(url);
        imageKeys.push(key);
      }
    }

    // 3. 内容存证 (针对文章类)
    let publishContentId: string | undefined;
    if (type === 'article' || type === 'weixin-gongzhonghao') {
      publishContentId = Array.from({ length: 24 }, () => Math.floor(Math.random() * 16).toString(16)).join('');
      const wrappedContent = `<html><body>${content || ''}</body></html>`;
      const storageRes = await fetch(`${API_URL}/storages/articles`, {
        method: 'POST',
        headers: { 'Authorization': API_KEY, 'Content-Type': 'application/json' },
        body: JSON.stringify({ publishContentId, title, content: wrappedContent, contentHtml: wrappedContent })
      });
      if (!storageRes.ok) throw new Error(`Storage failed: ${await storageRes.text()}`);
    }

    // 4. 构建业务表单 (Document-Driven: 将所有非核心参数收集到 form 中)
    const coreArgs = [
      'type', 'platforms', 'account_ids', 'title', 'content', 'description',
      'video_url', 'video_key', 'cover_url', 'cover_key', 'image_urls', 'image_keys'
    ];
    
    const contentPublishForm: Record<string, any> = {};
    Object.keys(params).forEach(key => {
      if (!coreArgs.includes(key)) {
        const value = params[key];
        // 尝试解析 JSON (如 category, activity 等)
        try {
          if (value.startsWith('[') || value.startsWith('{')) {
            contentPublishForm[key] = JSON.parse(value);
          } else if (value === 'true') {
            contentPublishForm[key] = true;
          } else if (value === 'false') {
            contentPublishForm[key] = false;
          } else if (!isNaN(Number(value)) && value.trim() !== '') {
            contentPublishForm[key] = Number(value);
          } else {
            contentPublishForm[key] = value;
          }
        } catch (e) {
          contentPublishForm[key] = value;
        }
      }
    });

    // 补全基础字段
    if (title && !contentPublishForm.title) contentPublishForm.title = title;
    if (content && !contentPublishForm.content && (type === 'article' || type === 'weixin-gongzhonghao')) {
      contentPublishForm.content = content;
    }
    // 补全描述字段 (针对文章摘要或视频/图文描述)
    if (!contentPublishForm.desc && (params.desc || params.description)) {
      contentPublishForm.desc = params.desc || params.description;
    }
    if (!contentPublishForm.description && (params.description || params.desc)) {
      contentPublishForm.description = params.description || params.desc;
    }
    
    // 针对文章格式的补全 (核心处理)
    if (type === 'article' || type === 'weixin-gongzhonghao' || platforms.some(p => p === '微信公众号' || p === '公众号')) {
      if (!contentPublishForm.articles) {
        contentPublishForm.articles = [{
          title: title,
          content: content,
          contentHtml: `<html><body>${content}</body></html>`,
          digest: contentPublishForm.digest || title?.substring(0, 120),
          isDraft: contentPublishForm.pubType === 0
        }];
      }

      // 针对微信公众号的专项补全 (嵌套结构要求)
      if (platforms.some(p => p === '微信公众号' || p === '公众号')) {
        const wechatArticle = contentPublishForm.articles[0];
        if (wechatArticle) {
          // 确保 WeChat 必须的 cover 对象存在
          if (!wechatArticle.cover && coverKey) {
            wechatArticle.cover = { key: coverKey, width: 900, height: 383, size: 0 };
          }
          // 确保 WeChat 必须的 createType & authorName
          if (wechatArticle.createType === undefined) {
            wechatArticle.createType = contentPublishForm.original === true ? 1 : 0;
          }
          if (wechatArticle.authorName === undefined) {
            wechatArticle.authorName = contentPublishForm.author || '';
          }
          // 补全公众号专用属性 (默认值)
          if (wechatArticle.quickRepost === undefined) wechatArticle.quickRepost = 1;
          if (wechatArticle.quickPrivateMessage === undefined) wechatArticle.quickPrivateMessage = 1;
        }

        // 兼容某些接口可能需要的 contentList
        if (!contentPublishForm.contentList) {
          contentPublishForm.contentList = contentPublishForm.articles;
        }

        // 群发逻辑补全
        if (contentPublishForm.notify === false) {
          contentPublishForm.notifySubscribers = 0;
        } else if (contentPublishForm.notifySubscribers === undefined) {
          contentPublishForm.notifySubscribers = 1;
        }
      }

      // 文章共用封面
      if (!contentPublishForm.covers && coverKey) {
        contentPublishForm.covers = [{ key: coverKey, width: 1200, height: 800, size: 0 }];
      }
    }

    // 针对图文格式的补全
    if (type === 'image-text') {
      if (!contentPublishForm.description) contentPublishForm.description = content;
      if (imageKeys.length > 0 && !contentPublishForm.images) {
        contentPublishForm.images = imageKeys.map(key => ({ key, width: 1200, height: 800, size: 0 }));
      }
      // 抖音等平台可能需要 covers 数组
      if (!contentPublishForm.covers && coverKey) {
        contentPublishForm.covers = [{ key: coverKey, width: 1200, height: 800, size: 0 }];
      } else if (!contentPublishForm.covers && imageKeys.length > 0) {
        // 如果没有显式封面，使用第一张图作为封面
        contentPublishForm.covers = [{ key: imageKeys[0], width: 1200, height: 800, size: 0 }];
      }
    }

    // 针对视频格式的补全
    if (type === 'video' && !contentPublishForm.description) {
      contentPublishForm.description = content;
    }

    // 5. 构造任务 Body
    const platformForms: Record<string, any> = {};
    platforms.forEach(p => { platformForms[p] = contentPublishForm; });

    const taskBody: any = {
      desc: title || content?.substring(0, 30),
      platforms,
      publishType: (type === 'weixin-gongzhonghao' || type === 'article') ? 'article' : (type === 'image-text' ? 'imageText' : type),
      publishChannel: 'cloud',
      isDraft: contentPublishForm.pubType === 0,
      coverKey,
      publishArgs: {
        platformForms,
        accountForms: accountIds.map(accountId => {
          const accForm: any = { platformAccountId: accountId, contentPublishForm };
          if (publishContentId) accForm.publishContentId = publishContentId;
          if (coverKey) {
            accForm.coverKey = coverKey;
            accForm.cover = { key: coverKey, width: 1200, height: 800, size: 0 };
          }
          if (videoKey) {
            accForm.video = { key: videoKey, width: 1920, height: 1080, size: 0 };
          }
          if (imageKeys.length > 0) {
            accForm.images = imageKeys.map(key => ({ key, width: 1200, height: 800, size: 0 }));
          }
          return accForm;
        })
      }
    };
    
    // 如果是视频，顶级也需要 videoKey
    if (videoKey) taskBody.videoKey = videoKey;

    // 6. 最终提交
    const publishRes = await fetch(`${API_URL}/taskSets/v2`, {
      method: 'POST',
      headers: { 'Authorization': API_KEY, 'Content-Type': 'application/json' },
      body: JSON.stringify(taskBody)
    });

    if (!publishRes.ok) throw new Error(`Publishing failed: ${await publishRes.text()}`);
    const result = await publishRes.json();

    console.log(JSON.stringify(result, null, 2));

  } catch (error) {
    console.error(JSON.stringify({ 
      error: "Universal Publish Engine Error", 
      details: (error as Error).message 
    }));
    process.exit(1);
  }
}

main();
