/**
 * 查询发布记录 (get-publish-records.ts)
 * 
 * 查询任务集 (TaskSet) 的历史发布记录。
 * 调用方式: node get-publish-records.ts --page=1 --size=10 [--status=allsuccessful] [--publish_type=video] [--keywords=test]
 */

async function main() {
  const args = process.argv.slice(2);
  const argMap: Record<string, string> = {};

  args.forEach((arg: string) => {
    const [key, value] = arg.split('=');
    if (key.startsWith('--')) {
      argMap[key.substring(2)] = value;
    }
  });

  const page = argMap.page || '1';
  const size = argMap.size || '10';
  const publishType = argMap.publish_type;
  const keywords = argMap.keywords;
  const status = argMap.status;
  const startTime = argMap.start_time;
  const endTime = argMap.end_time;

  const API_KEY = process.env.YIXIAOER_API_KEY;
  const API_URL = process.env.YIXIAOER_API_URL || 'https://www.yixiaoer.cn/api';

  if (!API_KEY) {
    console.error(JSON.stringify({ error: "Missing YIXIAOER_API_KEY environment variable." }));
    process.exit(1);
  }

  try {
    const url = new URL(`${API_URL}/v2/taskSets`);
    url.searchParams.append('page', page);
    url.searchParams.append('size', size);
    
    if (publishType) url.searchParams.append('publishType', publishType);
    if (keywords) url.searchParams.append('keyWords', keywords);
    if (status) url.searchParams.append('taskSetStatus', status);
    if (startTime) url.searchParams.append('publishStartTime', startTime);
    if (endTime) url.searchParams.append('publishEndTime', endTime);

    const res = await fetch(url.toString(), {
      method: 'GET',
      headers: {
        'Authorization': API_KEY,
        'Content-Type': 'application/json'
      }
    });

    if (!res.ok) {
      const errorText = await res.text();
      throw new Error(`Failed to fetch publish records. HTTP ${res.status}: ${errorText}`);
    }
    
    const result = await res.json();
    console.log(JSON.stringify(result.data || result, null, 2));

  } catch (error) {
    console.error(JSON.stringify({ 
      error: "Query Publish Records Error", 
      details: (error as Error).message 
    }));
    process.exit(1);
  }
}

main();

export {};

