#!/bin/bash

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

REPOS=(
    "frontend:https://github.com/lshy1993/genshin-quiz.git"
    "backend:https://github.com/lshy1993/genshin-quiz-backend.git"
    "vue-home:https://github.com/lshy1993/vue-home.git"
    "MoeLink:https://github.com/lshy1993/MoeLink.git"
    "AlphaSoul_Js:https://github.com/lshy1993/AlphaSoul_Js.git"
    "node-api:https://github.com/lshy1993/node-api.git"
    "liantui-hp:https://github.com/lshy1993/liantui-hp.git"
)

clone_or_update() {
    local name=$1
    local url=$2
    
    if [ -d "$name" ]; then
        print_status "更新项目: $name"
        cd "$name"
        if [ -d ".git" ]; then
            # 获取默认分支名称
            default_branch=$(git symbolic-ref refs/remotes/origin/HEAD 2>/dev/null | sed 's@^refs/remotes/origin/@@')
            if [ -z "$default_branch" ]; then
                # 如果无法获取默认分支，尝试从远程获取
                git remote set-head origin -a 2>/dev/null
                default_branch=$(git symbolic-ref refs/remotes/origin/HEAD 2>/dev/null | sed 's@^refs/remotes/origin/@@')
            fi
            
            # 如果还是无法获取，尝试常见的分支名称
            if [ -z "$default_branch" ]; then
                if git show-ref --verify --quiet refs/remotes/origin/main; then
                    default_branch="main"
                elif git show-ref --verify --quiet refs/remotes/origin/master; then
                    default_branch="master"
                else
                    print_warning "无法确定 $name 的默认分支，尝试使用 main"
                    default_branch="main"
                fi
            fi
            
            print_status "拉取 $name 的 $default_branch 分支"
            git pull origin "$default_branch"
        else
            print_warning "$name 目录存在但不是 Git 仓库，跳过更新"
        fi
        cd ..
    else
        print_status "克隆项目: $name"
        git clone "$url" "$name"
    fi
}

main() {
    print_status "🚀 开始设置 GenshinQuiz 开发环境..."
    
    # 确保 repos 目录存在
    mkdir -p repos
    cd repos
    
    for repo in "${REPOS[@]}"; do
        IFS=':' read -r name url <<< "$repo"
        clone_or_update "$name" "$url"
    done
    
    cd ..
    print_status "✅ 所有仓库设置完成！"
}

# 运行主函数
main "$@"