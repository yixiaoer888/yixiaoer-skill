/**
 * Get Publish Activities (get-publish-activities.ts)
 * 
 * 获取特定账号在特定发布模态下的征文活动。
 * 调用方式: node get-publish-activities.ts --account_id=XXX --type=article
 */

async function main() {
  const args = process.argv.slice(2);
  const accountId = args.find(a => a.startsWith('--account_id='))?.split('=')[1];
  const type = args.find(a => a.startsWith('--type='))?.split('=')[1];

  const API_KEY = process.env.YIXIAOER_API_KEY;
  const API_URL = process.env.YIXIAOER_API_URL || 'https://www.yixiaoer.cn/api';

  if (!API_KEY || !accountId || !type) {
    console.error(JSON.stringify({ 
      error: "Missing required parameters: --account_id, --type and YIXIAOER_API_KEY environment variable"
    }));
    process.exit(1);
  }

  try {
    const response = await fetch(`${API_URL}/web/config-data/activity-tasks`, {
      method: 'POST',
      headers: {
        'Authorization': API_KEY,
        'Content-Type': 'application/json',
        'x-account-id': accountId
      },
      body: JSON.stringify({
        openAccountId: accountId,
        publishType: type
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
      error: "Failed to get activities", 
      details: (error as Error).message 
    }));
    process.exit(1);
  }
}

main();
