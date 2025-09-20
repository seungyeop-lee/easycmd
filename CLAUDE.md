# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 프로젝트 개요

이 프로젝트는 Go 언어로 작성된 `easycmd` 라이브러리로, 명령어 실행을 간편하게 해주는 도구입니다. 외부 명령어를 실행하고 표준 입출력을 제어할 수 있는 래퍼 기능을 제공합니다.

## 개발 명령어

### 테스트 실행
```bash
# 모든 테스트 실행
go test ./...

# 특정 테스트 실행
go test ./test -run TestSimple

# 테스트 출력과 함께 실행
go test -v ./test
```

### 빌드
```bash
# 현재 플랫폼용 빌드
go build

# 모듈 의존성 관리
go mod tidy
go mod download
```

### 코드 품질 검증
```bash
# 코드 포맷팅
go fmt ./...

# 린팅 (golangci-lint 설치 필요)
golangci-lint run

# 정적 분석
go vet ./...
```

## 아키텍처

### 핵심 구조
- **모듈형 설계**: `easycmd.go`(메인 로직), `command.go`(명령어 타입), `config.go`(설정 관리), `parser.go`(명령어 파싱), `logger.go`(로깅)로 분리된 구조
- **타입 기반 설계**: Go의 타입 시스템을 활용한 안전한 명령어 처리
- **설정 주입 패턴**: 함수형 옵션 패턴을 사용한 유연한 설정
- **명령어 래핑**: bash/powershell 자동 래핑으로 크로스 플랫폼 지원

### 주요 타입
- `command`: 명령어 문자열을 감싸는 타입으로 bash/powershell 래핑 기능 제공
- `config`: 실행 디렉토리, 표준 입출력, 디버그 모드, 타임아웃, 환경변수 설정을 관리
- `Cmd`: 설정을 포함한 명령어 실행 인스턴스
- `configApply`: 함수형 옵션 패턴을 위한 함수 타입

### 메서드 분류
- `Run()`: 기본 명령어 실행
- `RunShell()`: bash 래핑된 명령어 실행
- `RunPowershell()`: PowerShell 래핑된 명령어 실행
- `RunWithDir()`, `RunShellWithDir()`, `RunPowershellWithDir()`: 특정 디렉토리에서 실행하는 변형들

### 설정 옵션
- `WithDebug()`: 디버그 모드 활성화 (명령어 파싱, 실행 과정 출력)
- `WithStdIn()`: 표준 입력 설정
- `WithStdOut()`: 표준 출력 설정
- `WithStdErr()`: 표준 에러 설정
- `WithTimeout()`: 명령어 실행 타임아웃 설정
- `WithEnv()`: 환경변수 설정

### 테스트 구조
- `test/e2e_test.go`: 실제 명령어 실행을 통한 통합 테스트 (타임아웃, 환경변수, 디버그 모드 포함)
- `parser_test.go`: 명령어 파싱 로직 단위 테스트 (인용부호, 공백 처리)
- `command_test.go`: command 타입의 메서드 단위 테스트 (bash/powershell 래핑 포함)
- 표준 출력 캡처를 통한 결과 검증
- 다양한 실행 방식(단순 명령어, 쉘 명령어, 멀티라인 쉘) 테스트

## 핵심 기능

### 명령어 실행 방식
1. **기본 실행**: `Run()` - 직접 명령어 실행
2. **Shell 실행**: `RunShell()` - bash로 래핑하여 실행 (복잡한 스크립트 지원)
3. **PowerShell 실행**: `RunPowershell()` - Windows PowerShell로 실행

### 고급 기능
- **디렉토리 변경**: 모든 실행 메서드에 `WithDir` 변형 제공
- **타임아웃 설정**: 긴 실행 명령어에 대한 시간 제한
- **환경변수**: 명령어 실행 시 특정 환경변수 설정
- **표준 입출력 제어**: 명령어의 stdin/stdout/stderr 완전 제어

### 디버그 모드
- `WithDebug()` 옵션으로 활성화
- 파싱된 명령어, 실행 명령어, 인수, 실행 디렉토리 출력
- 명령어 시작/완료/실패 상태 출력
- 별도 스트림으로 디버그 출력 분리 가능

## 코딩 컨벤션

- Go 표준 포맷팅 사용 (`go fmt`)
- 타입 안전성을 위한 별칭 타입 활용
- 에러 처리는 명시적으로 수행
- 설정은 함수형 옵션 패턴으로 주입