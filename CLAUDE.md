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
- **단일 파일 라이브러리**: `easycmd.go` 파일에 모든 핵심 로직이 포함됨
- **타입 기반 설계**: Go의 타입 시스템을 활용한 안전한 명령어 처리
- **설정 주입 패턴**: 함수형 옵션 패턴을 사용한 유연한 설정

### 주요 타입
- `command`: 명령어 문자열을 감싸는 타입으로 bash/powershell 래핑 기능 제공
- `Config`: 실행 디렉토리와 표준 입출력 설정을 관리
- `Cmd`: 설정을 포함한 명령어 실행 인스턴스

### 메서드 분류
- `Run()`: 기본 명령어 실행
- `RunShell()`: bash 래핑된 명령어 실행
- `RunPowershell()`: PowerShell 래핑된 명령어 실행
- `RunWithDir()`: 특정 디렉토리에서 실행하는 변형들

### 테스트 구조
- `test/e2e_test.go`: 실제 명령어 실행을 통한 통합 테스트
- 표준 출력 캡처를 통한 결과 검증
- 다양한 실행 방식(단순 명령어, 쉘 명령어, 멀티라인 쉘) 테스트

## 코딩 컨벤션

- Go 표준 포맷팅 사용 (`go fmt`)
- 타입 안전성을 위한 별칭 타입 활용
- 에러 처리는 명시적으로 수행
- 설정은 함수형 옵션 패턴으로 주입