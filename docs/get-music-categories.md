# 📄 获取音乐分类 Query 参数 (Get Music Categories Query)

获取在选择音乐素材时可选的分类列表。用于帮助用户按“场景”、“心情”或“曲风”快速定位目标音乐。

> [!TIP]
> **多级路径匹配**：对于具有层级结构的音乐分类，Agent **必须**使用返回对象中的 `child` 路径数组，以确保发布表单的完整性。

## 1. 触发场景 (Trigger)

- **意图辨析**：用户没有明确的歌名，但知道想要的曲风类型。
- **典型提示词**：
  - “看看抖音里有哪些音乐分类？”
  - “获取热门配乐的类别”

## 2. 交互协议 (Interactive Protocol)

1. **链式调用**：获取分类 ID 后，Agent 应主动引导用户使用该 ID 继续调用 `music` 接口获取具体曲目。
2. **铺平逻辑感知**：Agent 应知晓返回的列表已由脚本自动**铺平 (Flatten)**，可直接在扁平列表中进行关键词搜索。

## 3. 参数定义 (Parameters)

| 字段名 | 类型 | 必填 | 默认值 | 描述 |
| :--- | :--- | :--- | :--- | :--- |
| **`action`** | `string` | **是** | `music-category` | 固定值。 |
| **`account_id`** | `string` | **是** | - | 蚁小二账号 ID (`platformAccountId`)。 |

### 3.1 返回结果结构

| 字段 | 类型 | 描述 |
| :--- | :--- | :--- |
| `yixiaoerId` | `string` | 内部分类 ID |
| `yixiaoerName` | `string` | 分类名称 |
| `child` | `array` | **完整路径对象数组**。包含父子层级关系，需透传至发布表单。 |
| `raw` | `object` | 原始平台返回的分类对象结构。 |

## 4. 执行指令示例 (Command)

```bash
node scripts/api.ts --payload='{"action":"music-category","account_id":"67fb2f1735eeb3cf31db3d65"}'
```

## 5. 常见问题排查 (Troubleshooting)

| 现象 | 可能原因 | 处理建议 |
| :--- | :--- | :--- |
| **分类列表为空** | 该平台当前未定义公共音乐分类。 | 建议用户直接通过 `music` 接口进行关键词搜索。 |
| **ID 校验失败** | 跨平台使用了错误的分类 ID。 | 确保 `categoryId` 与当前发布的目标账号平台一致。 |

---
> [!NOTE]
> **自动化处理**：脚本 `api.ts` 会自动抓取多级嵌套分类并将其扁平化，Agent 开发者无需再行编写树状递归逻辑。
