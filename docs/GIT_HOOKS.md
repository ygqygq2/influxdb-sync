# Git Hooks 使用说明

## 🎯 功能说明

本项目配置了 Git pre-commit hook，在每次 `git commit` 前自动执行以下检查：

1. **📝 自动格式化**: 运行 `task fmt` 格式化 Go 代码
2. **🔎 格式检查**: 运行 `task fmt-check` 确保代码格式正确
3. **🔍 静态分析**: 运行 `task vet` 进行代码静态分析
4. **📦 依赖检查**: 运行 `task deps-check` 检查依赖状态

## 🚀 使用方法

### 自动安装 (推荐)

```bash
# pre-commit hook 已经存在，只需确认权限
task test-hooks

# 如果需要重新安装权限
task install-hooks
```

### 手动验证

```bash
# 测试 hook 是否正常工作
task test-hooks

# 查看可用的 hook 管理命令
task --list | grep hooks
```

## 📋 Hook 执行流程

当你运行 `git commit` 时：

```
git commit -m "你的提交信息"
    ↓
🔍 执行 pre-commit 检查...
    ↓
📝 格式化 Go 代码... (task fmt)
    ↓
🔎 检查代码格式... (task fmt-check)
    ↓
🔍 运行静态分析... (task vet)
    ↓
📦 检查依赖... (task deps-check)
    ↓
✅ pre-commit 检查通过！
    ↓
提交成功 ✅
```

## 🛠️ 管理命令

```bash
# 测试 hooks 状态
task test-hooks

# 重新安装 hooks 权限
task install-hooks

# 卸载 hooks (如果需要)
task uninstall-hooks
```

## 🚨 如果检查失败

如果 pre-commit 检查失败，会看到类似错误信息：

```
❌ 代码格式检查失败，请运行 task fmt
```

解决方法：

1. 按照提示运行相应的修复命令
2. 重新添加修改的文件到暂存区：`git add .`
3. 重新提交：`git commit -m "你的提交信息"`

## 💡 临时跳过 hooks

如果遇到紧急情况需要跳过检查：

```bash
git commit --no-verify -m "紧急提交 - 跳过检查"
```

> ⚠️ **注意**: 不建议经常使用 `--no-verify`，这会降低代码质量

## 🎯 优势

✅ **自动化**: 无需手动记住运行 `task fmt`  
✅ **一致性**: 确保所有提交都经过格式化和检查  
✅ **质量保证**: 避免提交格式不规范的代码  
✅ **团队协作**: 统一的代码质量标准
