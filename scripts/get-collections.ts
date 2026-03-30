/**
 * 获取合集列表 (get-collections.ts)
 * 
 * 获取账号已创建的合集列表。
 * 仅支持通过 --payload 传入 JSON 参数。
 * 
 * 调用方式:
 * node scripts/get-collections.ts --payload='{"account_id":"ACCOUNT_ID"}'
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

    const response = await fetch(`${API_URL}/web/config-data/collection-tasks`, {
      method: 'POST',
      headers: {
        'Authorization': API_KEY,
        'Content-Type': 'application/json',
        'x-account-id': accountId
      },
      body: JSON.stringify({ openAccountId: accountId })
    });

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`HTTP ${response.status}: ${errorText}`);
    }

    const result = await response.json();
    console.log(JSON.stringify(result.data || result, null, 2));

  } catch (error) {
    console.error(JSON.stringify({ 
      error: "Failed to query collections", 
      details: (error as Error).message 
    }));
    process.exit(1);
  }
}

main();

export {};
