# SSH Client

크로스 플랫폼 터미널 SSH 클라이언트 - 순수 Go 구현

[![Go Version](https://img.shields.io/badge/Go-1.16+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20macOS%20%7C%20Linux-lightgrey)](https://github.com/arkd0ng/sshclient)

## 특징

✅ **완전히 독립적**: OS의 `ssh` 명령어 없이 단일 바이너리로 완전 동작
✅ **권한 불필요**: root/sudo 권한 없이 일반 사용자 권한으로 실행
✅ **크로스 플랫폼**: Windows, macOS, Linux 완벽 지원
✅ **프로파일 시스템**: 자주 사용하는 서버 정보를 저장하고 빠르게 접속
✅ **SSH config 호환**: 기존 `~/.ssh/config` 파일 자동 읽기
✅ **AES-256 암호화**: 저장된 비밀번호를 자동 암호화/복호화
✅ **친절한 도움말**: 상황별 맞춤 도움말과 사용 예시

## 빠른 시작

### 설치

**1. Go가 설치되어 있는 경우**:
```bash
git clone https://github.com/arkd0ng/sshclient.git
cd sshclient
make build
```

**2. 바이너리 다운로드** (출시 예정):
```bash
# GitHub Releases에서 플랫폼별 바이너리 다운로드
```

### 기본 사용법

#### 방법 1: 전통적인 SSH 스타일
```bash
# 대화형 셸
bin/sshclient user@hostname

# 원격 명령 실행
bin/sshclient user@hostname ls -la
```

#### 방법 2: 프로파일 사용 (권장)
```bash
# 프로파일 생성
bin/sshclient profile add myserver

# 프로파일로 접속
bin/sshclient @myserver

# 원격 명령 실행
bin/sshclient @myserver uptime
```

#### 도움말
```bash
bin/sshclient -h              # 전체 도움말
bin/sshclient profile help    # 프로파일 도움말
```

## 빌드

### Makefile 사용 (권장)

```bash
# 현재 플랫폼 빌드
make build

# 모든 플랫폼 빌드
make build-all

# 특정 플랫폼 빌드
make windows       # Windows 64-bit
make linux         # Linux 64-bit
make darwin        # macOS Intel
make darwin-arm64  # macOS Apple Silicon

# 도움말
make help
```

### 수동 빌드

```bash
# Windows (64-bit)
GOOS=windows GOARCH=amd64 go build -o bin/sshclient.exe src/*.go

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o bin/sshclient-darwin src/*.go

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o bin/sshclient-darwin-arm64 src/*.go

# Linux (64-bit)
GOOS=linux GOARCH=amd64 go build -o bin/sshclient-linux src/*.go
```

## 문서

📚 **[사용자 가이드](docs/User-Guide.md)** - 빠른 시작, 프로파일 사용법, 실전 시나리오
📖 **[사용자 매뉴얼](docs/User-Manual.md)** - 명령어 레퍼런스, FAQ, 문제 해결, 아키텍처
📝 **[CHANGELOG](CHANGELOG.md)** - 버전별 변경 이력
🔧 **[CLAUDE.md](CLAUDE.md)** - 개발자 가이드

## 주요 기능

### 프로젝트 구조

```
sshclient/
├── src/               # 소스 코드
│   ├── main.go
│   ├── client.go
│   ├── config.go
│   ├── profile.go
│   └── crypto.go
├── bin/               # 빌드된 바이너리 (gitignored)
├── docs/              # 사용자 문서
│   ├── User-Guide.md
│   └── User-Manual.md
├── Makefile           # 빌드 자동화
├── README.md
├── CHANGELOG.md
├── CLAUDE.md
├── go.mod
└── go.sum
```

### 프로파일 시스템

```bash
# 프로파일 관리
bin/sshclient profile add webserver    # 대화형 프로파일 생성
bin/sshclient profile list              # 모든 프로파일 목록
bin/sshclient profile show webserver    # 프로파일 상세 정보
bin/sshclient profile remove webserver  # 프로파일 삭제
```

### 인증 방법

- **SSH 키 인증** (권장): 자동으로 `~/.ssh/id_rsa`, `id_ed25519`, `id_ecdsa` 검색
- **비밀번호 인증**: 접속 시 입력 또는 프로파일에 암호화하여 저장

### 보안

- **호스트 키 검증** (v1.2.1) - MITM 공격 방지
  - `~/.sshclient/known_hosts` 파일로 호스트 키 관리
  - 기존 `~/.ssh/known_hosts` 자동 복사/마이그레이션
  - 최초 접속 시 호스트 키 확인 후 저장
  - 키 변경 감지 시 경고
- **AES-256-GCM** 비밀번호 암호화
- **PBKDF2** (100,000 iterations) 키 파생
- 마스터 비밀번호 없이 자동 암호화/복호화

## 시스템 요구사항

- **OS**: Windows 10+, macOS 10.12+, Linux (현대적인 배포판)
- **Go**: 1.16+ (빌드 시)
- **메모리**: 최소 64MB
- **디스크**: 약 10MB

## 기술 스택

- **언어**: Go 1.16+
- **SSH 구현**: `golang.org/x/crypto/ssh`
- **터미널 제어**: `golang.org/x/term` (크로스 플랫폼)
- **설정 파일**: `gopkg.in/yaml.v3`

## 라이선스

MIT License - 자세한 내용은 [LICENSE](LICENSE) 파일을 참조하세요.

## 기여

이슈와 풀 리퀘스트를 환영합니다!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 지원

문제가 발생하거나 질문이 있으신가요?

- 📖 [사용자 매뉴얼](docs/User-Manual.md#문제-해결) 문제 해결 섹션 확인
- 🐛 [GitHub Issues](https://github.com/arkd0ng/sshclient/issues)에 이슈 등록
- 💬 [Discussions](https://github.com/arkd0ng/sshclient/discussions)에서 토론

## 개발 상태

- ✅ v1.2.0 - 프로파일 시스템, AES 암호화, 크로스 플랫폼 지원
- 🔄 v1.3.0 (예정) - 포트 포워딩, 점프 호스트 지원
- 🔄 v1.4.0 (예정) - GUI 지원

---

Made with ❤️ using [Go](https://golang.org)
