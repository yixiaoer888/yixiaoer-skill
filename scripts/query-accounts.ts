/**
 * Query Account List (query-accounts.ts)
 * 
 * 获取蚁小二平台绑定的账号列表。
 * 调用方式: node query-accounts.ts --platform=抖音 --name=昵称 --page=1 --size=20
 */

async function main() {
  const args = process.argv.slice(2);
  const params: Record<string, string> = {};

  args.forEach(arg => {
    if (arg.startsWith('--')) {
      const [key, value] = arg.slice(2).split('=');
      if (key && value !== undefined) {
        params[key] = value;
      }
    }
  });

  const API_KEY = process.env.YIXIAOER_API_KEY;
  const API_URL = process.env.YIXIAOER_API_URL || 'https://www.yixiaoer.cn/api';

  if (!API_KEY) {
    console.error(JSON.stringify({ error: "Missing YIXIAOER_API_KEY environment variable" }));
    process.exit(1);
  }

  try {
    const url = new URL(`${API_URL}/platform-accounts`);
    
    // 默认值
    if (!params.page) params.page = '1';
    if (!params.size) params.size = '20';

    // 允许的参数列表
    const allowedParams = [
      'page', 'size', 'platform', 'name', 'group', 
      'loginStatus', 'isolation', 'parentId', 'time'
    ];

    allowedParams.forEach(key => {
      if (params[key]) {
        url.searchParams.append(key, params[key]);
      }
    });

    // 处理数组参数 (特殊处理)
    // 如果用户传了 --platforms=抖音,又传了 --platforms=快手，这里只能拿到最后一个
    // 如果需要支持数组，通常建议在 CLI 中用逗号分隔，如 --platforms=抖音,快手
    if (params['platforms']) {
      params['platforms'].split(',').forEach(p => url.searchParams.append('platforms[]', p));
    }
    if (params['platformType']) {
      params['platformType'].split(',').forEach(t => url.searchParams.append('platformType[]', t));
    }

    const response = await fetch(url.toString(), {
      method: 'GET',
      headers: {
        'Authorization': API_KEY,
        'Content-Type': 'application/json'
      }
    });

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`HTTP ${response.status}: ${errorText}`);
    }

    const result = await response.json();
    // 蚁小二接口通常返回 { statusCode: 0, data: { data: [...], totalSize: ... } }
    // 或者直接返回 data 数组
    console.log(JSON.stringify(result.data || result, null, 2));

  } catch (error) {
    console.error(JSON.stringify({ 
      error: "Failed to query accounts", 
      details: (error as Error).message 
    }));
    process.exit(1);
  }
}

main();
