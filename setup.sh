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
    "frontend:https://github.com/lshy1993/genshin-quiz-frontend.git"
    "backend:https://github.com/lshy1993/genshin-quiz-backend.git"
    "vue-home:https://github.com/lshy1993/vue-home.git"
    "MoeLink:https://github.com/lshy1993/MoeLink.git"
    "AlphaSoul_Js:https://github.com/lshy1993/AlphaSoul_Js.git"
    "node-api:https://github.com/lshy1993/node-api.git"
)

clone_or_update() {
    local name=$1
    local url=$2
    
    if [ -d "$name" ]; then
        print_status "更新项目: $name"
        cd "$name"
        if [ -d ".git" ]; then
            git pull origin main || git pull origin master
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
    
    for repo in "${REPOS[@]}"; do
        IFS=':' read -r name url <<< "$repo"
        clone_or_update "$name" "$url"
    done
    
    print_status "✅ 所有仓库设置完成！"
}

# 运行主函数
main "$@"