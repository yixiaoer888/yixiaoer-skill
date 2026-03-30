/**
 * 获取地理位置 (get-locations.ts)
 * 
 * 获取发布时可选的地理位置列表。
 * 调用方式: node scripts/get-locations.ts --account_id=XXX --keyword=深圳 --type=1
 */

async function main() {
  const args = process.argv.slice(2);
  const argMap: Record<string, string> = {};

  args.forEach((arg: string) => {
    const [key, value] = arg.split('=');
    if (key.startsWith('--')) {
      argMap[key.substring(2)] = value;
    }
  });

  const accountId = argMap.account_id;
  const keyword = argMap.keyword || '';
  const type = parseInt(argMap.type || '1'); // 默认搜地点

  const API_KEY = process.env.YIXIAOER_API_KEY;
  const API_URL = process.env.YIXIAOER_API_URL || 'https://www.yixiaoer.cn/api';

  if (!API_KEY || !accountId) {
    console.error(JSON.stringify({ 
      error: "Missing required parameters: --account_id and YIXIAOER_API_KEY environment variable"
    }));
    process.exit(1);
  }

  try {
    const response = await fetch(`${API_URL}/web/config-data/location-tasks`, {
      method: 'POST',
      headers: {
        'Authorization': API_KEY,
        'Content-Type': 'application/json',
        'x-account-id': accountId
      },
      body: JSON.stringify({
        openAccountId: accountId,
        keyWord: keyword,
        locationType: type,
        nextPage: ""
      })
    });

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`HTTP ${response.status}: ${errorText}`);
    }

    const result = await response.json();
    console.log(JSON.stringify(result.data || result, null, 2));

  } catch (error) {
    console.error(JSON.stringify({ 
      error: "Failed to query locations", 
      details: (error as Error).message 
    }));
    process.exit(1);
  }
}

main();

export {};
