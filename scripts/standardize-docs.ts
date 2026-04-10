import fs from 'fs';
import path from 'path';

const baseDir = 'c:/work/yixiaoer/yixiaoer-universal/yixiaoer-skill/docs/publish';

type Category = 'article' | 'video' | 'image-text';

function getTemplates(category: Category, platform: string) {
  const templates = {
    video: {
      trigger: `## 触发场景 (Trigger)
- **意图辨析**：用户指定在“${platform}”平台分发视频内容时触发。
- **典型提示词**：
  - “把这个视频发布到${platform}”
  - “同步视频到${platform}”`,
      logic: `## 执行逻辑 (Logic Flow)
1. **意图确认**：确认目标平台为${platform}。
2. **参数装配**：识别并填充标题、描述等平台特定字段至 \`contentPublishForm\`。
3. **指令执行**：调用 \`node scripts/api.ts\`。`
    },
    article: {
      trigger: `## 触发场景 (Trigger)
- **意图辨析**：用户指定在“${platform}”平台发布文章内容时触发。
- **典型提示词**：
  - “发布这篇文章到${platform}”
  - “并在${platform}上同步更新”`,
      logic: `## 执行逻辑 (Logic Flow)
1. **内容处理**：确保文章正文符合${platform}要求的格式。
2. **参数装配**：提取标题、正文及封面信息至 \`contentPublishForm\`。
3. **指令执行**：调用 \`node scripts/api.ts\`。`
    },
    'image-text': {
      trigger: `## 触发场景 (Trigger)
- **意图辨析**：用户指定在“${platform}”平台发布图文动态时触发。
- **典型提示词**：
  - “发几张图到${platform}”
  - “同步这条动态到${platform}”`,
      logic: `## 执行逻辑 (Logic Flow)
1. **资源校验**：确保所有图片均已上传并获得 Key。
2. **参数装配**：填充描述及图片列表至 \`contentPublishForm\`。
3. **指令执行**：调用 \`node scripts/api.ts\`。`
    }
  };
  return templates[category];
}

function standardize(filePath: string) {
  let content = fs.readFileSync(filePath, 'utf-8');
  if (content.includes('## 触发场景 (Trigger)')) return;

  const fileName = path.basename(filePath, '.md');
  if (fileName === 'index') return;

  const categoryMatch = filePath.match(/publish[\\/](article|video|image-text)/);
  if (!categoryMatch) return;
  const category = categoryMatch[1] as Category;
  
  const platformName = fileName.charAt(0).toUpperCase() + fileName.slice(1);
  const { trigger, logic } = getTemplates(category, platformName);

  // Find the place after the title or IMPORTANT block
  const lines = content.split('\n');
  let insertPos = 1;
  for (let i = 0; i < lines.length; i++) {
    if (lines[i].includes('> [!IMPORTANT]') || lines[i].includes('> [!CAUTION]')) {
      // Find the end of this block
      while (i < lines.length && lines[i].trim() !== '') i++;
      insertPos = i + 1;
      break;
    }
  }

  lines.splice(insertPos, 0, '\n' + trigger + '\n\n' + logic + '\n');
  
  // Fix mandatory status in tables
  let newContent = lines.join('\n');
  newContent = newContent.replace(/\| `action` \| `string` \| 否 \|/g, '| `action` | `string` | **是** |');
  newContent = newContent.replace(/\| `formType` \| `string` \| 否 \|/g, '| `formType` | `string` | **是** |');
  newContent = newContent.replace(/\| `pubType` \| `number` \| 否 \|/g, '| `pubType` | `number` | **是** |');
  
  fs.writeFileSync(filePath, newContent);
  console.log(`Updated: ${filePath}`);
}

function walk(dir: string) {
  const files = fs.readdirSync(dir);
  for (const file of files) {
    const fullPath = path.join(dir, file);
    if (fs.statSync(fullPath).isDirectory()) {
      walk(fullPath);
    } else if (file.endsWith('.md')) {
      standardize(fullPath);
    }
  }
}

walk(baseDir);
