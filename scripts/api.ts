import * as fs from 'fs';
import * as path from 'path';

/**
 * 蚁小二 开放 API 助手
 */

export const API_KEY = process.env.YIXIAOER_API_KEY;
export const API_URL = process.env.YIXIAOER_API_URL || 'https://www.yixiaoer.cn/api';

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
export async function callApi(endpoint: string, options: RequestInit = {}) {
  if (!API_KEY) {
    throw new Error("Missing YIXIAOER_API_KEY environment variable");
  }

  const url = endpoint.startsWith('http') ? endpoint : `${API_URL}${endpoint.startsWith('/') ? '' : '/'}${endpoint}`;
  
  const headers: HeadersInit = {
    'Authorization': API_KEY,
    'Content-Type': 'application/json',
    ...(options.headers || {})
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
export async function uploadResource(urlOrPath: string, bucket: string = 'cloud-publish'): Promise<string> {
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
  const uploadInfo = await callApi(`/storages/${bucket}/upload-url?fileKey=${encodeURIComponent(fileName)}`);
  const data = uploadInfo.data || uploadInfo;
  const { serviceUrl, key } = data;

  if (!serviceUrl) {
    throw new Error(`Invalid upload info response: ${JSON.stringify(uploadInfo)}`);
  }

  // 3. 执行 PUT 上传
  const putRes = await fetch(serviceUrl, {
    method: 'PUT',
    body: buffer,
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' }
  });

  if (!putRes.ok) {
    throw new Error(`Failed to upload to OSS: ${await putRes.text()}`);
  }

  return key;
}

/**
 * 统一错误处理并输出到标准输出
 */
export function handleError(error: any, context: string) {
  console.error(JSON.stringify({ 
    error: `Failed to ${context}`, 
    details: error instanceof Error ? error.message : String(error) 
  }, null, 2));
  process.exit(1);
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
    const action = payload.action;

    if (!action) {
      throw new Error("Missing required field: action in payload");
    }

    let result: any;
    switch (action) {
      case 'publish': // 内容发布
        result = await callApi('/taskSets/v2', { 
          method: 'POST', 
          body: JSON.stringify(payload) 
        });
        break;

      case 'accounts': // 账号列表
        const accountUrl = new URL(`${API_URL}/v2/platform/accounts`);
        Object.keys(payload).forEach(key => {
          if (key !== 'action') accountUrl.searchParams.append(key, String(payload[key]));
        });
        result = await callApi(accountUrl.toString(), { method: 'GET' });
        break;

      case 'upload': // 资源上传
        const key = await uploadResource(payload.url, payload.bucket);
        result = { 
          key, 
          name: payload.url.startsWith('http') ? new URL(payload.url).pathname.split('/').pop() : payload.url.split(/[/\\]/).pop() 
        };
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
        result = await callApi('/web/config-data/category-tasks', { 
          method: 'POST', 
          headers: { 'x-account-id': payload.account_id },
          body: JSON.stringify({ openAccountId: payload.account_id, publishType: payload.type || 'video' })
        });
        break;

      case 'activities': // 活动查询
        result = await callApi('/web/config-data/activity-tasks', { 
          method: 'POST', 
          headers: { 'x-account-id': payload.account_id },
          body: JSON.stringify({ openAccountId: payload.account_id, publishType: payload.type || 1, categoryId: payload.categoryId, keyWord: payload.keyWord })
        });
        break;

      case 'locations': // POI 搜索
        result = await callApi('/web/config-data/location-tasks', { 
          method: 'POST', 
          headers: { 'x-account-id': payload.account_id },
          body: JSON.stringify({ openAccountId: payload.account_id, keyWord: payload.keyword || '', locationType: parseInt(payload.type || '1'), nextPage: "" })
        });
        break;

      case 'music': // 音乐素材
        result = await callApi('/web/config-data/music-tasks', { 
          method: 'POST', 
          headers: { 'x-account-id': payload.account_id },
          body: JSON.stringify({ openAccountId: payload.account_id, keyWord: payload.keyword || '', nextPage: "" })
        });
        break;

      case 'collections': // 合集查询
        result = await callApi('/web/config-data/collection-tasks', { 
          method: 'POST', 
          headers: { 'x-account-id': payload.account_id },
          body: JSON.stringify({ openAccountId: payload.account_id })
        });
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
        result = await callApi('/web/config-data/group-tasks', { 
          method: 'POST', 
          headers: { 'x-account-id': payload.account_id },
          body: JSON.stringify({ openAccountId: payload.account_id })
        });
        break;

      case 'goods': // 商品列表
        result = await callApi('/web/config-data/goods-tasks', { 
          method: 'POST', 
          headers: { 'x-account-id': payload.account_id },
          body: JSON.stringify({ openAccountId: payload.account_id, keyWord: payload.keyword || '', nextPage: "" })
        });
        break;

      case 'hot-events': // 热点列表
        result = await callApi('/web/config-data/hotevent-tasks', { 
          method: 'POST', 
          headers: { 'x-account-id': payload.account_id },
          body: JSON.stringify({ openAccountId: payload.account_id })
        });
        break;

      case 'challenges': // 挑战列表
        result = await callApi('/web/config-data/challenge-tasks', { 
          method: 'POST', 
          headers: { 'x-account-id': payload.account_id },
          body: JSON.stringify({ openAccountId: payload.account_id, keyWord: payload.keyword || '', nextPage: "" })
        });
        break;

      case 'miniapps': // 小程序列表
        result = await callApi('/web/config-data/miniapp-tasks', { 
          method: 'POST', 
          headers: { 'x-account-id': payload.account_id },
          body: JSON.stringify({ openAccountId: payload.account_id, keyWord: payload.keyword || '', nextPage: "" })
        });
        break;

      case 'syncapps': // 同步应用列表
        result = await callApi('/web/config-data/syncapp-tasks', { 
          method: 'POST', 
          headers: { 'x-account-id': payload.account_id },
          body: JSON.stringify({ openAccountId: payload.account_id })
        });
        break;

      case 'games': // 游戏列表
        result = await callApi('/web/config-data/game-tasks', { 
          method: 'POST', 
          headers: { 'x-account-id': payload.account_id },
          body: JSON.stringify({ openAccountId: payload.account_id, keyWord: payload.keyword || '', nextPage: "" })
        });
        break;

      default:
        throw new Error(`Unsupported action: ${action}`);
    }

    console.log(JSON.stringify(result.data || result, null, 2));

  } catch (error) {
    handleError(error, "execute api action");
  }
}

main();

export {};
