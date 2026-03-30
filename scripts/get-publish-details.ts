/**
 * 获取发布详情 (get-publish-details.ts)
 * 
 * 获取特定任务集的详细执行记录。
 * 调用方式: node scripts/get-publish-details.ts --payload='{"task_set_id":"TASK_SET_ID"}'
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
    const taskSetId = payload.task_set_id;

    if (!taskSetId) {
      throw new Error("Missing required field: task_set_id in payload");
    }

    const res = await fetch(`${API_URL}/v2/taskSets/${taskSetId}/tasks`, {
      method: 'GET',
      headers: {
        'Authorization': API_KEY,
        'Content-Type': 'application/json'
      }
    });

    if (!res.ok) {
      const errorText = await res.text();
      throw new Error(`HTTP ${res.status}: ${errorText}`);
    }

    const result = await res.json();
    console.log(JSON.stringify(result.data || result, null, 2));

  } catch (error) {
    console.error(JSON.stringify({ 
      error: "Failed to query publish details", 
      details: (error as Error).message 
    }));
    process.exit(1);
  }
}

main();

export {};
