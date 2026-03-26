const API_KEY = process.env.YIXIAOER_API_KEY;
const API_URL = process.env.YIXIAOER_API_URL || 'https://www.yixiaoer.cn/api';

async function main() {
  const args = process.argv.slice(2);
  const page = args.find(a => a.startsWith('--page='))?.split('=')[1] || '1';
  const size = args.find(a => a.startsWith('--size='))?.split('=')[1] || '10';
  const publishType = args.find(a => a.startsWith('--publish_type='))?.split('=')[1];
  const keyWords = args.find(a => a.startsWith('--keywords='))?.split('=')[1];

  if (!API_KEY) {
    console.error(JSON.stringify({ error: "Missing YIXIAOER_API_KEY environment variable." }));
    process.exit(1);
  }

  try {
    const url = new URL(`${API_URL}/v2/taskSets`);
    url.searchParams.append('page', page);
    url.searchParams.append('size', size);
    if (publishType) url.searchParams.append('publishType', publishType);
    if (keyWords) url.searchParams.append('keyWords', keyWords);

    const res = await fetch(url.toString(), {
      method: 'GET',
      headers: {
        'Authorization': API_KEY,
        'Content-Type': 'application/json'
      }
    });

    if (!res.ok) throw new Error(`Failed to fetch publish records: ${await res.text()}`);
    const result = await res.json();

    console.log(JSON.stringify(result, null, 2));

  } catch (error) {
    console.error(JSON.stringify({ 
      error: "Query Publish Records Error", 
      details: (error as Error).message 
    }));
    process.exit(1);
  }
}

main();
