# План рефакторинга интеграции Sentry

**Дата создания:** 2025-11-28
**Цель:** Добавить stack traces для всех ошибок, исправить недочеты интеграции Sentry

---

## Проблемы текущей интеграции

### Критические:
1. ❌ **Отсутствие stack traces** - не используется `pkg/errors`, все ошибки создаются через `fmt.Errorf`/`errors.New`
2. ❌ **Бесполезная функция BeforeSend** - просто возвращает event без изменений
3. ❌ **Нет дедупликации stack traces** - при многократном обёртывании дублируются frames
4. ❌ **Неправильные skip значения** - hardcoded значение `skip=4` может быть неточным

### Важные:
5. ⚠️ **Хардкод context keys** - magic strings `"request_id"`, `"user_id"` без type safety
6. ⚠️ **grpc/sentry.go создаёт ошибки через fmt.Errorf** - паники не имеют stack trace
7. ⚠️ **Разные timeout значения** - 2s и 5s в разных местах
8. ⚠️ **Magic numbers в convertLevel** - значения `12`, `11` без констант

### Улучшения:
9. 💡 **Отсутствие breadcrumbs** - не используется функционал Sentry для отслеживания событий
10. 💡 **Дублирующийся код** - повторяющаяся логика recover в grpc interceptors
11. 💡 **Проблемы в extractStackTrace** - хардкод типов ошибок, нет обработки wrapped errors

---

## Этапы рефакторинга

### ЭТАП 1: Создание централизованного пакета errors ⏳

**Цель:** Wrapper над `pkg/errors` с дедупликацией stack trace

**Файлы для создания:**
- `errors/errors.go` - основной файл с функциями New, Wrap, Wrapf
- `errors/dedup.go` - логика дедупликации stack trace
- `errors/sentinel.go` - sentinel errors (опционально)
- `errors/errors_test.go` - unit тесты

**Что реализуем:**
```go
package errors

// Re-export базовых функций из pkg/errors
func New(message string) error
func Wrap(err error, message string) error
func Wrapf(err error, format string, args ...interface{}) error
func Cause(err error) error
func Errorf(format string, args ...interface{}) error

// Новая функция с дедупликацией (из статьи incident.io)
func WrapWithDedup(err error, message string) error

// Проверки для совместимости с Go 1.13+
func Is(err, target error) bool
func As(err error, target interface{}) bool
```

**Логика дедупликации:**
- Проверяем, есть ли уже stack trace в ошибке
- Сравниваем текущий stack trace с существующим через сопоставление префиксов
- Добавляем новый trace только если он не является предком существующего

**Acceptance criteria:**
- ✅ Unit тесты покрывают 80%+ кода
- ✅ WrapWithDedup корректно дедуплицирует stack traces
- ✅ Совместимость с pkg/errors и стандартными errors

---

### ЭТАП 2: Typed context keys ⏳

**Цель:** Type-safe извлечение данных из контекста

**Файл:** `context/keys.go`

```go
package context

import "context"

type contextKey string

const (
    RequestIDKey contextKey = "request_id"
    UserIDKey    contextKey = "user_id"
)

// WithRequestID добавляет request ID в контекст
func WithRequestID(ctx context.Context, requestID string) context.Context {
    return context.WithValue(ctx, RequestIDKey, requestID)
}

// GetRequestID извлекает request ID из контекста
func GetRequestID(ctx context.Context) (string, bool) {
    requestID, ok := ctx.Value(RequestIDKey).(string)
    return requestID, ok
}

// WithUserID добавляет user ID в контекст
func WithUserID(ctx context.Context, userID string) context.Context {
    return context.WithValue(ctx, UserIDKey, userID)
}

// GetUserID извлекает user ID из контекста
func GetUserID(ctx context.Context) (string, bool) {
    userID, ok := ctx.Value(UserIDKey).(string)
    return userID, ok
}
```

**Acceptance criteria:**
- ✅ Type-safe API для работы с context
- ✅ Документация для всех публичных функций
- ✅ Unit тесты

---

### ЭТАП 3: Рефакторинг sentry/sentry.go ⏳

#### 3.1. Улучшить extractStackTrace

**Изменения:**
- Добавить обработку wrapped errors через `errors.Unwrap`
- Удалить хардкод типов ошибок (`*errors.errorString`, `*errors.fundamental`)
- Улучшить извлечение root cause
- Поддержка дедуплицированных ошибок из `core/errors`

```go
func extractStackTrace(err error) []sentrySDK.Exception {
    type stackTracer interface {
        StackTrace() errors.StackTrace
    }

    // Получаем тип ошибки
    errType := fmt.Sprintf("%T", err)
    if strings.HasPrefix(errType, "*errors.") {
        errType = "error"
    }

    exception := sentrySDK.Exception{
        Value: err.Error(),
        Type:  errType,
    }

    // Пробуем извлечь stack trace из любого уровня wrapped errors
    current := err
    for current != nil {
        if st, ok := current.(stackTracer); ok {
            exception.Stacktrace = convertStackTrace(st.StackTrace())
            break
        }
        current = errors.Unwrap(current)
    }

    return []sentrySDK.Exception{exception}
}
```

#### 3.2. Убрать/улучшить BeforeSend

**Варианты:**
1. Удалить полностью (если не нужна фильтрация)
2. Добавить полезную логику:

```go
BeforeSend: func(event *sentrySDK.Event, hint *sentrySDK.EventHint) *sentrySDK.Event {
    // Фильтрация чувствительных данных
    sensitiveKeys := []string{"password", "token", "secret", "api_key"}

    for _, key := range sensitiveKeys {
        delete(event.Extra, key)
        delete(event.Tags, key)
    }

    // Можно добавить дополнительную логику
    return event
},
```

#### 3.3. Использовать typed context keys

```go
// Вместо:
if requestID, ok := ctx.Value("request_id").(string); ok {
    event.Tags["request_id"] = requestID
}

// Использовать:
import ctxkeys "github.com/Rasikrr/core/context"

if requestID, ok := ctxkeys.GetRequestID(ctx); ok {
    event.Tags["request_id"] = requestID
}
if userID, ok := ctxkeys.GetUserID(ctx); ok {
    event.User = sentrySDK.User{ID: userID}
}
```

#### 3.4. Добавить константы для уровней логирования

```go
import "github.com/Rasikrr/core/log"

func convertLevel(level slog.Level) sentrySDK.Level {
    switch {
    case level >= log.LevelFatal:
        return sentrySDK.LevelFatal
    case level >= log.LevelSentry:
        return sentrySDK.LevelError
    case level >= slog.LevelError:
        return sentrySDK.LevelError
    case level >= slog.LevelWarn:
        return sentrySDK.LevelWarning
    case level >= slog.LevelInfo:
        return sentrySDK.LevelInfo
    default:
        return sentrySDK.LevelDebug
    }
}
```

**Acceptance criteria:**
- ✅ extractStackTrace обрабатывает wrapped errors
- ✅ Нет hardcoded строк
- ✅ Используются typed context keys
- ✅ Используются константы уровней из log пакета

---

### ЭТАП 4: Исправить captureStacktrace ⏳

**Проблема:** `skip=4` - магическое число, может быть неточным

**Решение 1:** Более умное определение skip

```go
func captureStacktrace(skip int) *sentrySDK.Stacktrace {
    const maxFrames = 32
    pcs := make([]uintptr, maxFrames)
    n := runtime.Callers(2, pcs) // Всегда начинаем с 2

    if n == 0 {
        return nil
    }

    frames := make([]sentrySDK.Frame, 0, n)
    skipped := 0

    for i := 0; i < n; i++ {
        pc := pcs[i]
        fn := runtime.FuncForPC(pc)
        if fn == nil {
            continue
        }

        // Пропускаем frames из sentry и log пакетов
        funcName := fn.Name()
        if skipped < skip ||
           strings.Contains(funcName, "github.com/Rasikrr/core/sentry") ||
           strings.Contains(funcName, "github.com/Rasikrr/core/log") {
            if !strings.Contains(funcName, "github.com/Rasikrr/core/sentry") &&
               !strings.Contains(funcName, "github.com/Rasikrr/core/log") {
                skipped++
            }
            continue
        }

        file, line := fn.FileLine(pc)
        frames = append(frames, sentrySDK.Frame{
            Function: funcName,
            Filename: file,
            Lineno:   line,
            InApp:    !strings.Contains(funcName, "runtime.") &&
                      !strings.Contains(funcName, "/vendor/"),
        })
    }

    // Sentry ожидает frames в обратном порядке
    for i, j := 0, len(frames)-1; i < j; i, j = i+1, j-1 {
        frames[i], frames[j] = frames[j], frames[i]
    }

    return &sentrySDK.Stacktrace{
        Frames: frames,
    }
}
```

**Решение 2:** Unit тест для проверки правильности skip

```go
func TestCaptureStacktrace(t *testing.T) {
    stack := captureStacktrace(0)
    require.NotNil(t, stack)
    require.NotEmpty(t, stack.Frames)

    // Первый frame должен быть из этой тестовой функции
    firstFrame := stack.Frames[0]
    assert.Contains(t, firstFrame.Function, "TestCaptureStacktrace")
}
```

**Acceptance criteria:**
- ✅ captureStacktrace корректно пропускает internal frames
- ✅ Unit тест подтверждает правильность работы
- ✅ InApp флаг корректно проставлен

---

### ЭТАП 5: Рефакторинг grpc/sentry.go ⏳

#### 5.1. Использовать errors.Wrap вместо fmt.Errorf

```go
// Вместо:
panicErr = fmt.Errorf("panic: %v", r)

// Использовать:
import "github.com/Rasikrr/core/errors"

panicErr = errors.Errorf("panic: %v", r)
```

#### 5.2. Вынести общую логику recover

```go
func handlePanic(hub *sentrySDK.Hub, r interface{}, tags map[string]string) {
    hub.WithScope(func(scope *sentrySDK.Scope) {
        scope.SetLevel(sentrySDK.LevelFatal)

        for key, value := range tags {
            scope.SetTag(key, value)
        }

        scope.SetContext("panic", map[string]interface{}{
            "value":      fmt.Sprintf("%v", r),
            "stacktrace": string(debug.Stack()),
        })

        var panicErr error
        switch x := r.(type) {
        case error:
            panicErr = x
        default:
            panicErr = errors.Errorf("panic: %v", r)
        }

        hub.CaptureException(panicErr)
    })
}

// Использование:
defer func() {
    if r := recover(); r != nil {
        handlePanic(hub, r, map[string]string{
            "grpc.type":    "unary",
            "grpc.service": service,
            "grpc.method":  method,
        })
        panic(r) // re-panic
    }
}()
```

#### 5.3. Централизовать timeouts

```go
// В sentry/config.go или sentry/constants.go
package sentry

const (
    // FlushTimeout - максимальное время ожидания отправки событий в Sentry
    FlushTimeout = 5 * time.Second
)
```

Использование:
```go
sentrySDK.Flush(sentry.FlushTimeout)
```

**Acceptance criteria:**
- ✅ Все ошибки создаются через `core/errors`
- ✅ Нет дублирующегося кода recover
- ✅ Timeouts вынесены в константы

---

### ЭТАП 6: Добавить breadcrumbs ⏳

**Цель:** Отслеживание событий, приведших к ошибке

**Файл:** `sentry/breadcrumbs.go`

```go
package sentry

import (
    "context"
    "time"

    sentrySDK "github.com/getsentry/sentry-go"
)

// AddBreadcrumb добавляет breadcrumb в текущий Sentry hub.
// Breadcrumbs помогают отследить события, приведшие к ошибке.
func AddBreadcrumb(ctx context.Context, category, message string, level sentrySDK.Level, data map[string]interface{}) {
    if !Enabled() {
        return
    }

    hub := sentrySDK.GetHubFromContext(ctx)
    if hub == nil {
        hub = sentrySDK.CurrentHub()
    }

    hub.AddBreadcrumb(&sentrySDK.Breadcrumb{
        Category:  category,
        Message:   message,
        Level:     level,
        Data:      data,
        Timestamp: time.Now(),
    }, nil)
}

// Convenience функции для частых случаев

func AddDatabaseBreadcrumb(ctx context.Context, query string, duration time.Duration) {
    AddBreadcrumb(ctx, "database", "Query executed", sentrySDK.LevelInfo, map[string]interface{}{
        "query":    query,
        "duration": duration.String(),
    })
}

func AddHTTPBreadcrumb(ctx context.Context, method, url string, statusCode int) {
    AddBreadcrumb(ctx, "http", "HTTP request", sentrySDK.LevelInfo, map[string]interface{}{
        "method":      method,
        "url":         url,
        "status_code": statusCode,
    })
}

func AddGRPCBreadcrumb(ctx context.Context, service, method string) {
    AddBreadcrumb(ctx, "grpc", "gRPC call", sentrySDK.LevelInfo, map[string]interface{}{
        "service": service,
        "method":  method,
    })
}
```

**Использование:**

```go
// В database запросах
func (r *Repository) FindUser(ctx context.Context, id string) (*User, error) {
    start := time.Now()
    query := "SELECT * FROM users WHERE id = $1"

    // ... выполнение запроса ...

    sentry.AddDatabaseBreadcrumb(ctx, query, time.Since(start))
    return user, nil
}

// В HTTP handlers
func HandleRequest(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    // ... обработка ...

    sentry.AddHTTPBreadcrumb(ctx, r.Method, r.URL.String(), statusCode)
}
```

**Acceptance criteria:**
- ✅ API для добавления breadcrumbs
- ✅ Convenience функции для частых случаев
- ✅ Документация с примерами использования

---

### ЭТАП 7: Замена fmt.Errorf на errors.Wrap ⏳

**Цель:** Заменить все 74+ вхождения `fmt.Errorf`/`errors.New` на `core/errors`

**План:**
1. Найти все файлы с `fmt.Errorf` и `errors.New`
2. Систематически пройти по каждому файлу
3. Применить паттерны замены

**Паттерны замены:**

```go
// Паттерн 1: fmt.Errorf с %w (wrapping)
fmt.Errorf("failed to connect: %w", err)
→
errors.Wrap(err, "failed to connect")

// Паттерн 2: fmt.Errorf с %v (не рекомендуется, но встречается)
fmt.Errorf("failed to connect: %v", err)
→
errors.Wrap(err, "failed to connect")

// Паттерн 3: fmt.Errorf с форматированием без wrapping
fmt.Errorf("invalid input: %s", input)
→
errors.Errorf("invalid input: %s", input)

// Паттерн 4: fmt.Errorf без параметров
fmt.Errorf("something went wrong")
→
errors.New("something went wrong")

// Паттерн 5: errors.New из стандартной библиотеки
errors.New("something went wrong")
→
errors.New("something went wrong") // теперь из core/errors

// Паттерн 6: Sentinel errors (объявления на уровне пакета)
var ErrNotFound = errors.New("not found")
→
var ErrNotFound = errors.New("not found") // из core/errors
```

**Процесс:**

```bash
# 1. Найти все файлы с fmt.Errorf
grep -r "fmt\.Errorf" --include="*.go" .

# 2. Найти все файлы с errors.New из стандартной библиотеки
grep -r "errors\.New" --include="*.go" .

# 3. Для каждого файла:
#    - Добавить импорт "github.com/Rasikrr/core/errors"
#    - Удалить импорт стандартного "errors" (если больше не используется)
#    - Заменить вызовы согласно паттернам
```

**Файлы для замены (по приоритету):**

1. **Критичные (бизнес-логика):**
   - `application/*.go`
   - `database/*.go`
   - `api/*.go`

2. **Инфраструктурные:**
   - `grpc/*.go`
   - `http/*.go`
   - `brokers/*.go`

3. **Утилиты и конфиги:**
   - `config/*.go`
   - `redis/*.go`
   - `metrics/*.go`

**Acceptance criteria:**
- ✅ Все `fmt.Errorf` с wrapping заменены на `errors.Wrap`
- ✅ Все `fmt.Errorf` без wrapping заменены на `errors.Errorf`
- ✅ Все `errors.New` используют `core/errors`
- ✅ Удалены неиспользуемые импорты стандартного `errors`
- ✅ Код компилируется без ошибок
- ✅ Все тесты проходят

---

### ЭТАП 8: Настройка golangci-lint ⏳

**Цель:** Автоматический enforcement использования `core/errors`

**Файл:** `.golangci.yml`

```yaml
linters-settings:
  # Запрет прямых импортов errors и pkg/errors
  gomodguard:
    blocked:
      modules:
        - errors:
            recommendations:
              - github.com/Rasikrr/core/errors
            reason: "Use core/errors package instead for automatic stack traces"
        - github.com/pkg/errors:
            recommendations:
              - github.com/Rasikrr/core/errors
            reason: "Use core/errors wrapper for stack trace deduplication"

  # Запрет использования fmt.Errorf
  forbidigo:
    forbid:
      - p: 'fmt\.Errorf'
        msg: 'Use errors.Wrap, errors.Wrapf or errors.Errorf from github.com/Rasikrr/core/errors'
      - p: '^errors\.New'
        pkg: '^errors$'
        msg: 'Use errors.New from github.com/Rasikrr/core/errors package'
      - p: '^errors\.Wrap'
        pkg: '^errors$'
        msg: 'Use errors.Wrap from github.com/Rasikrr/core/errors package'

  # Проверка правильного использования error wrapping
  errorlint:
    errorf: true
    asserts: true
    comparison: true

linters:
  enable:
    - gomodguard
    - forbidigo
    - errorlint
    - errcheck

  disable:
    # Отключаем конфликтующие линтеры, если есть
```

**Проверка:**

```bash
# Запустить линтер
golangci-lint run

# Должны увидеть ошибки, если где-то используется запрещенный код
```

**Acceptance criteria:**
- ✅ golangci-lint настроен
- ✅ Запрет импорта стандартного `errors` работает
- ✅ Запрет импорта `pkg/errors` работает
- ✅ Запрет `fmt.Errorf` работает
- ✅ CI/CD проверяет линтер

---

### ЭТАП 9: Написать unit тесты ⏳

**Цель:** Покрыть новый код тестами на 80%+

#### errors/errors_test.go

```go
package errors_test

import (
    "testing"

    "github.com/Rasikrr/core/errors"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
    err := errors.New("test error")
    assert.EqualError(t, err, "test error")

    // Проверяем, что есть stack trace
    type stackTracer interface {
        StackTrace() errors.StackTrace
    }
    _, ok := err.(stackTracer)
    assert.True(t, ok, "error should have stack trace")
}

func TestWrap(t *testing.T) {
    original := errors.New("original error")
    wrapped := errors.Wrap(original, "wrapped")

    assert.EqualError(t, wrapped, "wrapped: original error")
    assert.True(t, errors.Is(wrapped, original))
}

func TestWrapWithDedup(t *testing.T) {
    // Создаем ошибку с stack trace
    err1 := errors.New("original")

    // Оборачиваем несколько раз
    err2 := errors.WrapWithDedup(err1, "wrap 1")
    err3 := errors.WrapWithDedup(err2, "wrap 2")

    // TODO: Проверить, что stack trace не дублируется
    // Это зависит от реализации дедупликации
}

func TestErrorf(t *testing.T) {
    err := errors.Errorf("error with number: %d", 42)
    assert.EqualError(t, err, "error with number: 42")
}

func TestCause(t *testing.T) {
    original := errors.New("original")
    wrapped := errors.Wrap(original, "wrapped")

    cause := errors.Cause(wrapped)
    assert.Equal(t, original, cause)
}

func TestIs(t *testing.T) {
    original := errors.New("original")
    wrapped := errors.Wrap(original, "wrapped")

    assert.True(t, errors.Is(wrapped, original))
    assert.False(t, errors.Is(wrapped, errors.New("other")))
}

func TestAs(t *testing.T) {
    type customError struct {
        Code int
    }

    func (e *customError) Error() string {
        return fmt.Sprintf("error with code %d", e.Code)
    }

    original := &customError{Code: 42}
    wrapped := errors.Wrap(original, "wrapped")

    var target *customError
    assert.True(t, errors.As(wrapped, &target))
    assert.Equal(t, 42, target.Code)
}
```

#### sentry/sentry_test.go

```go
package sentry_test

import (
    "context"
    "log/slog"
    "testing"

    "github.com/Rasikrr/core/errors"
    "github.com/Rasikrr/core/sentry"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestExtractStackTrace(t *testing.T) {
    // Тестируем с ошибкой из pkg/errors
    err := errors.New("test error")

    exceptions := extractStackTrace(err)
    require.Len(t, exceptions, 1)

    ex := exceptions[0]
    assert.Equal(t, "test error", ex.Value)
    assert.NotNil(t, ex.Stacktrace)
    assert.NotEmpty(t, ex.Stacktrace.Frames)
}

func TestExtractStackTraceWrapped(t *testing.T) {
    original := errors.New("original")
    wrapped := errors.Wrap(original, "wrapped")

    exceptions := extractStackTrace(wrapped)
    require.Len(t, exceptions, 1)

    assert.Contains(t, exceptions[0].Value, "wrapped")
    assert.NotNil(t, exceptions[0].Stacktrace)
}

func TestCaptureStacktrace(t *testing.T) {
    stack := captureStacktrace(0)
    require.NotNil(t, stack)
    require.NotEmpty(t, stack.Frames)

    // Первый frame должен быть из этой функции
    firstFrame := stack.Frames[0]
    assert.Contains(t, firstFrame.Function, "TestCaptureStacktrace")
}

func TestConvertLevel(t *testing.T) {
    tests := []struct {
        name     string
        level    slog.Level
        expected sentrySDK.Level
    }{
        {"Fatal", log.LevelFatal, sentrySDK.LevelFatal},
        {"Sentry", log.LevelSentry, sentrySDK.LevelError},
        {"Error", slog.LevelError, sentrySDK.LevelError},
        {"Warn", slog.LevelWarn, sentrySDK.LevelWarning},
        {"Info", slog.LevelInfo, sentrySDK.LevelInfo},
        {"Debug", slog.LevelDebug, sentrySDK.LevelDebug},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := convertLevel(tt.level)
            assert.Equal(t, tt.expected, result)
        })
    }
}

func TestCaptureEventWithError(t *testing.T) {
    // Требует инициализации Sentry
    // Можно использовать mock или test DSN
}
```

#### context/keys_test.go

```go
package context_test

import (
    "context"
    "testing"

    ctxkeys "github.com/Rasikrr/core/context"
    "github.com/stretchr/testify/assert"
)

func TestRequestID(t *testing.T) {
    ctx := context.Background()

    // Initially no request ID
    _, ok := ctxkeys.GetRequestID(ctx)
    assert.False(t, ok)

    // Set request ID
    ctx = ctxkeys.WithRequestID(ctx, "test-request-123")

    // Get request ID
    requestID, ok := ctxkeys.GetRequestID(ctx)
    assert.True(t, ok)
    assert.Equal(t, "test-request-123", requestID)
}

func TestUserID(t *testing.T) {
    ctx := context.Background()

    // Initially no user ID
    _, ok := ctxkeys.GetUserID(ctx)
    assert.False(t, ok)

    // Set user ID
    ctx = ctxkeys.WithUserID(ctx, "user-456")

    // Get user ID
    userID, ok := ctxkeys.GetUserID(ctx)
    assert.True(t, ok)
    assert.Equal(t, "user-456", userID)
}
```

**Запуск тестов:**

```bash
# Запустить все тесты
go test ./...

# С покрытием
go test -cover ./...

# Детальное покрытие
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

**Acceptance criteria:**
- ✅ Покрытие тестами 80%+ для errors пакета
- ✅ Покрытие тестами 70%+ для sentry пакета
- ✅ Все тесты проходят
- ✅ CI/CD запускает тесты автоматически

---

### ЭТАП 10: Интеграционное тестирование ⏳

**Цель:** Проверить работу в реальных условиях

#### 10.1. Создать test endpoint

**Файл:** `http/handlers/sentry_test.go` (или в соответствующем месте)

```go
package handlers

import (
    "net/http"
    "time"

    "github.com/Rasikrr/core/errors"
    "github.com/Rasikrr/core/log"
    "github.com/Rasikrr/core/sentry"
)

// TestSentryHandler - эндпоинт для тестирования Sentry интеграции
// GET /test/sentry?type=simple|wrapped|panic|breadcrumbs
func TestSentryHandler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    testType := r.URL.Query().Get("type")

    switch testType {
    case "simple":
        // Тест 1: Простая ошибка
        err := errors.New("simple test error")
        log.Sentry(ctx, "Test: simple error", log.Attr{Key: "error", Value: err})
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Simple error sent to Sentry"))

    case "wrapped":
        // Тест 2: Обёрнутая ошибка
        err := simulateDeepError()
        log.Sentry(ctx, "Test: wrapped error", log.Attr{Key: "error", Value: err})
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Wrapped error sent to Sentry"))

    case "panic":
        // Тест 3: Паника (должна быть обработана middleware)
        panic("test panic from handler")

    case "breadcrumbs":
        // Тест 4: Ошибка с breadcrumbs
        sentry.AddBreadcrumb(ctx, "test", "Step 1: Started processing", sentrySDK.LevelInfo, nil)
        time.Sleep(100 * time.Millisecond)

        sentry.AddBreadcrumb(ctx, "test", "Step 2: Fetching data", sentrySDK.LevelInfo, nil)
        time.Sleep(100 * time.Millisecond)

        sentry.AddBreadcrumb(ctx, "test", "Step 3: Processing failed", sentrySDK.LevelWarning, nil)

        err := errors.New("error with breadcrumbs")
        log.Sentry(ctx, "Test: error with breadcrumbs", log.Attr{Key: "error", Value: err})
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Error with breadcrumbs sent to Sentry"))

    case "multiple_wrap":
        // Тест 5: Многократное обёртывание (проверка дедупликации)
        err := errors.New("base error")
        err = errors.Wrap(err, "layer 1")
        err = errors.Wrap(err, "layer 2")
        err = errors.Wrap(err, "layer 3")
        log.Sentry(ctx, "Test: multiple wrapping", log.Attr{Key: "error", Value: err})
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Multiple wrapped error sent to Sentry"))

    default:
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Unknown test type. Use: simple|wrapped|panic|breadcrumbs|multiple_wrap"))
    }
}

func simulateDeepError() error {
    err := level3Error()
    return errors.Wrap(err, "level 2 failed")
}

func level3Error() error {
    err := level4Error()
    return errors.Wrap(err, "level 3 failed")
}

func level4Error() error {
    return errors.New("level 4: database connection failed")
}
```

#### 10.2. Регистрация test endpoint

```go
// В вашем router setup (только для dev/staging!)
if env != enum.EnvironmentProd {
    router.HandleFunc("/test/sentry", handlers.TestSentryHandler).Methods("GET")
}
```

#### 10.3. Чек-лист проверки в Sentry UI

После запуска приложения и вызова test endpoints:

**Для типа "simple":**
- ✅ Событие появляется в Sentry
- ✅ Stack trace присутствует
- ✅ Stack trace показывает правильные file:line
- ✅ Нет фреймов из runtime/internal пакетов
- ✅ Сообщение: "simple test error"

**Для типа "wrapped":**
- ✅ Stack trace показывает всю цепочку вызовов
- ✅ Видны все уровни: level 4 → level 3 → level 2
- ✅ Нет дублирования одинаковых frames
- ✅ Сообщение содержит весь путь обёртывания

**Для типа "panic":**
- ✅ Паника перехвачена middleware
- ✅ Stack trace содержит место паники
- ✅ Level = Fatal
- ✅ Контекст panic сохранён

**Для типа "breadcrumbs":**
- ✅ В событии есть секция "Breadcrumbs"
- ✅ Видны все 3 шага: Step 1, Step 2, Step 3
- ✅ Временные метки корректны
- ✅ Помогает понять последовательность событий

**Для типа "multiple_wrap":**
- ✅ Stack trace НЕ дублируется
- ✅ Видно только уникальные frames
- ✅ Сообщение показывает всю цепочку wrapping

#### 10.4. Curl команды для тестирования

```bash
# Тест 1: Simple error
curl http://localhost:8080/test/sentry?type=simple

# Тест 2: Wrapped error
curl http://localhost:8080/test/sentry?type=wrapped

# Тест 3: Panic
curl http://localhost:8080/test/sentry?type=panic

# Тест 4: Breadcrumbs
curl http://localhost:8080/test/sentry?type=breadcrumbs

# Тест 5: Multiple wrapping
curl http://localhost:8080/test/sentry?type=multiple_wrap
```

#### 10.5. Проверка в production-like окружении

1. Развернуть на staging окружении
2. Включить Sentry с реальным DSN
3. Запустить тесты
4. Проверить в Sentry UI
5. Проверить метрики производительности (overhead от stack traces)

**Acceptance criteria:**
- ✅ Все 5 типов тестов проходят успешно
- ✅ Stack traces корректны в Sentry UI
- ✅ Нет дублирования frames
- ✅ Breadcrumbs работают
- ✅ Overhead производительности приемлем (<5ms на error)
- ✅ Test endpoints удалены или отключены для production

---

## Метрики успеха

После завершения всех этапов должны достичь:

### В Sentry UI:
- ✅ Stack traces для 100% ошибок
- ✅ Правильные file:line номера
- ✅ Нет дублирования frames
- ✅ Breadcrumbs перед ошибками
- ✅ Теги request_id и user_id присутствуют
- ✅ Правильные severity levels

### В коде:
- ✅ 0 вхождений `fmt.Errorf` для создания errors
- ✅ 0 прямых импортов стандартного `errors`
- ✅ 0 прямых импортов `pkg/errors`
- ✅ Все ошибки через `github.com/Rasikrr/core/errors`
- ✅ golangci-lint проходит без ошибок
- ✅ Нет magic numbers и magic strings

### Качество:
- ✅ Unit тесты покрывают 80%+ errors пакета
- ✅ Unit тесты покрывают 70%+ sentry пакета
- ✅ Интеграционные тесты проходят
- ✅ Документация обновлена
- ✅ CI/CD pipeline включает все проверки

---

## Порядок выполнения

**Рекомендованный порядок (последовательно):**

1. ✅ **ЭТАП 1** - Создать core/errors (фундамент)
2. ✅ **ЭТАП 2** - Создать typed context keys
3. ✅ **ЭТАП 3** - Рефакторинг sentry/sentry.go
4. ✅ **ЭТАП 4** - Исправить captureStacktrace
5. ✅ **ЭТАП 5** - Рефакторинг grpc/sentry.go
6. ✅ **ЭТАП 6** - Добавить breadcrumbs
7. ✅ **ЭТАП 7** - Массовая замена fmt.Errorf (самый долгий)
8. ✅ **ЭТАП 8** - Настроить линтер
9. ✅ **ЭТАП 9** - Написать тесты
10. ✅ **ЭТАП 10** - Интеграционное тестирование

**Можно параллелить:**
- ЭТАП 1 и ЭТАП 2 независимы
- ЭТАП 6 (breadcrumbs) можно делать параллельно с ЭТАП 7
- ЭТАП 9 (тесты) писать постепенно после каждого этапа

---

## Риски и митигация

### Риск 1: Большой объём изменений в ЭТАП 7
**Митигация:** Разбить на батчи по директориям, делать постепенно

### Риск 2: Breaking changes для существующего кода
**Митигация:** Создать core/errors как wrapper, не меняя API

### Риск 3: Performance overhead от stack traces
**Митигация:** Измерить performance до/после, оптимизировать если нужно

### Риск 4: Конфликты при параллельной разработке
**Митигация:** Делать рефакторинг в отдельной feature ветке, frequent rebase

---

## Ссылки

- [Статья incident.io про Golang errors](https://incident.io/blog/golang-errors)
- [pkg/errors GitHub](https://github.com/pkg/errors)
- [Sentry Go SDK docs](https://docs.sentry.io/platforms/go/)
- [Go error handling best practices](https://go.dev/blog/go1.13-errors)

---

**Создано:** 2025-11-28
**Версия:** 1.0
**Статус:** Готово к выполнению