# SSH Client 사용자 매뉴얼

SSH Client의 모든 기능과 옵션에 대한 상세한 레퍼런스 문서입니다.

## 목차

- [명령어 레퍼런스](#명령어-레퍼런스)
- [프로파일 설정 파일](#프로파일-설정-파일)
- [플래그 참조](#플래그-참조)
- [FAQ](#faq)
- [문제 해결](#문제-해결)
- [팁과 트릭](#팁과-트릭)
- [아키텍처](#아키텍처)
- [보안 주의사항](#보안-주의사항)

## 명령어 레퍼런스

### 프로파일 관리 명령어

#### `profile add <name>`

새로운 프로파일을 대화형으로 생성합니다.

```bash
./sshclient profile add myserver
```

**입력 항목**:
- Host: 호스트명 또는 IP 주소 (필수)
- User: SSH 사용자명 (필수)
- Port: SSH 포트 (선택, 기본값: 22)
- 인증 방법:
  - 1: SSH 키
  - 2: 비밀번호

#### `profile list` / `profile ls`

모든 프로파일 목록을 표시합니다.

```bash
./sshclient profile list
./sshclient profile ls  # 단축 명령
```

**출력**:
- 커스텀 프로파일 (`~/.sshclient/config.yaml`)
- SSH config 프로파일 (`~/.ssh/config`)

#### `profile show <name>`

특정 프로파일의 상세 정보를 표시합니다.

```bash
./sshclient profile show myserver
```

**출력 예시**:
```
Profile: myserver
────────────────────────────────────────
  Host:     example.com
  User:     root
  Port:     22
  Key:      /Users/username/.ssh/id_rsa
```

#### `profile remove <name>` / `profile rm <name>`

프로파일을 삭제합니다 (커스텀 프로파일만 가능).

```bash
./sshclient profile remove myserver
./sshclient profile rm myserver  # 단축 명령
```

**주의**: SSH config (`~/.ssh/config`)의 프로파일은 삭제할 수 없습니다.

### 연결 명령어

#### 프로파일로 연결

```bash
# 대화형 셸
./sshclient @profile

# 원격 명령 실행
./sshclient @profile command
./sshclient @profile "command with args"
```

#### 전통적인 SSH 스타일

```bash
# 대화형 셸
./sshclient user@host

# 원격 명령 실행
./sshclient user@host command
./sshclient user@host "command with args"

# SSH 키 지정
./sshclient user@host -key ~/.ssh/custom_key
```

#### 플래그 스타일

```bash
# 대화형 셸
./sshclient -host example.com -user myuser -i

# 원격 명령 실행
./sshclient -host example.com -user myuser -cmd "ls -la"

# SSH 키 지정
./sshclient -host example.com -user myuser -key ~/.ssh/custom_key -i
```

## 프로파일 설정 파일

### 커스텀 프로파일 (YAML 형식)

**위치**:
- macOS/Linux: `~/.sshclient/config.yaml`
- Windows: `C:\Users\사용자명\.sshclient\config.yaml`

**형식**:

```yaml
profiles:
  myserver:
    host: example.com         # 필수: 호스트명 또는 IP
    user: root                # 필수: SSH 사용자명
    port: "22"                # 선택: 포트 (기본값: 22)
    key: /path/to/key         # 선택: SSH 개인키 경로
    encrypted_password: "..." # 선택: 암호화된 비밀번호 (AES-256-GCM)

  # 예시 1: SSH 키 인증 (권장)
  webserver:
    host: web.example.com
    user: deploy
    port: "22"
    key: /Users/username/.ssh/id_rsa

  # 예시 2: 비밀번호 인증 (암호화)
  testserver:
    host: test.example.com
    user: admin
    port: "2222"
    encrypted_password: "kR7v+2Tm...암호화된데이터...=="

  # 예시 3: 최소 설정 (SSH 키 자동 감지)
  minimal:
    host: minimal.example.com
    user: user
```

**필드 설명**:

| 필드 | 타입 | 필수 | 설명 |
|------|------|------|------|
| `host` | string | ✅ | 호스트명 또는 IP 주소 |
| `user` | string | ✅ | SSH 사용자명 |
| `port` | string | ❌ | SSH 포트 (기본값: "22") |
| `key` | string | ❌ | SSH 개인키 경로 (절대 경로 권장) |
| `encrypted_password` | string | ❌ | AES-256-GCM 암호화된 비밀번호 |
| `password` | string | ❌ | **비권장**: 평문 비밀번호 (하위 호환성) |

### SSH Config 호환

기존 SSH config 파일도 자동으로 읽습니다:

**위치**:
- macOS/Linux: `~/.ssh/config`
- Windows: `C:\Users\사용자명\.ssh\config`

**지원하는 옵션**:
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

**지원하는 SSH config 옵션**:
- `Host` - 프로파일 이름
- `HostName` - 실제 호스트명 (설정되지 않으면 Host 값 사용)
- `User` - 사용자명
- `Port` - 포트 번호
- `IdentityFile` - SSH 키 경로 (`~` 확장 지원)

**지원하지 않는 옵션**:
- `ProxyCommand`, `ProxyJump` - 프록시 기능
- `LocalForward`, `RemoteForward` - 포트 포워딩
- `DynamicForward` - SOCKS 프록시
- `ServerAliveInterval`, `ServerAliveCountMax` - Keep-alive 설정

## 플래그 참조

### 연결 옵션

| 플래그 | 타입 | 기본값 | 설명 |
|--------|------|--------|------|
| `-host` | string | - | SSH 서버 호스트명 또는 IP 주소 |
| `-port` | string | "22" | SSH 서버 포트 |
| `-user` | string | - | SSH 사용자명 |
| `-password` | string | - | SSH 비밀번호 (비권장) |
| `-key` | string | - | SSH 개인키 파일 경로 |

### 실행 옵션

| 플래그 | 타입 | 기본값 | 설명 |
|--------|------|--------|------|
| `-i` | bool | false | 대화형 셸 시작 |
| `-cmd` | string | - | 원격에서 실행할 명령어 |
| `-version` | bool | - | 버전 정보 출력 |

**참고**: `user@host` 형식 사용 시 `-i` 플래그 없이도 대화형 모드가 기본값입니다.

## FAQ

### Q1: 프로파일과 SSH config 중 어느 것을 사용해야 하나요?

**커스텀 프로파일** (`~/.sshclient/config.yaml`)을 권장합니다:
- 대화형으로 쉽게 추가 (`profile add`)
- 비밀번호 자동 암호화
- 프로파일 관리 명령어 사용 가능 (list, show, remove)

**SSH config** (`~/.ssh/config`)는 다음 경우 유용합니다:
- 이미 많은 설정이 있는 경우
- 다른 SSH 도구와 설정을 공유하려는 경우
- 수동으로 파일을 직접 편집하는 것을 선호하는 경우

### Q2: 비밀번호를 저장하는 것이 안전한가요?

네, 안전합니다:
- **AES-256-GCM** 암호화 (산업 표준)
- **PBKDF2** (100,000 iterations)로 키 파생
- 설정 파일을 열어도 암호화된 문자열만 표시됩니다

하지만 더 높은 보안을 위해서는 **SSH 키 인증**을 권장합니다.

### Q3: 프로파일 우선순위는 어떻게 되나요?

프로파일 검색 순서:
1. **커스텀 프로파일** (`~/.sshclient/config.yaml`)
2. **SSH config** (`~/.ssh/config`)

같은 이름의 프로파일이 두 곳에 있으면 커스텀 프로파일이 우선됩니다.

### Q4: 프로파일을 백업하려면 어떻게 하나요?

```bash
# 백업
cp ~/.sshclient/config.yaml ~/.sshclient/config.yaml.backup

# 다른 머신으로 복사
scp ~/.sshclient/config.yaml user@newmachine:~/.sshclient/
```

**주의**: 암호화된 비밀번호는 이 프로그램의 내부 키로 암호화되어 있으므로, 다른 머신에서도 동일하게 작동합니다.

### Q5: 여러 SSH 키를 관리하려면?

프로파일별로 다른 SSH 키를 지정할 수 있습니다:

```yaml
profiles:
  work-server:
    host: work.example.com
    user: myname
    key: /Users/username/.ssh/work_rsa

  personal-server:
    host: personal.example.com
    user: myname
    key: /Users/username/.ssh/personal_rsa
```

### Q6: 프록시나 점프 호스트를 거쳐 연결할 수 있나요?

현재 버전은 직접 프록시나 점프 호스트 기능을 지원하지 않습니다.
대신 일반 `ssh` 명령의 `-J` 옵션을 사용하거나,
점프 호스트에 먼저 연결한 후 다시 sshclient를 사용하세요.

### Q7: 포트 포워딩을 사용할 수 있나요?

현재 버전은 포트 포워딩을 지원하지 않습니다.
일반 `ssh` 명령의 `-L`, `-R`, `-D` 옵션을 사용하세요.

## 문제 해결

### 문제: "Failed to read password: operation not supported on socket"

**원인**: 표준 입력이 터미널이 아닌 경우 (파이프나 리다이렉션 사용 시)

**해결**:
```bash
# ❌ 작동 안 함
echo "password" | ./sshclient user@host

# ✅ 작동
./sshclient user@host -key ~/.ssh/id_rsa

# ✅ 또는 프로파일에 비밀번호 저장
./sshclient profile add myserver
./sshclient @myserver
```

### 문제: "Failed to parse private key"

**원인**: 잘못된 SSH 키 파일 또는 암호화된 키

**해결**:
```bash
# 1. SSH 키 형식 확인
head -1 ~/.ssh/id_rsa
# "-----BEGIN RSA PRIVATE KEY-----" 또는
# "-----BEGIN OPENSSH PRIVATE KEY-----" 이어야 함

# 2. 암호화되지 않은 키 생성
ssh-keygen -t ed25519 -N "" -f ~/.ssh/id_ed25519_nopass

# 3. 기존 암호화된 키의 암호 제거
ssh-keygen -p -f ~/.ssh/id_rsa -N ""
```

### 문제: "Connection refused" 또는 "No route to host"

**원인**: SSH 서버가 실행 중이 아니거나 방화벽 차단

**해결**:
```bash
# 1. 호스트 도달 가능 여부 확인
ping example.com

# 2. SSH 포트 열려 있는지 확인
nc -zv example.com 22
telnet example.com 22

# 3. 올바른 포트 확인
nmap -p 22,2222 example.com
```

### 문제: 프로파일이 목록에 안 나타남

**원인**: YAML 형식 오류 또는 잘못된 파일 위치

**해결**:
```bash
# 1. 설정 파일 위치 확인 (macOS/Linux)
ls -la ~/.sshclient/config.yaml

# Windows
dir C:\Users\사용자명\.sshclient\config.yaml

# 2. YAML 형식 검증
cat ~/.sshclient/config.yaml

# 3. YAML 형식이 올바른지 확인
# - profiles: 키가 최상위에 있어야 함
# - 들여쓰기는 공백 2칸 또는 4칸 일관성 유지
# - 탭 문자 사용 금지
```

올바른 형식 예시:
```yaml
profiles:
  myserver:
    host: example.com
    user: root
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
# 기존 SSH config 파일을 그대로 사용
./sshclient @myhost  # ~/.ssh/config의 Host myhost 사용
```

프로파일 생성 필요 없이 바로 사용할 수 있습니다!

### 💡 Tip 4: 대량 서버 작업 스크립트

```bash
#!/bin/bash
# 모든 프로파일 서버에서 uptime 확인

for profile in $(./sshclient profile list | grep '@' | awk '{print $1}' | tr -d '@'); do
    echo "=== $profile ==="
    ./sshclient @$profile uptime
done
```

### 💡 Tip 5: 환경 변수로 설정 관리

```bash
# .env 파일
export SSH_PROFILE=production
export SSH_CMD="df -h"

# 사용
./sshclient @$SSH_PROFILE "$SSH_CMD"
```

## 아키텍처

### 프로젝트 구조

```
sshclient/
├── main.go          # CLI 인터페이스 및 메인 로직
├── client.go        # SSH 클라이언트 핵심 구현
├── config.go        # 프로파일 설정 파일 관리 (YAML + SSH config)
├── profile.go       # 프로파일 관리 명령어 (add, list, remove, show)
├── crypto.go        # AES-256 암호화/복호화 (v1.2.0)
├── go.mod           # Go 모듈 정의
├── go.sum           # 의존성 체크섬
├── README.md        # 프로젝트 소개
├── CHANGELOG.md     # 변경 이력
├── CLAUDE.md        # 개발 가이드
└── docs/            # 사용자 문서
    ├── User-Guide.md
    └── User-Manual.md
```

### 작동 원리

#### 1. 순수 Go SSH 구현

OS의 `ssh` 명령어를 사용하지 않고, `golang.org/x/crypto/ssh` 라이브러리로 SSH 프로토콜을 직접 구현:

1. **TCP 연결**: `net.Dial()`로 SSH 서버에 TCP 연결
2. **SSH 핸드셰이크**: `golang.org/x/crypto/ssh`로 SSH 프로토콜 핸드셰이크 및 암호화 통신
3. **인증**: 비밀번호 또는 SSH 키로 인증
4. **세션 관리**: SSH 세션 생성 및 명령 실행

#### 2. 프로파일 시스템

- YAML 파일(`~/.sshclient/config.yaml`)로 프로파일 저장
- SSH config(`~/.ssh/config`) 파서로 기존 설정 읽기
- 우선순위: 커스텀 프로파일 → SSH config

#### 3. 비밀번호 암호화

- 프로파일에 저장되는 비밀번호를 AES-256-GCM으로 암호화
- PBKDF2 (100,000 iterations)로 키 파생
- 내부 passphrase로 자동 암호화/복호화 (마스터 비밀번호 불필요)

### 의존성

- `golang.org/x/crypto/ssh` - SSH 프로토콜 구현
- `golang.org/x/term` - 터미널 제어 (크로스 플랫폼)
- `gopkg.in/yaml.v3` - YAML 설정 파일 파싱

### 빌드 정보

- **컴파일된 바이너리**: 약 6-7MB
- **플랫폼**: Windows, macOS, Linux (x86_64, ARM64)
- **Go 버전**: 1.16 이상

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

macOS/Linux:
```bash
~/.sshclient/           # 0700 (drwx------)
~/.sshclient/config.yaml # 0600 (-rw-------)
```

Windows:
- 사용자 홈 디렉토리의 기본 권한 사용
- 민감한 정보는 AES-256 암호화로 보호

### 권장사항

1. **SSH 키 인증 사용** (비밀번호보다 안전)
2. **비밀번호를 프로파일에 저장하지 않기** (접속 시 입력)
3. **정기적인 SSH 키 교체**
4. **프로파일 파일 백업 시 주의** (암호화되었지만 추가 보안 필요)

---

**추가 도움이 필요하면**:
- GitHub Issues: https://github.com/arkd0ng/sshclient/issues
- 개발자 가이드: [CLAUDE.md](../CLAUDE.md)
