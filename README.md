# easycmd

Go 언어로 작성된 명령어 실행 라이브러리입니다. 외부 명령어를 간편하게 실행하고 표준 입출력을 제어할 수 있는 래퍼 기능을 제공합니다.

## 설치

```bash
go get github.com/seungyeop-lee/easycmd
```

## 기본 사용법

### 간단한 명령어 실행

```go
package main

import (
    "github.com/seungyeop-lee/easycmd"
)

func main() {
    cmd := easycmd.New()
    err := cmd.Run("echo hello world")
    if err != nil {
        panic(err)
    }
}
```

### 표준 출력 캡처

```go
package main

import (
    "bytes"
    "fmt"
    "github.com/seungyeop-lee/easycmd"
)

func main() {
    out := &bytes.Buffer{}
    cmd := easycmd.New(easycmd.WithStdOut(out))

    err := cmd.Run("echo hello world")
    if err != nil {
        panic(err)
    }

    fmt.Println("출력:", out.String()) // 출력: hello world
}
```

## 고급 사용법

### Shell 명령어 실행

```go
cmd := easycmd.New()

// bash로 래핑된 명령어 실행
err := cmd.RunShell("(cd .. && pwd)")

// 멀티라인 shell 스크립트 실행
err = cmd.RunShell(`
    pwd
    ls -al
`)
```

### PowerShell 명령어 실행 (Windows)

```go
cmd := easycmd.New()
err := cmd.RunPowershell("Get-Location")
```

### 특정 디렉토리에서 실행

```go
cmd := easycmd.New()

// 기본 명령어를 특정 디렉토리에서 실행
err := cmd.RunWithDir("ls", "/tmp")

// Shell 명령어를 특정 디렉토리에서 실행
err = cmd.RunShellWithDir("pwd && ls", "/tmp")

// PowerShell 명령어를 특정 디렉토리에서 실행
err = cmd.RunPowershellWithDir("Get-ChildItem", "C:\\temp")
```

### 커스텀 설정

```go
import (
    "os"
    "strings"
    "time"
)

cmd := easycmd.New(
    easycmd.WithStdIn(strings.NewReader("input")), // 표준 입력 설정
    easycmd.WithStdOut(os.Stdout),                 // 표준 출력 설정
    easycmd.WithStdErr(os.Stderr),                 // 표준 에러 설정
    easycmd.WithTimeoutSeconds(30),                // 타임아웃 설정 (30초)
    easycmd.WithEnv([]string{"VAR=value"}),        // 환경변수 설정
)

err := cmd.Run("cat") // 표준 입력에서 "input"을 읽어서 출력
```

### 디버그 모드

```go
import (
    "bytes"
    "fmt"
)

// 기본 디버그 모드 (디버그 출력은 stderr로)
cmd := easycmd.New(easycmd.WithDebug())
err := cmd.Run("echo hello world")

// 커스텀 디버그 출력 스트림
debugOut := &bytes.Buffer{}
cmd = easycmd.New(easycmd.WithDebug(debugOut))
err = cmd.Run("echo hello world")

fmt.Println("디버그 출력:", debugOut.String())
```

### 타임아웃 설정

```go
// 권장: 직관적인 API 사용
cmd := easycmd.New(easycmd.WithTimeoutSeconds(5))    // 5초
err := cmd.Run("sleep 3") // 정상 완료

cmd = easycmd.New(easycmd.WithTimeoutMillis(1500))   // 1.5초
err = cmd.Run("sleep 3")  // 1.5초 후 타임아웃

// time.Duration 직접 사용도 가능
import "time"
cmd = easycmd.New(easycmd.WithTimeout(5 * time.Second))
err = cmd.Run("sleep 10") // 5초 후 타임아웃

if err != nil {
    fmt.Printf("명령어 실행 실패: %v\n", err)
}
```

⚠️ **주의사항**: `WithTimeout()`에 숫자만 전달하면 나노초 단위로 해석됩니다!
```go
// ❌ 잘못된 사용법
cmd := easycmd.New(easycmd.WithTimeout(5))  // 5나노초 (거의 즉시 타임아웃)

// ✅ 올바른 사용법
cmd := easycmd.New(easycmd.WithTimeoutSeconds(5))       // 5초
cmd := easycmd.New(easycmd.WithTimeout(5 * time.Second)) // 5초
```

### 환경변수 설정

```go
cmd := easycmd.New(
    easycmd.WithEnv([]string{
        "MY_VAR=hello",
        "ANOTHER_VAR=world",
    }),
)

// 환경변수를 사용하는 쉘 명령어 실행
err := cmd.RunShell("echo $MY_VAR $ANOTHER_VAR")
```

### 복합 설정 사용

```go
import (
    "bytes"
    "strings"
    "time"
)

out := &bytes.Buffer{}
debugOut := &bytes.Buffer{}
input := strings.NewReader("test input\n")

cmd := easycmd.New(
    easycmd.WithStdIn(input),
    easycmd.WithStdOut(out),
    easycmd.WithStdErr(os.Stderr),
    easycmd.WithDebug(debugOut),
    easycmd.WithTimeoutSeconds(10),                 // 10초 타임아웃
    easycmd.WithEnv([]string{"LANG=ko_KR.UTF-8"}),
)

err := cmd.RunShell("cat && echo 'processing...' >&2")
```

## API 레퍼런스

### 주요 메서드

- `New(configApplies ...configApply) *Cmd`: 새로운 Cmd 인스턴스 생성
- `Run(commandStr string) error`: 기본 명령어 실행
- `RunShell(commandStr string) error`: bash로 래핑된 명령어 실행
- `RunPowershell(commandStr string) error`: PowerShell로 래핑된 명령어 실행
- `RunWithDir(commandStr string, runDirStr string) error`: 특정 디렉토리에서 기본 명령어 실행
- `RunShellWithDir(commandStr string, runDirStr string) error`: 특정 디렉토리에서 Shell 명령어 실행
- `RunPowershellWithDir(commandStr string, runDirStr string) error`: 특정 디렉토리에서 PowerShell 명령어 실행

### 설정 함수

- `WithStdIn(reader io.Reader) configApply`: 표준 입력 설정
- `WithStdOut(writer io.Writer) configApply`: 표준 출력 설정
- `WithStdErr(writer io.Writer) configApply`: 표준 에러 설정
- `WithDebug(debugOut ...io.Writer) configApply`: 디버그 모드 활성화 및 디버그 출력 스트림 설정
- `WithTimeout(timeout time.Duration) configApply`: 명령어 실행 타임아웃 설정 (time.Duration)
- `WithTimeoutSeconds(seconds int) configApply`: 명령어 실행 타임아웃 설정 (초 단위) ⭐ 권장
- `WithTimeoutMillis(millis int) configApply`: 명령어 실행 타임아웃 설정 (밀리초 단위) ⭐ 권장
- `WithEnv(env []string) configApply`: 환경변수 설정

#### 디버그 모드 출력 내용

디버그 모드가 활성화되면 다음 정보들이 출력됩니다:

- 파싱된 명령어 문자열
- 실행될 명령어 이름
- 명령어 인수 배열
- 실행 디렉토리 (설정된 경우)
- 타임아웃 설정 (설정된 경우)
- 환경변수 개수 (설정된 경우)
- 명령어 실행 시작/완료/실패 메시지
- 명령어 실행 시간 측정

## 에러 처리

```go
cmd := easycmd.New()
err := cmd.Run("invalid-command")
if err != nil {
    fmt.Printf("명령어 실행 실패: %v\n", err)
}

// 빈 명령어에 대한 특별한 에러
if errors.Is(err, easycmd.EmptyCmdError) {
    fmt.Println("빈 명령어입니다")
}
```

### 타임아웃 에러 유형

타임아웃 설정 시 발생할 수 있는 두 가지 에러 유형:

```go
cmd := easycmd.New(easycmd.WithTimeoutMillis(1)) // 매우 짧은 타임아웃
err := cmd.Run("sleep 1")

if err != nil {
    errMsg := err.Error()

    if strings.Contains(errMsg, "명령어 시작 실패") {
        fmt.Println("명령어 시작 전 타임아웃 발생")
    } else if strings.Contains(errMsg, "명령어 실행 타임아웃") {
        fmt.Println("명령어 실행 중 타임아웃 발생")
    }
}
```

## 라이선스

이 프로젝트는 Apache License 2.0 하에 배포됩니다. 자세한 내용은 [LICENSE](LICENSE) 파일을 참조하세요.
