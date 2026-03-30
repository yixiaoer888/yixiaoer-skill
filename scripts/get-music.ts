/**
 * 获取音乐素材 (get-music.ts)
 * 
 * 获取发布时可选的音乐列表。
 * 调用方式: node scripts/get-music.ts --payload='{"account_id":"XXX","keyword":"周杰伦"}'
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

    const response = await fetch(`${API_URL}/web/config-data/music-tasks`, {
      method: 'POST',
      headers: {
        'Authorization': API_KEY,
        'Content-Type': 'application/json',
        'x-account-id': accountId
      },
      body: JSON.stringify({
        openAccountId: accountId,
        keyWord: payload.keyword || '',
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
      error: "Failed to query music", 
      details: (error as Error).message 
    }));
    process.exit(1);
  }
}

main();

export {};
