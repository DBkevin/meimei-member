#!/usr/bin/env bash

set -Eeuo pipefail

APP_ROOT="/opt/meimei-member"
RELEASE_FILE="${1:-${APP_ROOT}/release.tar.gz}"
SERVER_DIR="${APP_ROOT}/server"
WEB_DIR="${APP_ROOT}/web"
BACKUP_ROOT="${APP_ROOT}/backups"
RELEASE_ROOT="${APP_ROOT}/releases"
ENV_FILE="${SERVER_DIR}/.env"
CONFIG_FILE="${SERVER_DIR}/config.yaml"
SYSTEMD_SERVICE_FILE="/etc/systemd/system/member-api.service"
NGINX_CONF_FILE="/etc/nginx/conf.d/meimei-member.conf"
TIMESTAMP="$(date +%Y%m%d%H%M%S)"
CURRENT_RELEASE_DIR="${RELEASE_ROOT}/${TIMESTAMP}"
CURRENT_BACKUP_DIR="${BACKUP_ROOT}/${TIMESTAMP}"

log() {
  printf '[deploy][%s] %s\n' "$(date '+%F %T')" "$*"
}

ensure_safe_path() {
  case "$1" in
    "${APP_ROOT}"|${APP_ROOT}/*) ;;
    *)
      log "检测到不安全路径: $1"
      exit 1
      ;;
  esac
}

cleanup_failed_release() {
  if [[ -d "${CURRENT_RELEASE_DIR}" ]]; then
    rm -rf "${CURRENT_RELEASE_DIR}"
  fi
}

on_error() {
  local exit_code=$?
  log "部署失败，已停止执行。exit code=${exit_code}"
  cleanup_failed_release
  exit "${exit_code}"
}

trap on_error ERR

main() {
  log "开始部署 member-api"

  ensure_safe_path "${SERVER_DIR}"
  ensure_safe_path "${WEB_DIR}"
  ensure_safe_path "${BACKUP_ROOT}"
  ensure_safe_path "${RELEASE_ROOT}"
  ensure_safe_path "${CURRENT_RELEASE_DIR}"
  ensure_safe_path "${CURRENT_BACKUP_DIR}"
  ensure_safe_path "${WEB_DIR}/dist"

  if [[ ! -f "${RELEASE_FILE}" ]]; then
    log "未找到发布包: ${RELEASE_FILE}"
    exit 1
  fi

  if [[ ! -f "${ENV_FILE}" ]]; then
    log "缺少运行时环境文件: ${ENV_FILE}"
    log "请先参考 server/.env.example 手动创建 ${ENV_FILE}"
    exit 1
  fi

  if [[ ! -f "${CONFIG_FILE}" ]]; then
    log "缺少生产配置文件: ${CONFIG_FILE}"
    log "请先参考 server/config.prod.example.yaml 手动创建 ${CONFIG_FILE}，并完成数据库配置"
    exit 1
  fi

  if [[ ! -f "${SYSTEMD_SERVICE_FILE}" ]]; then
    log "未检测到 systemd 服务文件: ${SYSTEMD_SERVICE_FILE}"
    log "请先安装 deploy/member-api.service，并执行 systemctl daemon-reload && systemctl enable member-api"
    exit 1
  fi

  mkdir -p \
    "${SERVER_DIR}" \
    "${WEB_DIR}" \
    "${BACKUP_ROOT}" \
    "${RELEASE_ROOT}" \
    "${CURRENT_RELEASE_DIR}" \
    "${CURRENT_BACKUP_DIR}"

  log "解压 release.tar.gz"
  tar -xzf "${RELEASE_FILE}" -C "${CURRENT_RELEASE_DIR}"

  if [[ ! -f "${CURRENT_RELEASE_DIR}/server/member-api" ]]; then
    log "发布包缺少 server/member-api"
    exit 1
  fi

  if [[ ! -d "${CURRENT_RELEASE_DIR}/web/dist" ]]; then
    log "发布包缺少 web/dist"
    exit 1
  fi

  if [[ -f "${SERVER_DIR}/member-api" ]]; then
    log "备份旧后端二进制"
    mkdir -p "${CURRENT_BACKUP_DIR}/server"
    cp -a "${SERVER_DIR}/member-api" "${CURRENT_BACKUP_DIR}/server/member-api"
  fi

  if [[ -d "${WEB_DIR}/dist" ]]; then
    log "备份旧前端 dist"
    mkdir -p "${CURRENT_BACKUP_DIR}/web"
    cp -a "${WEB_DIR}/dist" "${CURRENT_BACKUP_DIR}/web/dist"
  fi

  log "替换 /opt/meimei-member/server/member-api"
  install -m 755 "${CURRENT_RELEASE_DIR}/server/member-api" "${SERVER_DIR}/member-api"

  log "替换 /opt/meimei-member/web/dist"
  rm -rf "${WEB_DIR}/dist"
  cp -a "${CURRENT_RELEASE_DIR}/web/dist" "${WEB_DIR}/dist"

  log "重启 member-api"
  sudo systemctl restart member-api
  sudo systemctl is-active --quiet member-api

  if [[ -f "${NGINX_CONF_FILE}" ]]; then
    log "检查 Nginx 配置"
    sudo nginx -t

    log "重新加载 Nginx"
    sudo systemctl reload nginx
  else
    log "未检测到 Nginx 配置: ${NGINX_CONF_FILE}"
    log "已跳过 Nginx 校验与重载，请先手动安装 deploy/nginx.conf.example 并替换真实域名"
  fi

  log "部署完成"
  log "备份目录: ${CURRENT_BACKUP_DIR}"
}

main "$@"
