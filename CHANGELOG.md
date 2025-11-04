# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.2.1] - 2025-11-04

### Added
- **호스트 키 검증** - MITM 공격 방지 (중요한 보안 개선!)
  - `~/.sshclient/known_hosts` 파일로 호스트 키 관리
  - 기존 `~/.ssh/known_hosts` 자동 복사 및 마이그레이션
  - 최초 접속 시 호스트 키 확인 후 저장
  - 재접속 시 호스트 키 자동 검증
  - 키 변경 감지 시 경고 및 사용자 확인
- Feature Roadmap 문서 추가 (`docs/Feature-Roadmap.md`)
  - 향후 추가 기능 계획
  - 우선순위 및 예상 시간
  - 구현 난이도 평가

### Changed
- 프로젝트 구조 개선
  - 소스 코드를 `src/` 디렉토리로 이동
  - 빌드 결과물을 `bin/` 디렉토리로 통합
- `Makefile` 추가로 빌드 자동화
  - `make build`, `make build-all` 등 편리한 명령어
  - 크로스 플랫폼 빌드 간소화

### Security
- **중요**: `ssh.InsecureIgnoreHostKey()` 제거
- 실제 호스트 키 검증으로 MITM 공격 방지
- 프로덕션 환경에서 안전하게 사용 가능

### Fixed
- 호스트 키 검증 부재로 인한 보안 취약점 해결

## [1.2.0] - 2025-11-04

### Added
- **프로파일 시스템**: 자주 사용하는 서버 정보를 저장하고 빠르게 접속
  - `profile add` - 대화형 프로파일 생성
  - `profile list` - 모든 프로파일 목록 보기
  - `profile show` - 프로파일 상세 정보
  - `profile remove` - 프로파일 삭제
- **SSH config 호환**: `~/.ssh/config` 파일 자동 읽기 및 사용
- **AES-256 비밀번호 암호화**: 프로파일에 저장되는 비밀번호를 AES-256-GCM으로 자동 암호화
  - PBKDF2 (100,000 iterations) 키 파생
  - 마스터 비밀번호 불필요 (자동 암호화/복호화)
- **친절한 도움말 시스템**: 상황별 맞춤 도움말과 사용 예시
- **전통적인 SSH 스타일 지원**: `user@host` 형식으로 직접 접속
- **SSH 키 자동 감지**: `~/.ssh/id_rsa`, `~/.ssh/id_ed25519`, `~/.ssh/id_ecdsa` 자동 검색
- **크로스 플랫폼 지원**: Windows, macOS, Linux 완벽 지원
  - `golang.org/x/term` 패키지로 터미널 제어 통합

### Changed
- 의존성 업데이트: `golang.org/x/crypto/ssh/terminal` → `golang.org/x/term`으로 변경 (크로스 플랫폼 호환성)
- 비밀번호 입력 시 `syscall.Stdin` → `os.Stdin.Fd()`로 변경 (Windows 호환성)
- 프로젝트 구조 개선:
  - 문서를 `docs/` 폴더로 분리
  - `User-Guide.md`, `User-Manual.md` 추가
  - README.md를 간결하게 재구성

### Security
- AES-256-GCM 암호화 구현 (평문 비밀번호 저장 방지)
- PBKDF2 키 파생 함수 적용 (100,000 iterations)
- 설정 파일 권한 자동 설정 (macOS/Linux: 0600)

### Breaking Changes
- v1.2.0 베타 버전(마스터 비밀번호 시스템)에서 생성한 암호화된 비밀번호는 새 버전과 호환되지 않음
  - 해결: 프로파일 재생성 필요 (`profile remove` → `profile add`)

## [1.1.0] - 2025-11-03

### Added
- 전통적인 SSH 스타일 지원: `sshclient user@host [command]`
- 명령줄 인자로 원격 명령 실행: `sshclient @profile command`
- 플래그 파싱 개선 (명령어와 플래그 분리)

### Changed
- 대화형 모드를 기본값으로 변경 (user@host 형식 사용 시)

## [1.0.0] - 2025-11-02

### Added
- 초기 릴리스
- SSH 클라이언트 핵심 기능:
  - 비밀번호 인증
  - SSH 키 인증 (RSA, Ed25519, ECDSA)
  - 대화형 셸 세션
  - 원격 명령 실행
  - SCP 파일 전송 (`CopyFile`, `DownloadFile`)
- 순수 Go 구현 (`golang.org/x/crypto/ssh`)
- 크로스 플랫폼 바이너리 빌드

### Security
- `ssh.InsecureIgnoreHostKey()` 사용 (테스트/개발 목적)

---

## Legend

- **Added**: 새로운 기능
- **Changed**: 기존 기능 변경
- **Deprecated**: 곧 제거될 기능
- **Removed**: 제거된 기능
- **Fixed**: 버그 수정
- **Security**: 보안 관련 변경

---

## Versioning

이 프로젝트는 [Semantic Versioning](https://semver.org/)을 따릅니다:

- **MAJOR** version: 호환되지 않는 API 변경
- **MINOR** version: 하위 호환되는 기능 추가
- **PATCH** version: 하위 호환되는 버그 수정

[1.2.0]: https://github.com/arkd0ng/sshclient/releases/tag/v1.2.0
[1.1.0]: https://github.com/arkd0ng/sshclient/releases/tag/v1.1.0
[1.0.0]: https://github.com/arkd0ng/sshclient/releases/tag/v1.0.0
