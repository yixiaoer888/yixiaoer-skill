/**
 * 获取音乐列表 (get-music.ts)
 * 
 * 获取发布时可选的音乐列表。
 * 调用方式: node scripts/get-music.ts --account_id=XXX --keyword=周杰伦
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
  const keyword = argMap.keyword || '';

  const API_KEY = process.env.YIXIAOER_API_KEY;
  const API_URL = process.env.YIXIAOER_API_URL || 'https://www.yixiaoer.cn/api';

  if (!API_KEY || !accountId) {
    console.error(JSON.stringify({ 
      error: "Missing required parameters: --account_id and YIXIAOER_API_KEY environment variable"
    }));
    process.exit(1);
  }

  try {
    const response = await fetch(`${API_URL}/web/config-data/music-tasks`, {
      method: 'POST',
      headers: {
        'Authorization': API_KEY,
        'Content-Type': 'application/json',
        'x-account-id': accountId
      },
      body: JSON.stringify({
        openAccountId: accountId,
        keyWord: keyword,
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
