# SSH Client 사용자 가이드

이 문서는 SSH Client를 처음 사용하는 사용자를 위한 가이드입니다.

## 목차

- [빠른 시작](#빠른-시작)
- [프로파일 시스템](#프로파일-시스템)
- [인증 방법](#인증-방법)
- [실전 사용 시나리오](#실전-사용-시나리오)

## 빠른 시작

### 1. 도움말 확인

언제든지 도움말을 볼 수 있습니다:

```bash
./sshclient -h
./sshclient profile help
```

### 2. 첫 연결 시도

**방법 1: 전통적인 SSH 스타일** (가장 간단)

```bash
./sshclient user@hostname
```

**방법 2: 프로파일 생성 후 사용** (자주 사용할 서버)

```bash
# 프로파일 생성
./sshclient profile add myserver

# 프로파일로 접속
./sshclient @myserver
```

### 3. 원격 명령 실행

```bash
# 전통적인 방식
./sshclient user@hostname ls -la

# 프로파일 사용
./sshclient @myserver uptime
```

## 프로파일 시스템

프로파일을 사용하면 자주 접속하는 서버의 정보를 저장하여 빠르게 접속할 수 있습니다.

### 프로파일 생성

```bash
./sshclient profile add webserver
```

대화형으로 다음 정보를 입력합니다:
- **Host**: 호스트명 또는 IP 주소 (예: `web.example.com` 또는 `192.168.1.100`)
- **User**: SSH 사용자명 (예: `root`, `admin`)
- **Port**: SSH 포트 (기본값: 22)
- **인증 방법**:
  - `1` = SSH 키 (권장)
  - `2` = 비밀번호

**SSH 키 인증 예시**:
```
Creating new profile: webserver

Host (hostname or IP): web.example.com
User: deploy
Port (default: 22): 22

Authentication method (1: SSH key, 2: Password): 1
Path to SSH key (default: ~/.ssh/id_rsa): [Enter]
Using default key: /Users/username/.ssh/id_rsa

✓ Profile 'webserver' created successfully!
```

**비밀번호 인증 예시**:
```
Creating new profile: testserver

Host (hostname or IP): test.example.com
User: admin
Port (default: 22): 2222

Authentication method (1: SSH key, 2: Password): 2
Password (leave empty to prompt on connect): [비밀번호 입력]

🔐 Password will be encrypted using AES-256-GCM
✅ Password encrypted and will be stored securely

✓ Profile 'testserver' created successfully!
```

### 프로파일 관리

**목록 보기**:
```bash
./sshclient profile list
```

출력 예시:
```
📋 Custom Profiles (~/.sshclient/config.yaml):
──────────────────────────────────────────────────────────────────────
  @webserver       deploy@web.example.com:22 (key: ~/.ssh/id_rsa)
  @dbserver        admin@db.example.com:3306 (password)
  @testserver      root@192.168.1.100:22 (key: ~/.ssh/test_key)

🔧 SSH Config Profiles (~/.ssh/config):
──────────────────────────────────────────────────────────────────────
  @production      deploy@prod.example.com:22 (key: ~/.ssh/prod_key)
```

**프로파일 상세 정보**:
```bash
./sshclient profile show webserver
```

**프로파일 삭제**:
```bash
./sshclient profile remove webserver
```

### 프로파일 사용하기

**대화형 셸**:
```bash
./sshclient @webserver
```

**원격 명령 실행**:
```bash
./sshclient @webserver ls -la
./sshclient @webserver "df -h && uptime"
```

## 인증 방법

### 방법 1: SSH 키 인증 (권장) 🔑

SSH 키는 가장 안전하고 편리한 인증 방법입니다.

**기본 SSH 키 사용**:
- 프로그램이 자동으로 `~/.ssh/id_rsa`, `~/.ssh/id_ed25519`, `~/.ssh/id_ecdsa`를 찾습니다
- SSH 키 경로를 생략하면 자동으로 감지됩니다

```bash
./sshclient user@hostname  # 자동으로 기본 키 사용
```

**특정 SSH 키 지정**:
```bash
./sshclient user@hostname -key ~/.ssh/custom_key
```

**SSH 키가 없는 경우**:
```bash
# SSH 키 생성 (Ed25519 권장)
ssh-keygen -t ed25519 -C "your_email@example.com"

# 공개키를 서버에 복사
ssh-copy-id -i ~/.ssh/id_ed25519.pub user@hostname
```

### 방법 2: 접속 시 비밀번호 입력 (간단)

```bash
./sshclient user@hostname
# Password for user@hostname: [비밀번호 입력]
```

### 방법 3: 프로파일에 암호화하여 저장 (편리) 🔐

프로파일 생성 시 비밀번호를 입력하면 AES-256-GCM으로 자동 암호화되어 저장됩니다.

```bash
./sshclient profile add myserver
# Authentication method: 2 (Password)
# Password: [비밀번호 입력]
# 🔐 Password will be encrypted using AES-256-GCM
# ✅ Password encrypted and will be stored securely
```

이후 접속 시 비밀번호 입력 없이 자동으로 복호화되어 사용됩니다:

```bash
./sshclient @myserver
# Connected successfully!  (비밀번호 자동 복호화)
```

**보안**:
- 비밀번호는 **AES-256-GCM**으로 자동 암호화되어 저장
- **PBKDF2** (100,000 iterations)로 키 파생
- 설정 파일을 열어도 암호화된 문자열만 표시
- 마스터 비밀번호 입력 불필요 (자동 암호화/복호화)

### 방법 4: 명령줄 인자 (비권장) ⚠️

```bash
./sshclient -host example.com -user myuser -password "mypass" -i
```

**경고**: 명령 기록(history)에 비밀번호가 남으므로 절대 권장하지 않습니다.

## 실전 사용 시나리오

### 시나리오 1: 웹 서버 로그 모니터링

```bash
#!/bin/bash
# monitor-web.sh

# 웹 서버 접속 및 실시간 로그 확인
./sshclient @webserver "tail -f /var/log/nginx/access.log"
```

### 시나리오 2: 데이터베이스 백업

```bash
#!/bin/bash
# backup-db.sh

# DB 서버에서 덤프 생성 및 로컬 다운로드
BACKUP_FILE="db_backup_$(date +%Y%m%d_%H%M%S).sql"

./sshclient @dbserver "mysqldump -u root -p mydb > /tmp/$BACKUP_FILE"
./sshclient @dbserver "cat /tmp/$BACKUP_FILE" > "./backups/$BACKUP_FILE"

echo "Backup saved to ./backups/$BACKUP_FILE"
```

### 시나리오 3: 여러 서버 동시 상태 확인

```bash
#!/bin/bash
# check-servers.sh

SERVERS=("web1" "web2" "db1" "cache1")

echo "=== Server Status Check ==="
for server in "${SERVERS[@]}"; do
    echo ""
    echo "[$server]"
    ./sshclient @$server "echo 'Uptime:' && uptime && echo 'Disk:' && df -h /" 2>/dev/null
done
```

### 시나리오 4: 배포 자동화

```bash
#!/bin/bash
# deploy.sh

PROFILE="production"
APP_DIR="/var/www/myapp"

echo "🚀 Deploying to production..."

# 1. Git pull
./sshclient @$PROFILE "cd $APP_DIR && git pull origin main"

# 2. Dependencies
./sshclient @$PROFILE "cd $APP_DIR && npm install"

# 3. Build
./sshclient @$PROFILE "cd $APP_DIR && npm run build"

# 4. Restart service
./sshclient @$PROFILE "sudo systemctl restart myapp"

echo "✅ Deployment complete!"
```

### 시나리오 5: 포트 포워딩 (로컬 개발)

```bash
#!/bin/bash
# forward-db.sh

# 원격 DB를 로컬 포트로 포워딩 (SSH 터널링)
# 주의: 현재 버전은 포트 포워딩 미지원
# 일반 ssh 명령과 함께 사용:

ssh -L 3306:localhost:3306 user@dbserver
```

### 시나리오 6: 대량 서버 설정 변경

```bash
#!/bin/bash
# update-config.sh

CONFIG_LINE="MaxConnections=1000"
CONFIG_FILE="/etc/myapp/config.ini"

SERVERS=($(./sshclient profile list | grep '@' | awk '{print $1}' | tr -d '@'))

for server in "${SERVERS[@]}"; do
    echo "Updating $server..."
    ./sshclient @$server "echo '$CONFIG_LINE' >> $CONFIG_FILE"
done
```

## SSH Config 통합

기존 `~/.ssh/config` 파일의 Host 항목을 자동으로 읽습니다:

```
Host production
    HostName prod.example.com
    User deploy
    Port 22
    IdentityFile ~/.ssh/prod_key

Host staging
    HostName staging.example.com
    User admin
    Port 2222
    IdentityFile ~/.ssh/staging_key
```

사용: `./sshclient @production` 또는 `./sshclient @staging`

**장점**:
- 기존 SSH 설정을 그대로 활용
- 다른 SSH 도구와 설정 공유
- 별도 설정 불필요

**주의**: SSH config의 프로파일은 읽기 전용이며, `profile remove` 명령으로 삭제할 수 없습니다.

## 다음 단계

- 📖 [사용자 매뉴얼](User-Manual.md) - 상세한 명령어 레퍼런스
- 🔧 [문제 해결 가이드](User-Manual.md#문제-해결) - 일반적인 문제와 해결 방법
- 💡 [팁과 트릭](User-Manual.md#팁과-트릭) - 효율적인 사용법

---

**참고**: 추가 도움이 필요하면 `./sshclient -h` 또는 GitHub Issues를 참고하세요.
