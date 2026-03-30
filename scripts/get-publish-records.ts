/**
 * 查询发布记录 (get-publish-records.ts)
 * 
 * 查询任务集 (TaskSet) 的历史发布记录。
 * 调用方式: node scripts/get-publish-records.ts --payload='{"page":1,"size":10}'
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
    const url = new URL(`${API_URL}/v2/taskSets`);

    // 将 Payload 中的字段映射为 SearchParams
    Object.keys(payload).forEach(key => {
      url.searchParams.append(key, String(payload[key]));
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
      error: "Failed to query publish records", 
      details: (error as Error).message 
    }));
    process.exit(1);
  }
}

main();

export {};
