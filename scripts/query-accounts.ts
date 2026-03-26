/**
 * Query Account List (query-accounts.ts)
 * 
 * 获取蚁小二平台绑定的账号列表。
 * 调用方式: node query-accounts.ts --platform=DouYin
 */

async function main() {
  const args = process.argv.slice(2);
  const platform = args.find(a => a.startsWith('--platform='))?.split('=')[1];
  const page = args.find(a => a.startsWith('--page='))?.split('=')[1] || '1';
  const size = args.find(a => a.startsWith('--size='))?.split('=')[1] || '20';

  const API_KEY = process.env.YIXIAOER_API_KEY;
  const API_URL = process.env.YIXIAOER_API_URL || 'https://www.yixiaoer.cn/api';

  if (!API_KEY) {
    console.error(JSON.stringify({ error: "Missing YIXIAOER_API_KEY environment variable" }));
    process.exit(1);
  }

  try {
    const url = new URL(`${API_URL}/platform-accounts`);
    url.searchParams.append('page', page);
    url.searchParams.append('size', size);
    if (platform) url.searchParams.append('platform', platform);

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
    // 假设结果在 data 字段中，直接输出
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
