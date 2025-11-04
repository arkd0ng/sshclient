# SSH Client

macOS용 터미널 SSH 클라이언트 - 순수 Go 구현

## 특징

✅ **프로파일 시스템** (v1.2.0 신규)
- `./sshclient @profile` - 저장된 프로파일로 빠른 접속
- 커스텀 프로파일: `~/.sshclient/config.yaml`
- SSH config 호환: `~/.ssh/config` 자동 읽기
- 대화형 프로파일 관리 (add, list, remove, show)

✅ **AES-256 비밀번호 암호화** (v1.2.0 신규)
- 프로파일에 저장되는 비밀번호를 AES-256-GCM으로 암호화
- 마스터 비밀번호로 보호
- 평문 비밀번호 저장 방지
- PBKDF2 키 파생으로 강력한 보안

✅ **전통적인 SSH 사용법 지원**
- `./sshclient user@host` - 익숙한 SSH 명령어 스타일
- `./sshclient user@host command` - 원격 명령 실행
- 플래그 방식도 여전히 지원

✅ **친절한 도움말 시스템**
- 상황별 맞춤 도움말과 사용 예시
- 메뉴얼 없이도 직관적으로 사용 가능
- `./sshclient -h` 또는 `./sshclient profile help`

✅ **완전히 독립적인 실행**
- OS의 `ssh` 명령어 사용 안 함
- OpenSSH 설치 불필요
- 단일 바이너리로 완전 동작

✅ **권한 불필요**
- root/sudo 권한 필요 없음
- 일반 사용자 권한으로 실행

✅ **순수 Go 구현**
- `golang.org/x/crypto/ssh` 사용
- 크로스 플랫폼 (macOS, Linux, Windows)

## 빌드

```bash
go build -o sshclient
```

## 빠른 시작

도움말이 필요하면 언제든지:
```bash
./sshclient -h              # 전체 도움말
./sshclient profile help    # 프로파일 도움말
```

## 사용법

### 🌟 프로파일 스타일 (v1.2.0 신규, 권장)

자주 접속하는 서버를 프로파일로 저장하여 빠르게 접속할 수 있습니다!

#### 1. 프로파일 생성 (대화형)

```bash
./sshclient profile add myserver
```

대화형으로 서버 정보를 입력합니다:
- Host (호스트명 또는 IP)
- User (SSH 사용자명)
- Port (기본값: 22)
- 인증 방식 (SSH 키 또는 비밀번호)

#### 2. 프로파일로 접속

```bash
# 대화형 셸
./sshclient @myserver

# 원격 명령 실행
./sshclient @myserver uptime
./sshclient @myserver "df -h"
./sshclient @myserver ls -la /var/log
```

#### 3. 프로파일 관리

```bash
# 모든 프로파일 목록 (커스텀 + SSH config)
./sshclient profile list

# 프로파일 상세 정보
./sshclient profile show myserver

# 프로파일 삭제
./sshclient profile remove myserver

# 도움말
./sshclient profile help
```

#### 4. SSH Config 호환

기존 `~/.ssh/config` 파일의 Host 항목을 자동으로 읽습니다:

```
# ~/.ssh/config
Host myserver
    HostName example.com
    User root
    Port 22
    IdentityFile ~/.ssh/id_rsa
```

사용: `./sshclient @myserver`

### 📌 전통적인 SSH 스타일

익숙한 SSH 명령어 스타일로 바로 접속할 수 있습니다!

#### 1. 대화형 셸 (기본)

```bash
./sshclient user@example.com
```

비밀번호 프롬프트가 나타나고 대화형 셸이 시작됩니다.

#### 2. 원격 명령 실행

```bash
./sshclient user@example.com ls -la
./sshclient user@example.com "df -h"
./sshclient user@example.com uptime
```

명령을 지정하면 실행 후 자동으로 종료됩니다.

#### 3. SSH 키 인증

```bash
./sshclient user@example.com -key ~/.ssh/id_rsa
./sshclient user@example.com -key ~/.ssh/id_ed25519 hostname
```

#### 4. 다른 포트 사용

```bash
./sshclient user@example.com -port 2222
```

### 📋 플래그 스타일 (기존 방식)

더 명시적인 옵션 지정이 필요한 경우 사용합니다.

#### 1. 비밀번호 인증 (대화형 셸)

```bash
./sshclient -host example.com -user myuser -i
```

#### 2. SSH 키 인증 (대화형 셸)

```bash
./sshclient -host example.com -user myuser -key ~/.ssh/id_rsa -i
```

#### 3. 원격 명령 실행

```bash
./sshclient -host example.com -user myuser -cmd "ls -la"
```

#### 4. 비밀번호를 명령줄에서 지정 (권장하지 않음)

```bash
./sshclient -host example.com -user myuser -password "mypassword" -i
```

**보안상 권장하지 않습니다.** 프롬프트로 입력하거나 SSH 키를 사용하세요.

## 옵션

### 위치 인자 (Positional Arguments)

| 형식 | 설명 | 예제 |
|------|------|------|
| `user@host` | SSH 접속 대상 (user@host 형식) | `root@example.com` |
| `[command...]` | 실행할 원격 명령 (선택사항) | `ls -la` |

### 플래그 옵션

| 옵션 | 설명 | 기본값 |
|------|------|--------|
| `-host` | SSH 서버 호스트명 또는 IP 주소 | - |
| `-user` | SSH 사용자명 | - |
| `-port` | SSH 서버 포트 | 22 |
| `-password` | 비밀번호 (보안상 비추천) | - |
| `-key` | SSH 개인키 파일 경로 | - |
| `-cmd` | 실행할 원격 명령 (플래그 스타일 사용 시) | - |
| `-i` | 대화형 셸 세션 시작 (플래그 스타일 사용 시) | false |
| `-version` | 버전 정보 출력 | - |

**참고**: `user@host` 형식 사용 시 `-i` 플래그 없이도 대화형 모드가 기본값입니다.

## 프로파일 설정 파일 상세

### 커스텀 프로파일 (YAML 형식)

위치: `~/.sshclient/config.yaml`

```yaml
profiles:
  myserver:
    host: example.com         # 필수: 호스트명 또는 IP
    user: root                # 필수: SSH 사용자명
    port: "22"                # 선택: 포트 (기본값: 22)
    key: /path/to/key         # 선택: SSH 개인키 경로
    password: "secret"        # 선택: 비밀번호 (비권장, 평문 저장)

  # 예시 1: SSH 키 인증 (권장)
  webserver:
    host: web.example.com
    user: deploy
    port: "22"
    key: /Users/shlee/.ssh/deploy_key

  # 예시 2: 기본 SSH 키 자동 사용 (키 경로 미지정)
  dbserver:
    host: db.example.com
    user: admin
    port: "3306"
    # key를 지정하지 않으면 ~/.ssh/id_rsa 등을 자동 탐색

  # 예시 3: 다른 포트 사용
  jumphost:
    host: jump.example.com
    user: root
    port: "2222"
    key: /Users/shlee/.ssh/id_rsa
```

### SSH Config 호환

기존 `~/.ssh/config` 파일도 자동으로 읽습니다:

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

**주의**: SSH config의 프로파일은 읽기 전용이며, `profile remove` 명령으로 삭제할 수 없습니다.

## 인증 방법 상세 가이드

### 1. SSH 키 인증 (권장)

**장점**: 안전하고 편리한 인증 방식

#### 방법 A: 기본 SSH 키 자동 사용

프로그램은 자동으로 다음 경로의 키를 찾습니다:
1. `~/.ssh/id_rsa`
2. `~/.ssh/id_ed25519`
3. `~/.ssh/id_ecdsa`

프로파일이나 명령줄에서 키를 지정하지 않으면 자동으로 위 키들을 시도합니다.

```bash
# 프로파일에 키 미지정 시 자동 탐색
./sshclient @myserver

# 출력:
# Trying default SSH key: /Users/shlee/.ssh/id_rsa
# Using key authentication with /Users/shlee/.ssh/id_rsa
# Connected successfully!
```

#### 방법 B: 특정 SSH 키 지정

```bash
# 명령줄에서 지정
./sshclient user@example.com -key ~/.ssh/custom_key

# 프로파일에 지정
./sshclient profile add myserver
# ... (대화형 입력)
# Authentication method: 1 (SSH key)
# Path to SSH key: /Users/shlee/.ssh/custom_key
```

#### SSH 키 생성 방법

아직 SSH 키가 없다면:

```bash
# RSA 키 생성 (호환성 좋음)
ssh-keygen -t rsa -b 4096 -C "your_email@example.com"

# Ed25519 키 생성 (최신, 더 안전)
ssh-keygen -t ed25519 -C "your_email@example.com"

# 생성된 키를 서버에 복사
ssh-copy-id user@example.com
```

### 2. 비밀번호 인증

#### 방법 A: 접속 시 프롬프트 (가장 안전)

```bash
./sshclient @myserver
# Password for user@host: [비밀번호 입력]
```

#### 방법 B: 프로파일에 암호화하여 저장 (권장) 🔐

```bash
./sshclient profile add myserver
# Authentication method: 2 (Password)
# Password: [비밀번호 입력]
# 🔐 Password will be encrypted using AES-256-GCM
# ✅ Password encrypted and will be stored securely
```

**보안**:
- 비밀번호는 **AES-256-GCM**으로 자동 암호화되어 저장
- **PBKDF2** (100,000 iterations)로 키 파생
- 설정 파일을 열어도 암호화된 문자열만 표시
- 마스터 비밀번호 입력 불필요 (자동 암호화/복호화)

**사용**:
```bash
./sshclient @myserver
# Connected successfully!  (비밀번호 자동 복호화)
```

#### 방법 C: 명령줄 인자 (비권장)

```bash
./sshclient -host example.com -user myuser -password "mypass" -i
```

**경고**: 명령 기록(history)에 비밀번호가 남으므로 절대 권장하지 않습니다.

## 주요 기능

### 1. 인증 방식

- **비밀번호 인증**: 프롬프트 또는 명령줄 옵션
- **SSH 키 인증**: RSA, Ed25519, ECDSA 키 지원

### 2. 실행 모드

- **대화형 셸** (`-i`): 일반 SSH처럼 대화형 터미널
- **명령 실행** (`-cmd`): 단일 명령 실행 후 종료

### 3. 파일 전송 (코드에 구현됨)

client.go에 SCP 기능이 포함되어 있습니다:
- `CopyFile()`: 로컬 → 원격 파일 전송
- `DownloadFile()`: 원격 → 로컬 파일 전송

## 동작 원리

이 SSH 클라이언트는 OS의 ssh 명령어를 호출하지 않고, Go의 SSH 라이브러리를 사용하여 직접 SSH 프로토콜을 구현합니다:

1. **네트워크 연결**: `net.Dial()`로 TCP 소켓 직접 생성
2. **SSH 프로토콜**: `golang.org/x/crypto/ssh`로 SSH 핸드셰이크 및 암호화 통신
3. **터미널 제어**: `golang.org/x/crypto/ssh/terminal`로 터미널 raw 모드 설정
4. **인증**: SSH 키 파일을 직접 파싱하여 인증 처리

## 실제 사용 예제

### 예제 1: 기본 접속 (비밀번호)

```bash
./sshclient root@sun.neteer.co.kr
```

실행하면:
```
Password for root@sun.neteer.co.kr: [비밀번호 입력]
Connected successfully!
Starting interactive shell... (Press Ctrl+D or type 'exit' to quit)
[root@sun ~]#
```

### 예제 2: 원격 명령 실행

```bash
./sshclient admin@192.168.1.100 uptime
./sshclient admin@192.168.1.100 "df -h"
./sshclient admin@192.168.1.100 ls -la /var/log
```

### 예제 3: SSH 키로 접속

```bash
./sshclient deploy@myserver.com -key ~/.ssh/deploy_key
./sshclient deploy@myserver.com -key ~/.ssh/id_rsa systemctl status nginx
```

### 예제 4: 다른 포트 사용

```bash
./sshclient user@example.com -port 2222
./sshclient user@example.com -port 2222 hostname
```

### 예제 5: 플래그 스타일 (기존 방식)

```bash
./sshclient -host 192.168.1.100 -user admin -i
./sshclient -host myserver.com -user deploy -key ~/.ssh/deploy_key -cmd "uptime"
```

## 실제 사용 시나리오

### 시나리오 1: 여러 서버를 자주 관리하는 경우

**상황**: 웹서버, DB서버, 백업서버를 자주 접속해야 함

**해결책**: 프로파일로 관리

```bash
# 초기 설정 (한 번만 실행)
./sshclient profile add web
./sshclient profile add db
./sshclient profile add backup

# 이후 접속 (간편!)
./sshclient @web
./sshclient @db "mysqldump -u root mydb > /backup/db.sql"
./sshclient @backup ls -lh /backup
```

### 시나리오 2: 일회성 서버 접속

**상황**: 처음 접속하거나 한 번만 접속할 서버

**해결책**: 전통적인 SSH 스타일 사용

```bash
# 프로파일 만들 필요 없이 바로 접속
./sshclient admin@temp-server.example.com

# 원격 명령 바로 실행
./sshclient admin@temp-server.example.com "cat /var/log/syslog | tail -50"
```

### 시나리오 3: 배포 스크립트에서 사용

**상황**: 자동화 스크립트에서 SSH 명령 실행

**해결책**: 프로파일 + 명령 실행

```bash
#!/bin/bash
# deploy.sh

echo "Deploying to production..."

# 최신 코드 pull
./sshclient @production "cd /var/www && git pull origin main"

# 의존성 설치
./sshclient @production "cd /var/www && npm install"

# 서비스 재시작
./sshclient @production "systemctl restart nginx"

echo "Deployment complete!"
```

### 시나리오 4: 서버 상태 모니터링

**상황**: 여러 서버의 상태를 빠르게 확인

**해결책**: 반복문으로 여러 서버 확인

```bash
#!/bin/bash
# check-servers.sh

for server in web db cache backup
do
  echo "=== Checking $server ==="
  ./sshclient @$server "uptime && df -h / && free -h"
  echo ""
done
```

### 시나리오 5: 점프 호스트를 통한 접속

**상황**: 보안상 점프 호스트(bastion)를 거쳐야 하는 경우

**해결책**: 프로파일 + SSH config 조합

`~/.ssh/config`:
```
Host jumphost
    HostName jump.example.com
    User admin
    IdentityFile ~/.ssh/jump_key

Host internal-server
    HostName 10.0.1.100
    User root
    ProxyJump jumphost
    IdentityFile ~/.ssh/internal_key
```

사용:
```bash
./sshclient @internal-server
```

### 시나리오 6: 다중 서버 로그 수집

**상황**: 여러 서버의 로그를 한 번에 수집

**해결책**: 스크립트로 자동화

```bash
#!/bin/bash
# collect-logs.sh

TIMESTAMP=$(date +%Y%m%d_%H%M%S)
LOGDIR="./logs_$TIMESTAMP"
mkdir -p "$LOGDIR"

for server in web1 web2 web3 db
do
  echo "Collecting logs from $server..."
  ./sshclient @$server "tail -1000 /var/log/app.log" > "$LOGDIR/$server.log"
done

echo "Logs collected in $LOGDIR"
```

## FAQ (자주 묻는 질문)

### Q1: 비밀번호를 매번 입력하기 귀찮습니다

**A**: SSH 키 인증을 사용하세요. 한 번 설정하면 비밀번호 없이 접속할 수 있습니다.

```bash
# 1. SSH 키 생성 (아직 없는 경우)
ssh-keygen -t rsa -b 4096

# 2. 서버에 공개키 복사
ssh-copy-id user@server

# 3. 이제 비밀번호 없이 접속 가능
./sshclient user@server
```

### Q2: 여러 서버를 관리할 때 프로파일과 SSH config 중 뭘 써야 하나요?

**A**: 둘 다 장단점이 있습니다:

| 방식 | 장점 | 단점 |
|------|------|------|
| **커스텀 프로파일** (`~/.sshclient/config.yaml`) | - 이 프로그램 전용<br>- 대화형으로 쉽게 추가<br>- `profile` 명령으로 관리 편리 | - 다른 SSH 도구와 공유 안 됨 |
| **SSH config** (`~/.ssh/config`) | - 모든 SSH 도구와 공유<br>- 표준 형식<br>- 기존 설정 재사용 | - 수동으로 파일 편집 필요<br>- 읽기 전용 (삭제 불가) |

**권장**:
- 일반적인 경우 → SSH config 사용 (다른 도구와 호환)
- 이 프로그램만 사용 → 커스텀 프로파일 (관리 편리)

### Q3: 프로파일 이름과 SSH config의 Host 이름이 같으면 어떻게 되나요?

**A**: 커스텀 프로파일이 우선순위가 높습니다.

```
우선순위: ~/.sshclient/config.yaml > ~/.ssh/config
```

### Q4: 기본 SSH 키를 자동으로 못 찾는 것 같아요

**A**: 다음을 확인하세요:

```bash
# 1. SSH 키 존재 확인
ls -la ~/.ssh/id_*

# 2. 키 권한 확인 (600이어야 함)
chmod 600 ~/.ssh/id_rsa

# 3. 공개키가 서버에 등록되어 있는지 확인
ssh-copy-id user@server

# 4. 디버깅: 키 경로 명시적으로 지정해보기
./sshclient user@server -key ~/.ssh/id_rsa
```

### Q5: 프로파일 설정 파일 위치를 변경할 수 있나요?

**A**: 현재는 `~/.sshclient/config.yaml`로 고정되어 있습니다. 환경 변수로 변경하는 기능은 향후 추가 예정입니다.

### Q6: OpenSSH와 비교해서 어떤 차이가 있나요?

**A**: 주요 차이점:

| 항목 | sshclient | OpenSSH |
|------|-----------|---------|
| **설치** | 단일 바이너리, 설치 불필요 | OS 패키지 관리자로 설치 |
| **의존성** | 없음 (모두 포함) | 시스템 라이브러리 의존 |
| **권한** | 일반 사용자 | 일반 사용자 |
| **프로파일** | 대화형 관리 (`profile add`) | 수동 파일 편집 |
| **호스트 키 검증** | ⚠️ 현재 비활성화 (개발 중) | ✅ 완전 지원 |
| **포트 포워딩** | ❌ 미지원 | ✅ 지원 |
| **ProxyJump** | ❌ 미지원 | ✅ 지원 |

**용도**:
- **sshclient**: 개발/테스트 환경, 간편한 서버 관리
- **OpenSSH**: 프로덕션, 고급 기능 필요 시

### Q7: Windows나 Linux에서도 사용할 수 있나요?

**A**: 네! Go로 작성되어 크로스 플랫폼을 지원합니다.

```bash
# Linux용 빌드
GOOS=linux GOARCH=amd64 go build -o sshclient-linux

# Windows용 빌드
GOOS=windows GOARCH=amd64 go build -o sshclient.exe

# macOS용 빌드 (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o sshclient-macos-arm64
```

## 트러블슈팅

### 문제: "Failed to read password: operation not supported on socket"

**원인**: 비대화형 환경(파이프, 리다이렉션)에서 비밀번호 입력 시도

**해결**:
```bash
# ❌ 작동 안 함
echo "ls" | ./sshclient user@server

# ✅ 해결 방법 1: SSH 키 사용
./sshclient user@server -key ~/.ssh/id_rsa ls

# ✅ 해결 방법 2: 프로파일에 키 등록
./sshclient profile add server  # 키 경로 지정
./sshclient @server ls
```

### 문제: "Failed to parse private key"

**원인**: 잘못된 SSH 키 파일 또는 암호화된 키

**해결**:
```bash
# 1. 키 파일 형식 확인
head -1 ~/.ssh/id_rsa
# 출력: -----BEGIN RSA PRIVATE KEY----- 또는
#      -----BEGIN OPENSSH PRIVATE KEY-----

# 2. 암호화되지 않은 키 생성
ssh-keygen -t rsa -b 4096 -N ""

# 3. 기존 암호화된 키의 암호 제거
ssh-keygen -p -f ~/.ssh/id_rsa -N ""
```

### 문제: "Connection refused" 또는 "No route to host"

**원인**: 네트워크 문제 또는 잘못된 호스트/포트

**해결**:
```bash
# 1. 호스트 이름 확인
ping example.com

# 2. 포트 확인 (SSH는 기본 22번)
nc -zv example.com 22

# 3. 방화벽 확인
# 서버에서: sudo ufw status
# 서버에서: sudo firewall-cmd --list-all

# 4. SSH 서비스 상태 확인
# 서버에서: systemctl status sshd
```

### 문제: 프로파일이 목록에 안 나타남

**원인**: YAML 형식 오류 또는 파일 권한 문제

**해결**:
```bash
# 1. 설정 파일 확인
cat ~/.sshclient/config.yaml

# 2. YAML 형식 검증 (온라인: yamllint.com)

# 3. 파일 권한 확인 및 수정
ls -la ~/.sshclient/config.yaml
chmod 600 ~/.sshclient/config.yaml

# 4. 설정 파일 재생성
rm ~/.sshclient/config.yaml
./sshclient profile add myserver
```

### 문제: "Permission denied (publickey)"

**원인**: SSH 키가 서버에 등록되지 않음

**해결**:
```bash
# 1. 공개키를 서버에 복사
ssh-copy-id -i ~/.ssh/id_rsa.pub user@server

# 또는 수동으로:
# 2. 공개키 내용 확인
cat ~/.ssh/id_rsa.pub

# 3. 서버의 ~/.ssh/authorized_keys에 추가
# (서버에서): echo "복사한_공개키" >> ~/.ssh/authorized_keys

# 4. 권한 설정 (서버에서)
chmod 700 ~/.ssh
chmod 600 ~/.ssh/authorized_keys
```

### 문제: 터미널이 깨지거나 이상하게 동작

**원인**: 터미널 상태가 제대로 복원되지 않음

**해결**:
```bash
# 터미널 리셋
reset

# 또는
stty sane
```

### 문제: 프로파일은 있는데 "profile not found" 에러

**원인**: 대소문자 불일치 또는 공백 문제

**해결**:
```bash
# 1. 정확한 프로파일 이름 확인
./sshclient profile list

# 2. 대소문자 정확히 일치시켜 사용
./sshclient @MyServer  # ❌
./sshclient @myserver  # ✅
```

### 문제: "Error decrypting password" 에러 (v1.2.0 이전 버전에서 업그레이드한 경우)

**원인**: v1.2.0 이전 버전에서 마스터 비밀번호로 암호화된 비밀번호는 새 버전과 호환되지 않음

**해결**:
```bash
# 1. 영향받는 프로파일 제거
./sshclient profile remove old-profile

# 2. 프로파일 재생성 (새로운 자동 암호화 사용)
./sshclient profile add old-profile
# 비밀번호 재입력 - 이제 마스터 비밀번호 없이 자동 암호화됨

# 참고: v1.2.0부터는 마스터 비밀번호 없이 자동 암호화/복호화됩니다
```

## 팁과 트릭

### 💡 Tip 1: 별칭(alias)으로 더 간편하게

`~/.bashrc` 또는 `~/.zshrc`에 추가:

```bash
alias sssh='./sshclient'
alias ssh-prod='./sshclient @production'
alias ssh-dev='./sshclient @development'
alias ssh-db='./sshclient @database'
```

사용:
```bash
sssh @myserver
ssh-prod
ssh-db "show databases"
```

### 💡 Tip 2: 프로파일 백업

```bash
# 백업
cp ~/.sshclient/config.yaml ~/.sshclient/config.yaml.backup

# 다른 머신으로 복사
scp ~/.sshclient/config.yaml user@newmachine:~/.sshclient/
```

### 💡 Tip 3: 기존 SSH config 활용

이미 `~/.ssh/config`에 많은 서버가 설정되어 있다면:

```bash
# 모든 SSH config 호스트 확인
./sshclient profile list

# SSH config의 호스트를 그대로 사용
./sshclient @production
./sshclient @staging
```

### 💡 Tip 4: 서버 그룹별 관리

프로파일 이름에 접두사 사용:

```yaml
profiles:
  prod-web:
    host: web.prod.example.com
  prod-db:
    host: db.prod.example.com
  dev-web:
    host: web.dev.example.com
  dev-db:
    host: db.dev.example.com
```

사용:
```bash
./sshclient @prod-web
./sshclient @dev-db
```

### 💡 Tip 5: 원격 명령 체이닝

```bash
# 여러 명령을 && 로 연결
./sshclient @web "cd /var/www && git pull && npm install && pm2 restart all"

# 조건부 실행
./sshclient @db "mysqldump mydb > backup.sql && echo 'Backup successful' || echo 'Backup failed'"
```

## 파일 구조

```
sshclient/
├── main.go          # CLI 인터페이스 및 메인 로직
├── client.go        # SSH 클라이언트 핵심 구현
├── config.go        # 프로파일 설정 관리 (YAML + SSH config)
├── profile.go       # 프로파일 관리 명령어
├── crypto.go        # AES-256 암호화/복호화 (v1.2.0)
├── go.mod           # Go 모듈 정의
├── go.sum           # 의존성 체크섬
├── sshclient        # 컴파일된 바이너리
├── README.md        # 사용자 문서
└── CLAUDE.md        # 개발 가이드
```

## 의존성

- `golang.org/x/crypto/ssh` - SSH 프로토콜 구현
- `golang.org/x/term` - 터미널 제어
- `gopkg.in/yaml.v3` - YAML 설정 파일 파싱

## 빌드 정보

- **컴파일된 바이너리**: `sshclient`
- **크기**: 약 6.1MB
- **플랫폼**: macOS (x86_64)
- **Go 버전**: 최신 안정 버전

## 보안 주의사항

### ✅ 보안 기능 (v1.2.0)

**AES-256 비밀번호 암호화**:
- 프로파일에 저장되는 비밀번호를 **AES-256-GCM**으로 자동 암호화
- **PBKDF2** (100,000 iterations)로 키 파생
- 암호화된 데이터만 디스크에 저장
- 마스터 비밀번호 입력 불필요 (자동 처리)

**설정 파일 형식 예시**:
```yaml
myserver:
    host: example.com
    user: root
    encrypted_password: "kR7v+2Tm...암호화된데이터...==" # AES-256-GCM 암호화
```

### ⚠️  주의사항

**테스트/개발용**:
- `ssh.InsecureIgnoreHostKey()` 사용 중 (호스트 키 검증 안 함)
- 프로덕션 환경에서는 호스트 키 검증 구현 필요

**암호화**:
- 비밀번호는 자동으로 암호화/복호화됩니다
- 프로그램 내부 키를 사용하여 기본적인 보안 제공
- 추가 보안이 필요한 경우 SSH 키 인증 사용 권장

### 🔒 프로파일 파일 권한

보안을 위해 설정 파일이 적절한 권한으로 자동 생성됩니다:
- `~/.sshclient/`: 0700 (소유자만 접근)
- `~/.sshclient/config.yaml`: 0600 (소유자만 읽기/쓰기)
- `~/.sshclient/master.hash`: 0600 (마스터 비밀번호 해시)

### 📝 프로덕션용 개선 사항

1. ✅ **비밀번호 암호화**: AES-256-GCM 구현 완료 (v1.2.0)
2. **호스트 키 검증**: `known_hosts` 파일 사용 (향후 추가)
3. **연결 타임아웃**: 환경에 맞게 조정
4. **에러 처리**: 더 자세한 에러 메시지

## 라이선스

테스트/학습용 프로젝트

## 버전

- **v1.2.0** (현재) - 프로파일 시스템, AES 암호화, 친절한 도움말
  - YAML 기반 프로파일 관리 시스템
  - SSH config (`~/.ssh/config`) 자동 읽기 지원
  - 프로파일 관리 명령어 (add, list, remove, show)
  - `@profile` 형식으로 빠른 접속
  - **AES-256-GCM 비밀번호 암호화** 🔐
  - PBKDF2 키 파생 (100,000 iterations)
  - 마스터 비밀번호 기반 보안
  - 기본 SSH 키 자동 감지
  - 상황별 맞춤 도움말과 사용 예시
  - 메뉴얼 없이도 사용 가능한 친절한 UX

- **v1.1.0** - 전통적인 SSH 스타일 지원
  - `user@host` 형식 파싱
  - 명령줄 인자로 원격 명령 실행 지원
  - 기존 플래그 방식과 호환성 유지

- **v1.0.0** - 초기 릴리스
  - 기본 SSH 클라이언트 기능
  - 비밀번호 및 SSH 키 인증
  - 대화형 셸 및 명령 실행
  - SCP 파일 전송 기능 (코드에 포함)
