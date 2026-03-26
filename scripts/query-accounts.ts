/**
 * Query Account List (query-accounts.ts)
 * 
 * 获取蚁小二平台绑定的账号列表。
 * 该脚本可通过 Node.js 直接运行，无需辅助 package。
 * 
 * 调用方式: node query-accounts.ts [flags]
 */

async function main() {
  const args = process.argv.slice(2);
  const platformArg = args.find(a => a.startsWith('--platform='))?.split('=')[1];
  const statusArg = args.find(a => a.startsWith('--status='))?.split('=')[1];

  // 模拟从环境变量或配置文件获取 API 信息
  const API_KEY = process.env.YIXIAOER_API_KEY || 'MOCK_KEY';
  const API_URL = process.env.YIXIAOER_API_URL || 'https://www.yixiaoer.cn/api';

  try {
    // 使用 Node.js 18+ 内置的 fetch，不依赖外部库
    /*
    const response = await fetch(`${API_URL}/accounts?platform=${platformArg || ''}&status=${statusArg || ''}`, {
      headers: { 'Authorization': `Bearer ${API_KEY}` }
    });
    const data = await response.json();
    console.log(JSON.stringify(data, null, 2));
    */

    // 演示模拟输出
    const mockAccounts = [
      { accountId: "acc_001", uid: "douyin_user_1", platform: "DouYin", name: "运营大咖", status: "active" },
      { accountId: "acc_002", uid: "sph_user_2", platform: "VideoChannel", name: "生活记录官", status: "expired" },
      { accountId: "acc_003", uid: "xhs_user_3", platform: "XiaoHongShu", name: "美妆博主", status: "active" }
    ];

    const filtered = mockAccounts.filter(acc => {
      let match = true;
      if (platformArg && acc.platform !== platformArg) match = false;
      if (statusArg && acc.status !== statusArg) match = false;
      return match;
    });

    console.log(JSON.stringify(filtered, null, 2));

  } catch (error) {
    console.error(JSON.stringify({ error: "Failed to fetch accounts", details: (error as Error).message }));
    process.exit(1);
  }
}

main();
