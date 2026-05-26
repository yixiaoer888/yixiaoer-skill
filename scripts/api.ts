import * as fs from 'fs';
import * as path from 'path';

/**
 * 蚁小二 开放 API 助手
 */

export const API_KEY = process.env.YIXIAOER_API_KEY;
export const API_URL = process.env.YIXIAOER_API_URL || 'https://www.yixiaoer.cn/api';

function isBlankString(value: any): boolean {
  return typeof value !== 'string' || value.trim().length === 0;
}

/**
 * 解析并获取 --payload 参数中的 JSON 对象
 */
export function getPayload<T = any>(): T {
  const args = process.argv.slice(2);
  const payloadArg = args.find((a: string) => a.startsWith('--payload='))?.split('=')[1];

  if (!payloadArg) {
    throw new Error("Missing required parameter: --payload");
  }

  try {
    return JSON.parse(payloadArg) as T;
  } catch (error) {
    throw new Error(`Invalid JSON in --payload: ${(error as Error).message}`);
  }
}

/**
 * 低级 API 调用封装
 */
export async function callApi(endpoint: string, options: RequestInit = {}): Promise<any> {
  if (!API_KEY) {
    throw new Error("Missing YIXIAOER_API_KEY environment variable");
  }

  const url = endpoint.startsWith('http') ? endpoint : `${API_URL}${endpoint.startsWith('/') ? '' : '/'}${endpoint}`;

  const headers: Record<string, string> = {
    'Authorization': API_KEY,
    'Content-Type': 'application/json',
    ...((options.headers || {}) as Record<string, string>)
  };

  const response = await fetch(url, {
    ...options,
    headers
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(`HTTP ${response.status}: ${errorText}`);
  }

  const result = await response.json();
  return result;
}

/**
 * 助手函数: 通用资源上传
 */
export async function uploadResource(
  urlOrPath: string,
  bucket: string = 'cloud-publish',
  contentType: string,
  size?: number
): Promise<string> {
  if (!bucket) {
    throw new Error('Missing required bucket for upload');
  }

  let buffer: ArrayBuffer;
  let fileName = 'file.jpg';

  // 1. 获取资源内容
  if (urlOrPath.startsWith('http')) {
    const res = await fetch(urlOrPath);
    if (!res.ok) throw new Error(`HTTP error downloading file during sync upload: ${res.status}`);
    buffer = await res.arrayBuffer();
    const urlObj = new URL(urlOrPath);
    fileName = urlObj.pathname.split('/').pop() || 'image.jpg';
    if (!fileName.includes('.')) fileName += '.jpg';
  } else {
    const absolutePath = path.isAbsolute(urlOrPath) ? urlOrPath : path.resolve(process.cwd(), urlOrPath);
    if (!fs.existsSync(absolutePath)) {
      throw new Error(`Local file not found: ${absolutePath}`);
    }
    const fileBuffer = fs.readFileSync(absolutePath);
    buffer = fileBuffer.buffer.slice(fileBuffer.byteOffset, fileBuffer.byteOffset + fileBuffer.byteLength);
    fileName = path.basename(absolutePath);
  }

  // 2. 获取预签名上传地址 (使用 callApi)
  const queryParams = new URLSearchParams();
  queryParams.append('fileKey', fileName);
  if (contentType) queryParams.append('contentType', contentType);
  if (size) queryParams.append('size', String(size));

  const uploadInfo = await callApi(`/storages/${bucket}/upload-url?${queryParams.toString()}`);
  const data = uploadInfo.data || uploadInfo;
  const { serviceUrl, key } = data;

  if (!serviceUrl) {
    throw new Error(`Invalid upload info response: ${JSON.stringify(uploadInfo)}`);
  }

  // 3. 执行 PUT 上传 (注意: Content-Type 必须与获取 URL 时一致)
  const putRes = await fetch(serviceUrl, {
    method: 'PUT',
    body: buffer,
    headers: { 'Content-Type': contentType }
  });

  if (!putRes.ok) {
    const errorText = await putRes.text();
    if (errorText.includes('SignatureDoesNotMatch')) {
      throw new Error(`Failed to upload to OSS (SignatureDoesNotMatch): 预签名的 Content-Type (${contentType}) 与实际上传的 Header 不匹配。请确保两者完全一致。`);
    }
    throw new Error(`Failed to upload to OSS: ${errorText}`);
  }

  return key;
}

/**
 * 统一错误处理并输出到标准输出
 */
export function handleError(error: any, context: string, errorCode: string = 'YIXIAOER_REMOTE_ERR') {
  let finalErrorCode = errorCode;
  const message = error instanceof Error ? error.message : String(error);

  // 识别特定错误类型并自动归类
  if (message.includes('Missing required') || message.includes('Invalid JSON')) {
    finalErrorCode = 'YIXIAOER_USAGE_ERR';
  } else if (message.includes('不支持') || message.includes('Unsupported') || message.includes('Invalid platform')) {
    finalErrorCode = 'YIXIAOER_USAGE_ERR';
  } else if (message.includes('API_KEY')) {
    finalErrorCode = 'YIXIAOER_AUTH_ERR';
  }

  const response = {
    success: false,
    errorCode: finalErrorCode,
    message: message.includes('不支持') ? message : `Failed to ${context}`,
    details: message,
    suggestion: message.includes('不支持')
      ? "请检查发布类型 (publishType) 组合是否正确。例如小红书仅支持 'image-text' 或 'video'，不支持 'article'。"
      : "请依次检查: 1. 技能版本号是否一致; 2. 请求参数是否符合 DTO 规范; 3. 查阅 [避坑指南](./docs/troubleshooting-guide.md)。"
  };

  console.error(JSON.stringify(response, null, 2));
  process.exit(1);
}

/**
 * 平台与类型及其动作的支持矩阵预检
 */
function validateSupport(payload: any) {
  const { action, publishType, platforms } = payload;

  // 1. 基础动作校验
  const supportedActions = [
    'publish', 'save-draft', 'accounts', 'upload', 'material', 'records', 'details',
    'categories', 'activities', 'locations', 'music', 'music-category',
    'collections', 'groups', 'goods', 'hot-events', 'challenges',
    'miniapps', 'syncapps', 'games', 'proxies', 'proxy-areas',
    'account-overviews', 'content-overviews', 'update-account'
  ];

  if (action && !supportedActions.includes(action)) {
    throw new Error(`不支持的动作 (Unsupported action): "${action}"。请参考 SKILL.md 中的动作列表。`);
  }

  // 2. 必填字段快速预检 (针对 publish 动作)
  if (action === 'publish' || action === 'save-draft') {
    if (!payload.publishArgs) {
      throw new Error(`"publish" 动作缺少必填的 "publishArgs" 对象。请参考对应平台的公开 DTO 架构文档。`);
    }
    const { accountForms } = payload.publishArgs;
    if (!accountForms || !Array.isArray(accountForms) || accountForms.length === 0) {
      throw new Error(`"publishArgs" 中缺少必填的 "accountForms" 数组。你必须至少指定一个目标账号。`);
    }

    accountForms.forEach((form: any, index: number) => {
      // 检查账号 ID
      if (!form.platformAccountId && !form.account_id) {
        throw new Error(`accountForms[${index}] 缺少必填的 "platformAccountId"。请先通过 "accounts" 动作获取并将 ID 填入。`);
      }

      // 检查发布表单
      if (!form.contentPublishForm) {
        throw new Error(`accountForms[${index}] 缺少必填的 "contentPublishForm"。这是存放标题、描述等核心信息的容器。`);
      }

      // 深度检查必填核心业务字段 (通识性必填)
      const cpf = form.contentPublishForm;
      if (!cpf.title && payload.publishType !== 'video') {
        // 视频类型在某些平台可能可选标题或由 description 替代，但文章和图文通常必填
        console.warn(`[Warning] accountForms[${index}].contentPublishForm 中未检测到 "title"。绝大多数平台发布文章/图文均要求标题，请核对。`);
      }

      // 资源存在性校验
      if (payload.publishType === 'video' && !form.video) {
        throw new Error(`"publishType" 为 "video" 时，accountForms[${index}] 必须包含 "video" 资源信息（key, size, width, height, duration）。`);
      }
      if (payload.publishType === 'article' && isBlankString(cpf.content)) {
        throw new Error(`"publishType" 为 "article" 时，contentPublishForm 必须包含非空白的 "content" 正文。`);
      }
      if (payload.publishType === 'image-text' && (!form.images || form.images.length === 0)) {
        throw new Error(`"publishType" 为 "image-text" 时，accountForms[${index}] 必须包含 "images" 数组。`);
      }
    });

    // 检查全局平台声明
    if (!payload.platforms || (Array.isArray(payload.platforms) && payload.platforms.length === 0)) {
      throw new Error(`"publish" 动作缺少必填的 "platforms" (平台列表) 声明。`);
    }
  }
}

function buildTaskSetBody(payload: Record<string, any>, forceDraft: boolean = false) {
  const { action: _action, ...body } = payload;
  if (!('publishChannel' in body)) {
    body.publishChannel = 'cloud';
  }
  if (forceDraft) {
    body.isDraft = true;
  }
  return body;
}

/**
 * 主执行入口 (Execution Entry)
 */
async function main() {
  // 1. 检查是否为直接执行 (ts-node 或 node scripts/api.ts)
  const isMain = process.argv[1]?.replace(/\\/g, '/').endsWith('scripts/api.ts') ||
    process.argv[1]?.replace(/\\/g, '/').endsWith('scripts/api');

  if (!isMain) return;

  try {
    const payload = getPayload();

    // 动作与平台支持预检测
    validateSupport(payload);

    const action = payload.action;

    if (!action) {
      throw new Error("Missing required field: action in payload");
    }

    let result: any;
    switch (action) {
      case 'publish':    // 内容发布
        {
          const taskSetBody = buildTaskSetBody(payload);
          // 在 publish 动作中，强制关闭 isDraft，确保其始终走向发布流程
          // 平台草稿通过 contentPublishForm.pubType=0 控制，而非 isDraft=true
          taskSetBody.isDraft = false;

          result = await callApi('/taskSets/v2', {
            method: 'POST',
            body: JSON.stringify(taskSetBody)
          });
        }
        break;

      case 'save-draft': // 保存为蚁小二草稿
        {
          const taskSetBody = buildTaskSetBody(payload, true); // 强制 isDraft = true
          result = await callApi('/taskSets/drafts', {
            method: 'PUT',
            body: JSON.stringify(taskSetBody)
          });
        }
        break;

      case 'accounts': // 账号列表
        const accountUrl = new URL(`${API_URL}/v2/platform/accounts`);
        Object.keys(payload).forEach(key => {
          if (key !== 'action') accountUrl.searchParams.append(key, String(payload[key]));
        });
        result = await callApi(accountUrl.toString(), { method: 'GET' });
        break;

      case 'upload': // 资源上传
        if (!payload.contentType) {
          throw new Error("Missing required field: contentType for action: upload");
        }
        const uploadBucket = payload.bucket || 'cloud-publish';
        const uploadKey = await uploadResource(payload.url, uploadBucket, payload.contentType, payload.size);
        result = {
          key: uploadKey,
          name: payload.url.startsWith('http') ? new URL(payload.url).pathname.split('/').pop() : payload.url.split(/[/\\]/).pop(),
          bucket: uploadBucket
        };
        break;

      case 'material': // 上传到素材库
        if (!payload.filePath) throw new Error("Missing filePath for action: material");
        if (!payload.fileName) throw new Error("Missing fileName for action: material");
        if (typeof payload.width !== 'number') throw new Error("Missing numeric width for action: material");
        if (typeof payload.height !== 'number') throw new Error("Missing numeric height for action: material");
        if (!payload.type) throw new Error("Missing type for action: material");
        result = await callApi('/material', {
          method: 'POST',
          body: JSON.stringify({
            thumbPath: payload.thumbPath,
            filePath: payload.filePath,
            fileName: payload.fileName,
            width: payload.width,
            height: payload.height,
            type: payload.type
          })
        });
        break;

      case 'records': // 发布记录
        const recordUrl = new URL(`${API_URL}/v2/taskSets`);
        Object.keys(payload).forEach(key => {
          if (key !== 'action') recordUrl.searchParams.append(key, String(payload[key]));
        });
        result = await callApi(recordUrl.toString(), { method: 'GET' });
        break;

      case 'details': // 任务详情
        if (!payload.task_set_id) throw new Error("Missing task_set_id for action: details");
        result = await callApi(`/v2/taskSets/${payload.task_set_id}/tasks`, { method: 'GET' });
        break;

      case 'categories': // 分类查询
        const categoryUrl = new URL(`${API_URL}/platform-accounts/${payload.account_id}/categories`);
        categoryUrl.searchParams.append('publishType', payload.type || 'video');
        result = await callApi(categoryUrl.toString(), { method: 'GET' });
        break;

      case 'activities': // 活动查询
        const activityUrl = new URL(`${API_URL}/platform-accounts/${payload.account_id}/activities`);
        activityUrl.searchParams.append('publishType', payload.type || 'video');
        if (payload.categoryId) activityUrl.searchParams.append('categoryId', payload.categoryId);
        if (payload.keyword || payload.keyWord) activityUrl.searchParams.append('keyWord', payload.keyword || payload.keyWord);
        result = await callApi(activityUrl.toString(), { method: 'GET' });
        break;

      case 'locations': // POI 搜索
        const locationUrl = new URL(`${API_URL}/platform-accounts/${payload.account_id}/location`);
        if (payload.keyword || payload.keyWord) locationUrl.searchParams.append('keyWord', payload.keyword || payload.keyWord);
        if (payload.type) locationUrl.searchParams.append('locationType', String(payload.type));
        if (payload.nextPage) locationUrl.searchParams.append('nextPage', payload.nextPage);
        result = await callApi(locationUrl.toString(), { method: 'GET' });
        break;

      case 'music': // 音乐素材
        if (!payload.account_id) throw new Error("Missing account_id for action: music");
        const musicUrl = new URL(`${API_URL}/platform-accounts/${payload.account_id}/music`);
        if (payload.keyword || payload.keyWord) musicUrl.searchParams.append('keyWord', payload.keyword || payload.keyWord);
        if (payload.categoryId) musicUrl.searchParams.append('categoryId', payload.categoryId);
        if (payload.categoryName) musicUrl.searchParams.append('categoryName', payload.categoryName);
        if (payload.nextPage) musicUrl.searchParams.append('nextPage', payload.nextPage);
        result = await callApi(musicUrl.toString(), { method: 'GET' });
        break;

      case 'music-category': // 音乐分类
        if (!payload.account_id) throw new Error("Missing account_id for action: music-category");
        result = await callApi(`/platform-accounts/${payload.account_id}/music/category`, { method: 'GET' });
        break;

      case 'collections': // 合集查询
        const collectionUrl = new URL(`${API_URL}/platform-accounts/${payload.account_id}/collections`);
        collectionUrl.searchParams.append('publishType', payload.type || 'video');
        result = await callApi(collectionUrl.toString(), { method: 'GET' });
        break;

      case 'proxies': // 代理列表
        const proxyUrl = new URL(`${API_URL}/proxys`);
        proxyUrl.searchParams.append('size', String(payload.size || 9999));
        result = await callApi(proxyUrl.toString(), { method: 'GET' });
        break;

      case 'proxy-areas': // 内置代理地区列表 (默认代理)
        result = await callApi('/daili/areas', { method: 'GET' });
        break;

      case 'update-account': // 更新账号信息 (如设置代理)
        if (!payload.account_id) throw new Error("Missing account_id for action: update-account");
        // 支持更新 kuaidailiArea 或 proxyId
        const updateBody: any = {};
        if ('kuaidailiArea' in payload) updateBody.kuaidailiArea = payload.kuaidailiArea;
        if ('proxyId' in payload) updateBody.proxyId = payload.proxyId;
        if ('remark' in payload) updateBody.remark = payload.remark;
        if ('groups' in payload) updateBody.groups = payload.groups;

        result = await callApi(`/platform-accounts/${payload.account_id}`, {
          method: 'PATCH',
          body: JSON.stringify(updateBody)
        });
        break;

      case 'content-overviews': // 作品数据
        const contentOverviewUrl = new URL(`${API_URL}/contents/overviews`);
        Object.keys(payload).forEach(key => {
          if (key !== 'action') {
            if (Array.isArray(payload[key])) {
              payload[key].forEach((v: any) => contentOverviewUrl.searchParams.append(key, String(v)));
            } else {
              contentOverviewUrl.searchParams.append(key, String(payload[key]));
            }
          }
        });
        result = await callApi(contentOverviewUrl.toString(), { method: 'GET' });
        break;

      case 'account-overviews': // 账号数据 (V2)
        const accountOverviewUrl = new URL(`${API_URL}/platform-accounts/overviews-v2`);
        Object.keys(payload).forEach(key => {
          if (key !== 'action') {
            if (Array.isArray(payload[key])) {
              payload[key].forEach((v: any) => accountOverviewUrl.searchParams.append(key, String(v)));
            } else {
              accountOverviewUrl.searchParams.append(key, String(payload[key]));
            }
          }
        });
        result = await callApi(accountOverviewUrl.toString(), { method: 'GET' });
        break;

      case 'groups': // 群聊列表
        result = await callApi(`/platform-accounts/${payload.account_id}/group-chats`, { method: 'GET' });
        break;

      case 'goods': // 商品列表
        const goodsUrl = new URL(`${API_URL}/platform-accounts/${payload.account_id}/goods`);
        if (payload.keyword || payload.keyWord) goodsUrl.searchParams.append('keyWord', payload.keyword || payload.keyWord);
        if (payload.nextPage) goodsUrl.searchParams.append('nextPage', payload.nextPage);
        result = await callApi(goodsUrl.toString(), { method: 'GET' });
        break;

      case 'hot-events': // 热点列表
        const hotEventUrl = new URL(`${API_URL}/platform-accounts/${payload.account_id}/hot-events`);
        hotEventUrl.searchParams.append('publishType', payload.type || 'video');
        result = await callApi(hotEventUrl.toString(), { method: 'GET' });
        break;

      case 'challenges': // 挑战列表
        const challengeUrl = new URL(`${API_URL}/platform-accounts/${payload.account_id}/challenges`);
        challengeUrl.searchParams.append('publishType', payload.type || 'video');
        if (payload.keyword || payload.keyWord) challengeUrl.searchParams.append('keyWord', payload.keyword || payload.keyWord);
        if (payload.nextPage) challengeUrl.searchParams.append('nextPage', payload.nextPage);
        result = await callApi(challengeUrl.toString(), { method: 'GET' });
        break;

      case 'miniapps': // 小程序列表
        const miniappUrl = new URL(`${API_URL}/platform-accounts/${payload.account_id}/mini-apps`);
        if (payload.keyword || payload.keyWord) miniappUrl.searchParams.append('keyWord', payload.keyword || payload.keyWord);
        result = await callApi(miniappUrl.toString(), { method: 'GET' });
        break;

      case 'syncapps': // 同步应用列表
        result = await callApi(`/platform-accounts/${payload.account_id}/sync-apps`, { method: 'GET' });
        break;

      case 'games': // 游戏列表
        const gameUrl = new URL(`${API_URL}/platform-accounts/${payload.account_id}/games`);
        if (payload.keyword || payload.keyWord) gameUrl.searchParams.append('keyWord', payload.keyword || payload.keyWord);
        result = await callApi(gameUrl.toString(), { method: 'GET' });
        break;



      default:
        throw new Error(`Unsupported action: ${action}`);
    }

    const finalResult = {
      success: true,
      action: action,
      version: "1.6.4", // 与 SKILL.md 保持同步
      data: result.data || result
    };

    console.log(JSON.stringify(finalResult, null, 2));

  } catch (error) {
    handleError(error, "execute api action");
  }
}

main();

export { };
