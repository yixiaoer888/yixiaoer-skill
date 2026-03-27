/**
 * 获取发布分类 (get-publish-categories.ts)
 * 
 * 获取特定账号在特定发布下可以使用的分类列表。
 * 调用方式: node get-publish-categories.ts --account_id=XXX --type=article
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
  const type = argMap.type;

  const API_KEY = process.env.YIXIAOER_API_KEY;
  const API_URL = process.env.YIXIAOER_API_URL || 'https://www.yixiaoer.cn/api';

  if (!API_KEY || !accountId || !type) {
    console.error(JSON.stringify({ 
      error: "Missing required parameters: --account_id, --type and YIXIAOER_API_KEY environment variable" 
    }));
    process.exit(1);
  }

  try {
    // 标准 API 路径: GET platform-accounts/:id/categories
    const response = await fetch(`${API_URL}/platform-accounts/${accountId}/categories?publishType=${type}`, {
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
    console.log(JSON.stringify(result.data || result, null, 2));

  } catch (error) {
    console.error(JSON.stringify({ 
      error: "Failed to get categories", 
      details: (error as Error).message 
    }));
    process.exit(1);
  }
}

main();

export {};


