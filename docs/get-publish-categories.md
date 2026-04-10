# 获取发布分类/话题 (Get Categories)

获取特定平台账号支持的文章分类、视频分类或话题标签列表。

## 触发场景 (Trigger)
- **意图辨析**：在准备发布内容时，为了确保 `platformSettings` 中的分类/话题符合平台标准，不产生非法值，必须预先查询合法值。
- **典型提示词**：
  - “获取这个抖音号的视频分类”
  - “查看百家号支持的文章类别”
  - “查询可以挂载的话题”

## 参数定义 (Parameters)

### 参数列表 (Payload Properties)

| 字段名 | 类型 | 是否必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `action` | `string` | **是** | 固定值：`categories` |
| `account_id` | `string` | **是** | 蚁小二账号 ID (ObjectId) |
| `type` | `string` | 否 | 发布类型：`video` (默认) 或 `article` |

## 执行逻辑 (Logic Flow)
1. **身份锚定**：识别目标账号 `account_id`（通过 `accounts` action 获取）。
2. **场景对齐**：根据发布内容决定 `type`。
3. **参数装配**：构造 `action: "categories"` 负载。
4. **指令执行**：调用 `node scripts/api.ts --payload='{...}'`。
5. **值注入**：将返回的分类 `id` 或 `name` 填入发布 Payload 的对应位置。

## 返回数据说明 (Response Details)

返回包含分类对象（`Category` 结构）的数组。
- **铺平处理 (Flatten)**：脚本会自动将多级嵌套的分类（二级分类等）铺平为一级数组，方便 Agent 直接匹配。
- **`child` 字段 (重要)**：每一个分类对象都会附带一个 `child` 数组。
  - **内容**：包含了从根节点到当前节点的完整路径序列。例如：`[{ "yixiaoerName": "科技", ... }, { "yixiaoerName": "人工智能", ... }]`。
  - **填表规范**：对于要求多级分类的平台（如 Bilibili），Agent **必须**直接将该 `child` 数组填入发布 Payload 的 `category` 字段。

### 复杂对象：CategoryItem
- `yixiaoerId`: 内部分类 ID。
- `yixiaoerName`: 分类名称。
- `child`: **完整路径对象数组**。如果分类有父子级关系，发布表单时通常需要在此处填入整个生成的 `child`。
- `raw`: 原始平台返回的分类对象。

## 调用指令 (Command)

```bash
node scripts/api.ts --payload='{"action":"categories","account_id":"YOUR_ACCOUNT_ID","type":"video"}'
```
