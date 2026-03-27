/**
 * 获取当前团队信息 (get-team-info.ts)
 * 
 * 获取当前用户的团队元数据及限额配额。
 * 调用方式: node get-team-info.ts
 */

async function main() {
  const API_KEY = process.env.YIXIAOER_API_KEY;
  const API_URL = process.env.YIXIAOER_API_URL || 'https://www.yixiaoer.cn/api';

  if (!API_KEY) {
    console.error(JSON.stringify({ error: "Missing YIXIAOER_API_KEY environment variable" }));
    process.exit(1);
  }

  try {
    // 步骤 1: 获取当前用户信息，提取最后一次使用的团队 ID (latestTeamId)
    const userInfoRes = await fetch(`${API_URL}/users/info`, {
      method: 'GET',
      headers: {
        'Authorization': API_KEY,
        'Content-Type': 'application/json'
      }
    });

    if (!userInfoRes.ok) {
        const errorText = await userInfoRes.text();
        throw new Error(`Failed to fetch user info. HTTP ${userInfoRes.status}: ${errorText}`);
    }

    const userData = await userInfoRes.json();
    const userInfo = userData.data || userData;
    const teamId = userInfo.latestTeamId;

    if (!teamId) {
      console.error(JSON.stringify({ error: "User is not associated with any team" }));
      process.exit(1);
    }

    // 步骤 2: 获取该特定团队的详细信息
    const teamInfoRes = await fetch(`${API_URL}/teams/${teamId}`, {
      method: 'GET',
      headers: {
        'Authorization': API_KEY,
        'Content-Type': 'application/json'
      }
    });

    if (!teamInfoRes.ok) {
        const errorText = await teamInfoRes.text();
        throw new Error(`Failed to fetch team info. HTTP ${teamInfoRes.status}: ${errorText}`);
    }

    const teamData = await teamInfoRes.json();
    const finalResult = {
      user: {
        id: userInfo.id,
        nickName: userInfo.nickName,
      },
      team: teamData.data || teamData
    }

    console.log(JSON.stringify(finalResult, null, 2));

  } catch (error) {
    console.error(JSON.stringify({ 
      error: "Failed to query team information", 
      details: (error as Error).message 
    }));
    process.exit(1);
  }
}

main();

export {};

