/**
 * 获取征文活动 (get-publish-activities.ts)
 * 
 * 获取特定账号在特定发布类型下的征文活动列表。
 * 调用方式: node scripts/get-publish-activities.ts --payload='{"account_id":"XXX","type":1}'
 */

async function main() {
  const API_KEY = process.env.YIXIAOER_API_KEY;
  const API_URL = process.env.YIXIAOER_API_URL || 'https://www.yixiaoer.cn/api';

  const args = process.argv.slice(2);
  const payloadArg = args.find(a => a.startsWith('--payload='))?.split('=')[1];

  if (!API_KEY) {
    console.error(JSON.stringify({ error: "Missing YIXIAOER_API_KEY environment variable" }));
    process.exit(1);
  }

  if (!payloadArg) {
    console.error(JSON.stringify({ error: "Missing required parameter: --payload" }));
    process.exit(1);
  }

  try {
    const payload = JSON.parse(payloadArg);
    const accountId = payload.account_id;

    if (!accountId) {
      throw new Error("Missing required field: account_id in payload");
    }

    const response = await fetch(`${API_URL}/web/config-data/activity-tasks`, {
      method: 'POST',
      headers: {
        'Authorization': API_KEY,
        'Content-Type': 'application/json',
        'x-account-id': accountId
      },
      body: JSON.stringify({
        openAccountId: accountId,
        publishType: payload.type || 1,
        categoryId: payload.categoryId,
        keyWord: payload.keyWord
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
      error: "Failed to query activities", 
      details: (error as Error).message 
    }));
    process.exit(1);
  }
}

main();

export {};
