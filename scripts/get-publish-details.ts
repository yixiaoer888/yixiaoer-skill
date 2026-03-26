const API_KEY = process.env.YIXIAOER_API_KEY;
const API_URL = process.env.YIXIAOER_API_URL || 'https://www.yixiaoer.cn/api';

async function main() {
  const args = process.argv.slice(2);
  const id = args.find(a => !a.startsWith('--')) || args.find(a => a.startsWith('--id='))?.split('=')[1];

  if (!id) {
    console.error(JSON.stringify({ error: "Missing required parameter: taskSetId" }));
    process.exit(1);
  }

  if (!API_KEY) {
    console.error(JSON.stringify({ error: "Missing YIXIAOER_API_KEY environment variable." }));
    process.exit(1);
  }

  try {
    const res = await fetch(`${API_URL}/v2/taskSets/${id}/tasks`, {
      method: 'GET',
      headers: {
        'Authorization': API_KEY,
        'Content-Type': 'application/json'
      }
    });

    if (!res.ok) throw new Error(`Failed to fetch publish details: ${await res.text()}`);
    const result = await res.json();

    console.log(JSON.stringify(result, null, 2));

  } catch (error) {
    console.error(JSON.stringify({ 
      error: "Query Publish Details Error", 
      details: (error as Error).message 
    }));
    process.exit(1);
  }
}

main();
