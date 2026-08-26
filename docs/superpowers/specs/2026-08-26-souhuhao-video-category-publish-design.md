# 搜狐号视频分类发布修复设计

## 背景

搜狐号视频发布要求 `category` 使用父分类到子分类的完整级联路径。现有 CLI 在发布边界只会把查询结果对象包装成前端的 `id/text/raw` 结构；当调用方已经传入 `id/text`，但 `raw` 直接是搜狐平台原始对象时，包装逻辑会误认为结构完整，最终服务端读取不到预期的嵌套 `raw.id`，产生：

```text
Cannot read properties of undefined (reading 'id')
```

另一个风险是调用方只传入末级分类。搜狐号需要父子路径，单独传入子分类会被平台判定为频道不存在。

## 目标与范围

本次只改 CLI 的“搜狐号视频”发布流程，目标是：

1. 客户可以通过分类查询结果明确看到父分类、子分类及其 ID。
2. 发布前使用目标搜狐号账号的最新分类数据，不依赖手写或过期分类对象。
3. 客户传入完整查询对象、前端包装对象、平台原始分类对象，或只传入可识别的末级分类时，CLI 自动生成正确的父子路径和嵌套 `raw`。
4. 分类不存在、名称歧义或查询失败时，在 CLI 侧阻止发布，返回结构化错误、可用分类路径和下一步查询命令。
5. 不改变其他平台的分类归一化逻辑，也不修改调用方原始 payload。

不在范围内：修改搜狐平台服务端、改变视频文件处理、改变非搜狐号平台的分类语义、增加交互式输入流程。

## 方案

### 分类数据展示

保留 `yxer query categories <account_id> --type video --json` 的 JSON 输出协议，继续返回 CLI 规范化后的分类树。增加 `--paths` 视图供客户直接查看可发布的完整路径；该选项只改变 `data` 内的查询视图，不改变默认输出：

```bash
yxer query categories <account_id> --type video --paths --json
```

`--paths` 的 `data` 结构为 `{ "categories": <分类树>, "paths": [...] }`，每条路径至少包含 `path`、任意深度的 `nodes` 数组，以及可直接回填的完整 `category` 数组；两层路径额外提供 `parentId`、`parentName`、`childId`、`childName` 便于客户阅读：

```json
{
  "categories": [
    {
      "yixiaoerId": "43",
      "yixiaoerName": "生活",
      "raw": { "id": 43 },
      "child": [
        { "yixiaoerId": "58", "yixiaoerName": "生活百态", "raw": { "id": 58, "channelId": 43 } }
      ]
    }
  ],
  "paths": [
    {
      "path": "生活 > 生活百态",
      "parentId": "43",
      "parentName": "生活",
      "childId": "58",
      "childName": "生活百态",
      "nodes": [
        { "id": "43", "name": "生活" },
        { "id": "58", "name": "生活百态" }
      ],
      "category": ["<完整父分类对象>", "<完整子分类对象>"]
    }
  ]
}
```

stdout 仍然只输出 JSON；不向 stdout 混入纯文本树。同步更新搜狐号命令帮助和技能文档，明确 `child` 就是子分类，发布时需要保留从父到子的完整路径。发布前校验失败时，错误详情额外返回由最新查询结果生成的可用路径列表，例如 `生活 > 生活百态`，避免客户只能看到远端的模糊异常。

### 发布前解析

在 `internal/api` 增加搜狐号视频分类解析能力，供发布准备阶段调用：

1. 从每个 `accountForms[].platformAccountId` 读取目标账号。
2. 搜狐号视频发布必须针对每个 `accountForms[].platformAccountId` 独立调用一次 `categories` 查询接口。一个账号的分类结果不得复用给另一个账号。查询失败、返回空树或结构不完整时禁止回退到 payload、缓存或手写分类；使用 CLI 已有的规范化分类树作为唯一数据源。
3. 从 `accountForms[].contentPublishForm.category` 数组的每个候选项提取 ID、名称或原始对象：支持 canonical 查询对象、正确的前端包装对象、旧的 `raw` 直挂平台对象，以及只提供 `id`/`text`/`name` 的候选项。匹配顺序为 ID 精确匹配，再按名称匹配；同名候选超过一个时返回歧义错误。
4. 以最后一个候选项作为选中的末级项，向上构造从根到末级的完整路径；如果 payload 同时提供父级候选，则逐级校验其 ID/名称与查询路径一致。父级分类全部直接来自查询结果，不能通过硬编码或只复制平台 ID 生成。
5. 同一候选项同时提供多个身份字段时，非空 ID 必须与其 `raw.id` 或查询对象 ID 一致，非空名称必须与查询对象名称一致；发生冲突时返回 `sohuhao_category_invalid`，不能静默选择其中一个字段。
6. 查询 canonical 分类对象的契约是：`yixiaoerId` 和 `yixiaoerName` 是 CLI 稳定身份字段，`raw` 是完整的平台原始分类对象；`raw.id` 必须存在，且字符串化后等于 `yixiaoerId`。搜狐号的 `raw` 不要求有 `text`，但不得丢弃平台返回的其他字段（例如 `channelId`）。
7. 将路径转换为服务端需要的结构：每一项包含 `id`、`text` 和完整的 CLI 查询对象 `raw`；如果存在子项，保留 `children`。最终请求示例为：

   ```json
   "category": [
     {
       "id": "43",
       "text": "生活",
       "raw": { "yixiaoerId": "43", "yixiaoerName": "生活", "raw": { "id": 43, "cmsChannelId": 206 } }
     },
     {
       "id": "58",
       "text": "生活百态",
       "raw": { "yixiaoerId": "58", "yixiaoerName": "生活百态", "raw": { "id": 58, "channelId": 43 } }
     }
   ]
   ```

   这里有两个不同阶段的字段路径：请求线上是 `wireItem.raw.raw.id`；服务端的 `parsePlatformCategory` 抽取 `wireItem.raw` 后，搜狐实现读取的是 `category.raw.id`。两者指向同一个平台分类 ID。对每个发布分类项断言 `wireItem.raw.raw.id` 存在，才能避免原始的 `undefined.id`。
8. 解析结果写入发布请求副本。用户传入的 payload 和预校验所使用的原始对象不被修改。

如果已有分类是正确的完整路径，也重新以当前查询结果中的对象组装，确保 ID、名称和嵌套 `raw` 一致。搜狐号 CLI 发布只有在上述查询和解析成功后才允许进入 API 发布边界；API 边界的通用分类包装不承担查询、补父级或绕过失败的职责。

### 错误处理

分类查询或解析失败时不创建发布任务。错误必须使用现有 `yxerrors.Error` 结构，并包含失败账号、请求分类和可用路径。定义以下稳定的 `code/category`：

| 场景 | code | category | retryable |
| --- | --- | --- | --- |
| 分类查询失败 | `sohuhao_category_query_failed` | `sohuhao_category` | `true` |
| 查询结果为空或结构无效 | `sohuhao_category_data_invalid` | `sohuhao_category` | `false` |
| 分类不存在 | `sohuhao_category_not_found` | `sohuhao_category` | `false` |
| 分类名称歧义 | `sohuhao_category_ambiguous` | `sohuhao_category` | `false` |

错误详情至少为：

```json
{
  "accountId": "<account_id>",
  "requested": [{ "id": "58", "text": "生活百态" }],
  "availablePaths": [
    { "path": "生活 > 生活百态", "parentId": "43", "childId": "58" }
  ],
  "cause": "<optional underlying cause>"
}
```

分类选择错误时 `availablePaths` 由本次最新查询生成；查询失败、空树或结构无效时 `availablePaths` 固定为空数组，并在 `cause` 中保留底层查询/结构原因，不使用旧 payload 代替。三层及以上分类统一使用 `nodes` 表示完整路径，`path` 是节点名称用 ` > ` 连接的展示文本；`parentId`/`childId` 仅作为两层路径的便捷字段，不限制路径深度：

```json
{
  "path": "生活 > 生活百态 > 访谈",
  "nodes": [
    { "id": "43", "name": "生活" },
    { "id": "58", "name": "生活百态" },
    { "id": "99", "name": "访谈" }
  ],
  "category": ["<完整对象>", "<完整对象>", "<完整对象>"]
}
```

完整错误输出仍由现有输出层包装为：

```json
{
  "ok": false,
  "error": {
    "code": "sohuhao_category_not_found",
    "category": "sohuhao_category",
    "message": "搜狐号视频分类不存在",
    "hint": "请从 availablePaths 中选择有效的父子分类路径。",
    "retryable": false,
    "nextCommand": "yxer query categories <account_id> --type video --json",
    "details": {
      "accountId": "<account_id>",
      "requested": [{ "id": "missing" }],
      "availablePaths": []
    }
  }
}
```

查询失败时同一结构中的 `code` 为 `sohuhao_category_query_failed`、`retryable` 为 `true`，并保留 `details.cause`。每个错误都提供 `hint`；`nextCommand` 固定为 `yxer query categories <account_id> --type video --json`。

这样错误会在本地发布前出现，而不是创建任务后才收到搜狐号的 JavaScript 异常。

## 数据流

```text
payload.json
    ↓
发布准备：读取账号与分类选择
    ↓
query categories（每个搜狐号账号，失败不得回退）
    ↓
按 ID/名称匹配末级分类，补齐父级路径
    ↓
生成 id/text/raw（raw 内含完整 canonical 对象）
    ↓
schema + preflight + dry-run 展示最终请求（stdout 仅 JSON）
    ↓
真实发布；只有分类查询与解析成功后才进入 API 边界
```

## 测试与验收

采用 TDD，先增加失败测试，再实现：

1. `--paths` 输出父分类、子分类、ID、任意深度的 `nodes` 和完整路径，且 stdout 始终是可解析 JSON。
2. API 分类查询结果包含父分类和子分类时，末级选择生成完整两级路径。
3. 第二次失败任务的旧格式（`raw` 直接为平台原始对象）被修复为嵌套 canonical `raw`，并断言最终每项 `raw.raw.id` 存在。
4. 只传末级子分类时自动补齐父分类；三层以上分类也能生成完整路径。
5. ID 匹配优先于名称匹配；名称重复时返回 `sohuhao_category_ambiguous` 并列出可用路径。
6. 分类不存在、查询失败、空结果和 raw 缺少 `id` 时不调用 `/taskSets/v2`，错误包含稳定字段、可用路径和 `nextCommand`。
7. 同一候选项的 `id`、`raw.id`、名称不一致，或 payload 父子路径与查询结果不一致时，返回 `sohuhao_category_invalid`，不创建任务。
8. 多账号 payload 为每个 `platformAccountId` 分别返回不同分类树时，断言每个 `contentPublishForm.category` 只使用自身账号对应的查询结果；查询失败不得使用另一个账号的分类或 payload 旧值。
9. 正确完整路径仍能发布；非搜狐平台完全跳过新解析器并保持现有分类归一化测试通过。
10. dry-run 会查询分类但绝不调用 `/taskSets/v2`，stdout 只输出 JSON，输出补齐后的发布请求和分类归一化信息；查询失败时返回同样的结构化错误，诊断信息只能写入 stderr；真实发布只在所有校验通过后执行。

验收标准：使用搜狐号账号的真实分类查询结果，原先会报 `Cannot read properties of undefined (reading 'id')` 的 payload 在 `validate`、`--dry-run` 和真实发布流程中均能得到可识别的父子分类结构；分类不合法时不会创建远端任务。
