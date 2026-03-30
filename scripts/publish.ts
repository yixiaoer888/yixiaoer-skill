/**
 * 通用发布引擎 (Universal Publishing Engine)
 * 
 * 仅支持通过 --payload 传入完整的符合 CloudTaskPushRequest 结构的 JSON 对象。
 * 
 * 调用方式:
 * node scripts/publish.ts --payload='{...}'
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
    const taskBody = JSON.parse(payloadArg);

    const response = await fetch(`${API_URL}/taskSets/v2`, {
      method: 'POST',
      headers: {
        'Authorization': API_KEY,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(taskBody)
    });

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`HTTP ${response.status}: ${errorText}`);
    }

    const result = await response.json();
    console.log(JSON.stringify(result, null, 2));

  } catch (error) {
    console.error(JSON.stringify({ 
      error: "Failed to submit publish task", 
      details: (error as Error).message 
    }));
    process.exit(1);
  }
}

main();

export {};
