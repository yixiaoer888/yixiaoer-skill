/**
 * Query Account List (query-accounts.ts)
 * 
 * 获取蚁小二平台绑定的账号列表。
 * 仅支持通过 --payload 传入完整的过滤参数 JSON 对象。
 * 
 * 调用方式:
 * node scripts/query-accounts.ts --payload='{"platform":"抖音","name":"昵称","page":1,"size":20}'
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
    const url = new URL(`${API_URL}/v2/platform/accounts`);

    // 将 Payload 中的字段映射为 SearchParams
    Object.keys(payload).forEach(key => {
      const val = payload[key];
      if (Array.isArray(val)) {
        val.forEach(v => url.searchParams.append(`${key}[]`, String(v)));
      } else {
        url.searchParams.append(key, String(val));
      }
    });

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

export {};
