/**
 * 获取征文活动 (get-publish-activities.ts)
 * 
 * 获取特定账号在特定发布类型下的征文活动列表。
 * 调用方式: node get-publish-activities.ts --account_id=XXX --type=video [--categoryId=YYY] [--keyWord=ZZZ]
 */

async function main() {
  const args = process.argv.slice(2);
  const argMap: Record<string, string> = {};

  args.forEach(arg => {
    const [key, value] = arg.split('=');
    if (key.startsWith('--')) {
      argMap[key.substring(2)] = value;
    }
  });

  const accountId = argMap.account_id;
  const type = argMap.type;
  const categoryId = argMap.categoryId;
  const keyWord = argMap.keyWord;

  const API_KEY = process.env.YIXIAOER_API_KEY;
  const API_URL = process.env.YIXIAOER_API_URL || 'https://www.yixiaoer.cn/api';

  if (!API_KEY || !accountId || !type) {
    console.error(JSON.stringify({ 
      error: "Missing required parameters: --account_id, --type and YIXIAOER_API_KEY environment variable"
    }));
    process.exit(1);
  }

  try {
    // 构造查询参数
    const queryParams = new URLSearchParams({
      publishType: type
    });
    if (categoryId) queryParams.append('categoryId', categoryId);
    if (keyWord) queryParams.append('keyWord', keyWord);

    // 标准 API 路径: GET platform-accounts/:id/activities
    const response = await fetch(`${API_URL}/platform-accounts/${accountId}/activities?${queryParams.toString()}`, {
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
      error: "Failed to get activities", 
      details: (error as Error).message 
    }));
    process.exit(1);
  }
}

main();

export {};


