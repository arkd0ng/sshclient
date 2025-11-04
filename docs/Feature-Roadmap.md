# SSH Client - 기능 로드맵 및 개선 제안

현재 프로젝트의 상태를 검토하고 추가할만한 기능과 개선사항을 제안합니다.

## 현재 구현 상태 (v1.2.0)

### ✅ 완료된 핵심 기능
- [x] SSH 클라이언트 (비밀번호/키 인증)
- [x] 대화형 셸
- [x] 원격 명령 실행
- [x] 프로파일 시스템
- [x] SSH config 호환
- [x] AES-256 비밀번호 암호화
- [x] 크로스 플랫폼 지원 (Windows, macOS, Linux)
- [x] SCP 파일 전송 (구현됨, CLI 미노출)

### ⚠️ 보안 주의사항 (현재)
- `ssh.InsecureIgnoreHostKey()` 사용 중 (호스트 키 검증 안 함)
- 프로덕션 환경에서는 위험

---

## 우선순위별 개선 제안

### 🔴 높음 (High Priority) - 보안 및 안정성

#### 1. 호스트 키 검증 구현 ⚠️ **중요**
**현재 문제**:
```go
// client.go:31, 61
HostKeyCallback: ssh.InsecureIgnoreHostKey(),  // 보안 취약!
```

**제안**:
```go
// 호스트 키 저장 및 검증
~/.sshclient/known_hosts 파일 사용
- 최초 접속 시 호스트 키 저장
- 재접속 시 호스트 키 검증
- 키 변경 감지 시 경고
```

**구현 난이도**: 중간
**예상 시간**: 4-6 시간

**가치**: ⭐⭐⭐⭐⭐
- MITM 공격 방지
- 프로덕션 사용 가능

#### 2. 연결 타임아웃 설정
**현재**: 무한 대기 가능
**제안**:
```yaml
# 프로파일에 타임아웃 추가
profiles:
  myserver:
    host: example.com
    user: root
    timeout: 30  # 초 단위
```

**구현 난이도**: 낮음
**예상 시간**: 1-2 시간

**가치**: ⭐⭐⭐⭐
- 응답 없는 서버 대응
- 스크립트 안정성 향상

#### 3. 재연결 로직
**제안**:
```bash
# 연결 끊김 시 자동 재연결
sshclient @myserver --retry 3 --retry-delay 5
```

**구현 난이도**: 중간
**예상 시간**: 3-4 시간

**가치**: ⭐⭐⭐⭐
- 불안정한 네트워크 대응

#### 4. 테스트 코드 작성 🧪
**현재**: 테스트 코드 전무
**제안**:
```go
// src/client_test.go
// src/config_test.go
// src/crypto_test.go
// src/profile_test.go

// 단위 테스트
// 통합 테스트
// 벤치마크
```

**구현 난이도**: 중간
**예상 시간**: 8-12 시간

**가치**: ⭐⭐⭐⭐⭐
- 버그 예방
- 리팩토링 안전성
- CI/CD 통합

---

### 🟡 중간 (Medium Priority) - 사용성 개선

#### 5. SCP 파일 전송 CLI 노출
**현재**: CopyFile(), DownloadFile() 함수 구현됨, 명령줄 인터페이스 없음
**제안**:
```bash
# 파일 업로드
sshclient @myserver upload local.txt /remote/path/

# 파일 다운로드
sshclient @myserver download /remote/file.txt ./local/

# 디렉토리 업로드/다운로드
sshclient @myserver upload -r ./localdir/ /remote/dir/
```

**구현 난이도**: 낮음
**예상 시간**: 2-3 시간

**가치**: ⭐⭐⭐⭐
- SCP 기능 활용
- 파일 전송 편의성

#### 6. 프로파일 편집 기능
**제안**:
```bash
# 프로파일 편집
sshclient profile edit myserver

# 대화형 또는 에디터 열기
```

**구현 난이도**: 낮음
**예상 시간**: 2 시간

**가치**: ⭐⭐⭐
- 프로파일 수정 편의성

#### 7. 프로파일 복사/백업
**제안**:
```bash
# 프로파일 복사
sshclient profile copy myserver myserver-backup

# 프로파일 내보내기/가져오기
sshclient profile export myserver > myserver.yaml
sshclient profile import myserver.yaml
```

**구현 난이도**: 낮음
**예상 시간**: 2-3 시간

**가치**: ⭐⭐⭐
- 프로파일 관리 편의성

#### 8. 연결 로그 기록
**제안**:
```yaml
# 옵션으로 로그 활성화
logging:
  enabled: true
  path: ~/.sshclient/logs/
  level: info  # debug, info, warn, error

# 로그 내용
# - 연결 시도
# - 인증 방법
# - 명령 실행 기록 (옵션)
# - 오류
```

**구현 난이도**: 중간
**예상 시간**: 3-4 시간

**가치**: ⭐⭐⭐
- 디버깅 편의성
- 감사 추적

#### 9. 명령 히스토리
**제안**:
```bash
# 최근 실행한 명령 보기
sshclient history

# 명령 재실행
sshclient history run 5
```

**구현 난이도**: 중간
**예상 시간**: 3-4 시간

**가치**: ⭐⭐⭐
- 반복 작업 편의성

#### 10. 프로파일 그룹/태그
**제안**:
```yaml
profiles:
  web1:
    host: web1.example.com
    user: deploy
    tags: [production, web]

  web2:
    host: web2.example.com
    user: deploy
    tags: [production, web]

# 사용
sshclient profile list --tag production
sshclient @tag:web uptime  # 모든 web 태그 서버에 실행
```

**구현 난이도**: 중간-높음
**예상 시간**: 5-7 시간

**가치**: ⭐⭐⭐⭐
- 대량 서버 관리

---

### 🟢 낮음 (Low Priority) - 고급 기능

#### 11. 포트 포워딩 (로컬/원격) 🚀
**제안**:
```bash
# 로컬 포트 포워딩
sshclient @myserver -L 8080:localhost:80

# 원격 포트 포워딩
sshclient @myserver -R 9090:localhost:3000

# 다이나믹 포트 포워딩 (SOCKS 프록시)
sshclient @myserver -D 1080
```

**구현 난이도**: 높음
**예상 시간**: 8-12 시간

**가치**: ⭐⭐⭐⭐⭐
- 데이터베이스 접속
- 개발 환경 구성
- 보안 터널링

#### 12. 점프 호스트 (ProxyJump)
**제안**:
```yaml
# 프로파일에 점프 호스트 설정
profiles:
  internal-server:
    host: internal.local
    user: admin
    jump_host: bastion

  bastion:
    host: bastion.example.com
    user: jump-user

# 사용
sshclient @internal-server
# bastion을 거쳐서 internal-server 접속
```

**구현 난이도**: 높음
**예상 시간**: 10-15 시간

**가치**: ⭐⭐⭐⭐⭐
- 보안 강화된 인프라 접근

#### 13. 터널 유지 (Keep-Alive)
**제안**:
```yaml
# 프로파일에 keep-alive 설정
profiles:
  myserver:
    host: example.com
    user: root
    keep_alive:
      interval: 60  # 60초마다 패킷 전송
      max_count: 3  # 3회 실패 시 종료
```

**구현 난이도**: 중간
**예상 시간**: 2-3 시간

**가치**: ⭐⭐⭐
- 장시간 세션 유지

#### 14. SSH 에이전트 지원
**제안**:
```bash
# SSH 에이전트의 키 사용
sshclient @myserver --agent

# 에이전트 포워딩
sshclient @myserver --agent-forward
```

**구현 난이도**: 중간-높음
**예상 시간**: 6-8 시간

**가치**: ⭐⭐⭐⭐
- 키 관리 편의성

#### 15. 설정 파일 마이그레이션
**제안**:
```bash
# SSH config → sshclient 프로파일 변환
sshclient import ssh-config ~/.ssh/config

# 다른 도구에서 가져오기
sshclient import putty sessions.reg
```

**구현 난이도**: 중간
**예상 시간**: 4-6 시간

**가치**: ⭐⭐⭐
- 마이그레이션 편의성

#### 16. 명령 템플릿/스크립트
**제안**:
```yaml
# 프로파일에 명령 템플릿 저장
profiles:
  web1:
    host: web1.example.com
    user: deploy
    scripts:
      deploy: |
        cd /var/www/app
        git pull
        npm install
        pm2 restart app

      backup: |
        tar -czf backup-$(date +%Y%m%d).tar.gz /var/www/app

# 사용
sshclient @web1 run deploy
sshclient @web1 run backup
```

**구현 난이도**: 중간
**예상 시간**: 4-5 시간

**가치**: ⭐⭐⭐⭐
- 배포 자동화

#### 17. 병렬 명령 실행
**제안**:
```bash
# 여러 서버에 동시 명령 실행
sshclient @web1,@web2,@web3 uptime

# 태그 기반
sshclient @tag:production "systemctl restart nginx"

# 결과 집계
sshclient @tag:web "df -h /" --format table
```

**구현 난이도**: 높음
**예상 시간**: 8-10 시간

**가치**: ⭐⭐⭐⭐⭐
- 대규모 서버 관리

#### 18. 인터랙티브 서버 선택 UI
**제안**:
```bash
# 대화형 서버 선택
sshclient connect
# → 화살표 키로 프로파일 선택
# → 필터링 지원
# → 최근 접속 서버 표시
```

**구현 난이도**: 중간
**예상 시간**: 4-6 시간

**가치**: ⭐⭐⭐
- 사용자 경험 개선

---

## 문서 개선 제안

### 📚 추가 문서

#### 1. CONTRIBUTING.md
- 기여 가이드라인
- 코드 스타일
- PR 프로세스
- 이슈 템플릿

#### 2. Security-Guide.md
- 보안 모범 사례
- 취약점 보고 절차
- 보안 설정 가이드

#### 3. Architecture.md
- 상세한 아키텍처 문서
- 모듈 간 의존성
- 데이터 흐름
- 확장 가이드

#### 4. API-Reference.md (라이브러리로 사용 시)
- 공개 API 문서
- 사용 예제
- Go 패키지로 사용하는 방법

---

## 품질 개선

### 🧪 테스트 커버리지
- [ ] 단위 테스트 (목표: 80%)
- [ ] 통합 테스트
- [ ] E2E 테스트
- [ ] 벤치마크 테스트

### 🔧 CI/CD
- [ ] GitHub Actions 설정
- [ ] 자동 빌드
- [ ] 자동 테스트
- [ ] 자동 릴리스

### 📦 릴리스 프로세스
- [ ] GitHub Releases
- [ ] 플랫폼별 바이너리 자동 빌드
- [ ] Homebrew tap (macOS)
- [ ] Chocolatey (Windows)
- [ ] APT/YUM 저장소 (Linux)

---

## 우선순위 요약

### 즉시 구현 권장 (1-2주)
1. ⭐⭐⭐⭐⭐ 호스트 키 검증 (보안)
2. ⭐⭐⭐⭐⭐ 테스트 코드 작성
3. ⭐⭐⭐⭐ SCP CLI 노출
4. ⭐⭐⭐⭐ 연결 타임아웃

### 단기 목표 (1-2개월)
5. ⭐⭐⭐⭐⭐ 포트 포워딩
6. ⭐⭐⭐⭐ 재연결 로직
7. ⭐⭐⭐⭐ 프로파일 그룹/태그
8. ⭐⭐⭐ 연결 로그

### 중장기 목표 (3-6개월)
9. ⭐⭐⭐⭐⭐ 점프 호스트
10. ⭐⭐⭐⭐⭐ 병렬 명령 실행
11. ⭐⭐⭐⭐ 명령 템플릿/스크립트
12. ⭐⭐⭐⭐ SSH 에이전트

---

## 결론

**현재 상태**: 기본 기능은 매우 잘 구현되어 있음
**주요 약점**: 보안 (호스트 키 검증), 테스트 부재
**강점**: 크로스 플랫폼, 깔끔한 코드, 좋은 문서

**추천 로드맵**:
1. **v1.2.1** (버그픽스) - 호스트 키 검증, 타임아웃
2. **v1.3.0** (기능 추가) - SCP CLI, 포트 포워딩, 테스트
3. **v1.4.0** (고급 기능) - 점프 호스트, 병렬 실행, 스크립트
4. **v2.0.0** (메이저 업데이트) - 전면 리팩토링, API 안정화

이 로드맵을 따르면 **프로덕션 레벨의 엔터프라이즈 SSH 클라이언트**로 성장할 수 있습니다! 🚀
