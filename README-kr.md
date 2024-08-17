다음은 `README.md` 파일에 포함할 내용을 한국어로 작성되었습니다.

---

# gord - GORM을 활용한 제네릭 레포지토리 패턴

## 개요

`gord`는 Go 언어에서 GORM을 기반으로 구현된 제네릭 레포지토리 패턴입니다. 이 패키지는 공통 CRUD(Create, Read, Update, Delete) 작업을 위한 깔끔하고 재사용 가능한 인터페이스를 제공하여 Go 애플리케이션의 데이터 접근 레이어를 더 구조화되고 유지 보수하기 쉽게 만듭니다.

## 주요 기능

- **제네릭 CRUD 작업**: 모든 데이터 모델에서 작동하는 제네릭 CRUD 메서드 제공
- **타입 안전성**: Go의 제네릭 기능을 활용하여 모델과 ID에 대한 타입 안전성 보장
- **GORM 통합**: GORM의 강력한 기능을 그대로 활용
- **유연한 업데이트**: 타입 안전성을 갖춘 유연한 업데이트 가능
- **대량 작업 지원**: 대량 삭제 및 저장 작업 지원
- **존재 여부 확인**: ID로 레코드의 존재 여부를 확인하는 메서드 제공

## 설치

먼저 Go 프로젝트에 GORM을 설치합니다:

```bash
go get -u gorm.io/gorm
go get -u gorm.io/driver/sqlite # 사용 중인 데이터베이스 드라이버로 교체
```

그런 다음, 이 패키지를 프로젝트에 포함시킵니다:

```bash
go get github.com/sangjinsu/gord
```

## 사용법

### 1. 모델 정의

먼저 레포지토리로 관리할 모델을 정의합니다.

```go
package models

import "gorm.io/gorm"

type User struct {
    gorm.Model
    Name  string
    Email string
}
```

### 2. 레포지토리 초기화

해당 모델에 대한 레포지토리 인스턴스를 생성합니다.

```go
package main

import (
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "github.com/yourusername/gord"
    "yourapp/models"
)

func main() {
    db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    if err != nil {
        panic("데이터베이스 연결에 실패했습니다.")
    }

    db.AutoMigrate(&models.User{})

    userRepository := gord.Repository[models.User, uint]{tx: db}

    // userRepository를 사용해 CRUD 작업을 수행합니다.
}
```

### 3. CRUD 작업 수행

이제 레포지토리 인스턴스를 사용하여 다양한 작업을 수행할 수 있습니다.

```go
// 새로운 사용자 생성
user := models.User{Name: "John Doe", Email: "john.doe@example.com"}
err := userRepository.Create(user)

// ID로 사용자 찾기
foundUser, err := userRepository.FindByID(1)

// 사용자 업데이트
updates := gord.UpdateMap{
    "Name": "Jane Doe",
}
err := userRepository.Updates(foundUser, updates)

// 사용자 삭제
err := userRepository.Delete(foundUser)

// 사용자 존재 여부 확인
exists, err := userRepository.ExistByID(1)
```

### 4. 레포지토리 커스터마이징

커스텀 메서드가 필요한 경우, `Repository` 구조체를 임베드하여 메서드를 추가할 수 있습니다.

```go
type UserRepository struct {
    gord.Repository[models.User, uint]
}

func (r UserRepository) FindByEmail(email string) (models.User, error) {
    var user models.User
    err := r.tx.Where("email = ?", email).First(&user).Error
    return user, err
}
```

### 5. 업데이트 검증

`gord`는 업데이트 작업에서 허용된 타입만 사용되도록 검증하는 메커니즘을 제공합니다.

```go
updates := gord.UpdateMap{
    "Name": "New Name",
    // "InvalidField": []int{1, 2, 3}, // 이 경우 검증 오류가 발생합니다.
}

if err := updates.valid(); err != nil {
    fmt.Println("검증 실패:", err)
} else {
    userRepository.Updates(user, updates)
}
```

## 기여 방법

기여는 언제나 환영입니다! 풀 리퀘스트를 제출하거나 아이디어를 논의할 이슈를 열어주세요.

## 라이선스

이 프로젝트는 MIT 라이선스 하에 라이선스됩니다. 자세한 내용은 `LICENSE` 파일을 참조하세요.
