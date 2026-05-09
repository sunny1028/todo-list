# 列表共享 — 设计文档

日期：2026-05-10

## 概述

为待办事项应用增加列表共享功能。用户可生成邀请码将列表分享给他人，支持"仅查看"和"可编辑"两种权限。被分享者通过邀请码加入，在侧边栏中看到共享列表。

采用方案 C：轻量起步，先做共享链接 + 简单权限，后续可演进到工作空间模型。

## 数据模型

新增两张表：

### list_shares

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint, PK | |
| list_id | uint, FK→lists, index | 被共享的列表 |
| code | string, unique, not null | 8位随机邀请码 |
| permission | string, not null | `view` 或 `edit` |
| created_at | time | |

### list_share_members

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint, PK | |
| share_id | uint, FK→list_shares | |
| user_id | uint, FK→users | 加入的用户 |
| joined_at | time | |

级联删除：删除列表 → 级联 share → 级联 members。撤销 share 时删除所有关联 members。

用户可访问列表的条件：`list.user_id == user_id` 或 `list 存在有效 share 且 user 在 members 中`。

## API

在 `/api/lists` 组下新增：

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| POST | `/api/lists/:id/share` | 创建共享，返回 code | 列表所有者 |
| GET | `/api/lists/:id/share` | 查看共享信息（code, permission, members） | 列表所有者 |
| DELETE | `/api/lists/:id/share` | 撤销共享，级联删除 members | 列表所有者 |
| POST | `/api/lists/join` | 通过 `{ code }` 加入列表 | 已登录用户 |

### 权限校验规则

- 创建/查看/撤销共享：仅列表所有者
- 加入列表：任何已登录用户（code 有效即可）
- 访问共享列表内容时，根据 permission 限制写入操作
- `view` 权限：可读所有数据，不可创建/修改/删除
- `edit` 权限：完整读写

### 现有接口改动

- `GET /api/lists`：返回自有列表 + 已加入的共享列表
- `GET /api/todos`：校验增加 member 权限（所有者或 member 均可访问）
- `POST/PUT/PATCH/DELETE /api/todos`：增加写入权限校验（edit member 或 owner）
- 所有 todo 相关操作均需校验：owner 或 share member (edit)

## 前端

### 共享弹窗（ShareDialog.vue）

- 列表所有者可见分享按钮（侧边栏列表项旁，悬停显示）
- 弹窗内容：邀请码 + 复制链接按钮 + 权限选择 + 成员列表 + 撤销按钮
- 共享链接格式：`{origin}/join?code=A3K9M2X1`

### 加入列表（JoinDialog.vue）

- 侧边栏底部「加入列表」按钮 → 弹出输入框 → 输入邀请码
- 加入成功后侧边栏刷新显示共享列表
- 支持通过 URL `/join?code=XXX` 直接加入

### 侧边栏共享列表

- 共享列表显示 👥 图标与自有列表区分
- 选中共享列表时，根据权限禁用/启用编辑操作
- 只读列表：TodoForm 隐藏、编辑按钮禁用、删除按钮禁用

## 边界情况

- 同一用户不能通过邀请码加入自己的列表
- 同一用户重复加入同一列表 → 幂等，返回已加入
- 撤销共享后，所有 members 的侧边栏立即移除该列表
- 列表被删除时，级联清理 share 和 members
- 创建 share 时若已存在 → 返回现有 share（一个列表只有一个 share）

## 后续迭代方向

- 成员退出列表 `DELETE /api/lists/:id/leave`
- 共享列表内显示编辑者头像/名称
- 活动日志（谁完成了什么任务）
- WebSocket 实时同步
- 工作空间模型（多列表打包共享）
