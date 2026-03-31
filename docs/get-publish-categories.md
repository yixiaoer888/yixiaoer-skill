# 获取发布分类 (Get Publish Categories)

该指令用于获取特定账号在特定分发模态（视频、文章等）下的分类列表。这些分类 ID 是发布时 `category` 参数的基础。

## 场景描述 (Usage)

- "列出我的抖音账号 A 可选的所有视频分类。"
- "看看这个百家号账号都有哪些文章分类。"

## 参数列表 (Payload Properties)

| 字段名 | 类型 | 是否必填 | 描述 |
| :--- | :--- | :--- | :--- |
| `account_id` | `string` | **是** | 蚁小二媒体账号 ID |
| `type` | `ContentTypeEnum` | **是** | 发布模态：`video` (视频), `article` (文章), `dynamic` (动态), `imageText` (图文) |

### 枚举值定义

#### ContentTypeEnum (发布类型)
- `video`: 视频
- `article`: 文章
- `dynamic`: 动态 (爬虫使用)
- `imageText`: 图文 (前端使用)

## 脚本逻辑 (Backend)

- **脚本路径**: `../scripts/get-publish-categories.ts`
- **实际接口**: `GET https://www.yixiaoer.cn/api/v2/platform/accounts/:id/categories`
- **环境要求**: 需设置 `YIXIAOER_API_KEY` 环境变量。

## 输出结果 (Output)

脚本返回分类列表数组，每个对象包含以下核心字段：

| 字段名 | 类型 | 描述 |
| :--- | :--- | :--- |
| `yixiaoerId` | `string` | 分类在本系统的唯一标识 |
| `yixiaoerName` | `string` | 分类名称 |
| `child` | `array` | 子分类列表（若支持层级分类，如抖音/快手） |
| `raw` | `object` | 平台原始分类数据 |

### 调用指令 (Command)
```bash
node scripts/api.ts --payload='{"action":"categories","account_id":"64dxxx","type":"article"}'
```
