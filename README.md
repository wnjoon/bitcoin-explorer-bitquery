# bitcoin-explorer-bitquery


## 요구사항

1. 현재의 비트코인 블록체인 정보를 먼저 Google 스프레드 시트에 옮길 것
2. Google 스프레드시트의 정보를 읽어 블록 생성 난이도 예측에 관한 데이터를 웹으로 시각화하여 제공
    - 가상의 Frontend Engineer가 해당 데이터를 전달 받아 테이블 및 시계열 차트를 구성할 예정 = 시각화에 필요한 자료 제공이 목적으로 판단

<br>

## 실행 방법

### 1. 앱 실행

```
$ go run cmd/main.go
```

### 2. 웹서버 주소

```
http://localhost:4000/api/v1/update?limit={개수}
```

### 3. 구글 스프레드시트 주소

```
# 작업 주소
https://docs.google.com/spreadsheets/d/1_1qFNF9ZqgydWgCJE_n5BsqfX1SNNYgmY7fL4iKk3EI/edit#gid=0

# 공유 주소
https://docs.google.com/spreadsheets/d/1_1qFNF9ZqgydWgCJE_n5BsqfX1SNNYgmY7fL4iKk3EI/edit?usp=sharing
```
- 인증정보 : [credentials.json](credentials.json)


<br>

## 설계

### 1. 언어 선택 이유

본 어플리케이션은 golang을 기반으로 작성되었다. 기존에는 Java를 이용하여로 하였으나 BitQuery를 다룰 수 있는 Java 기반의 API를 찾지 못하였고, 그 중 개인적으로 사용이 가능한 golang을 선정하였다.

### 2. 범위 선정

다양한 기능을 최대한 반영하여 개발하고자 하였으나, 시간이 한정적이고 이 중에서 과제를 할 수 있는 시간 또한 제한적이어서 가장 중요한 기능 순으로 간략하게 우선순위를 정해보았다.

|우선순위|기능|비고|
|------|---|---|
|1|BitQuery에서 비트코인 데이터를 가져오는 기능||
|1|BitQuery에서 가져온 데이터를 구조체 형태로 변경하는 기능|변경된 데이터를 조회가 가능한 데이터베이스에 넣기 위함|
|2|Google spreadsheet에 데이터를 넣는 기능||
|2|BitQuery에서 효과적으로 데이터를 가져올 수 있는 방안 확인||
|2|개발된 내역을 API 형태로 생성 후 호출||
|3|개발된 API에 보안(json 또는 기타 기반의 인증 과정) 적용||

> 시간 상 구글스프레드시트가 아닌 로컬에 데이터베이스와 유사한 기능을 구현하고 먼저 테스트하려 했으나, 생각보다 BitQuery와의 연동이 일찍 해결되어 구글스프레드시트에 직접 데이터를 주입하였다.

### 3.시스템 구성

- api : rest api 서버 구동
- cmd : main
- config : 앱(서버) 구동에 필요한 config 파일 및 이를 설정하기 위한 기능
    - BitQuery, Google 스프레드시트의 접근에 필요한 값들 포함
- pkg
    - explorer : bitQuery를 활용한 비트코인 모니터링 기능
    - storage : 데이터를 google spreadsheet에 주입하는 기능
- credential.json : Google 스프레드시트에 사용할 인증데이터
    - json을 직접 주입하지 않고, base64 형태로 값을 변환하여 사용(조금 더 나은 보안성 목적)
    - config.yaml -> storage.google.spreadsheet.credential

<br>

## 개발

### 1. BitQuery

현 시점에 가장 높은 블록 높이를 조회한 후, 원하는 갯수만큼 더 작은 값의 블록을 조회하도록 구현하였다. [금일 날짜(2022/11/05 기준)에 생성된 비트코인 정보](https://explorer.bitquery.io/bitcoin?from=2022-11-05&till=2022-11-05)의 하단에서 이에 대한 API를 제공받을 수 있었고, 이를 커스터마이즈해서 원하는 값을 얻었다.

BitQuery에서 데이터를 받아오는 방식은 아래와 같다.

1. 현 시점을 기준으로 가장 높은 값의 블록 번호를 받아온다 (getHigestBlock)
2. 해당 블록을 기준으로 사용자가 입력한 갯수(limit)만큼의 블록을 내림차순 순으로 가져온다 (UpdateBitcoinInfo)
3. 가져온 데이터는 구글 스프레드시트에 저장할 수 있도록 구조체화 시킨다.

이렇게 기능을 설계한 이유는 블록 생성 난이도는 계속 증가하기 때문에, 사용자가 2023년 1월 1일자의 블록 생성 난이도를 조회하려면 현 시점부터 지난 시점까지의 블록 생성 난이도 기록을 원하는 수 만큼 조회할 수 있어야 한다고 생각하였기 때문이다.

BitQuery에 대한 설정 정보는 config.yaml에서 수정 가능하다.

### 2. Google SpreadSheet

구글 스프레드시트에는 조회한 데이터를 그대로 저장하는 용도로만 사용한다. 사용자가 해당 데이터를 이용하여 테이블 또는 시계열 차트를 생성할 수 있도록 기초 데이터를 제공하는 용도 까지만 기능을 구현하였으며, 프론트엔드에서 발생할 수 있는 기능적인 부분은 제외하였다.

구글 스프레드시트는 아래와 같이 구성된다. 아래 3가지 정보는 모두 BitQuery에서 가져온다.

|Height|Difficulty|Timestamp|
|------|---|---|

- Height : 블록 높이
- Difficulty : 블록 난이도
- Timestamp : 블록 생성 시간

구글 스프레드시트에 대한 설정 정보는 config.yaml에서 수정 가능하다.

### 3. API Server

본 과제에서 사용자에게 제공되어야 하는 가장 중요한 기능은 **"<u>현 시점부터 사용자가 원하는 시점 까지의 최신 비트코인 블록에 대한 정보, 그 중에서도 블록 높이와 난이도, 생성된 시간</u>"** 이라고 판단하였다. 다른 추가 기능을 도출할 수도 있지만, 우선 주어진 기능을 먼저 개발하자는 취지로 시작하였고, 결국 해당 기능 하나만 존재한다.

```
http://localhost:4000/api/v1/update?limit={개수}
```

서버에 대한 접속 정보(포트)는 config.yaml에서 수정 가능하다.

<br>

## 향후 보완점

### 1. Go Test의 부재 

본 어플리케이션은 go test가 존재하지 않는다. 어플리케이션 개발에 사용한 개인 PC가 Mac으로 최근 Ventura로 업그레이드를 하고나서 go test 실행 시 아래와 같은 오류가 발생하였다.

```sh
xcrun: error: invalid active developer path
```

구글링([Mac 업그레이드 후 xcrun: error: invalid active developer path 에러 해결하기 - hahwul](https://www.hahwul.com/2019/11/18/how-to-fix-xcrun-error-after-macos-update/)) 결과 'xcode-select --install' 명령어를 통해 해결 가능하다고 하였고 이를 적용해보았으나, 어딘가 충돌이 존재하는지 해결되지 않았다 ☠️

우선 과제를 해결해야 하는 시간이 짧은 관계로, main에 직접 한땀한땀 기능을 실행하면서 앱을 개발하는 방향으로 진행하였다. 추후 시간이 된다면 해당 기능을 기반으로 더욱 다양한 테스트를 진행해봐야 할 것 같다.

### 2. 변수 및 함수 네이밍

초반에 업무 설계에 혼동이 있어서 중간에 한번 네이밍과 업무 설계 방안을 변경하였다. 그로 인해 조금 더 깔끔한 네이밍을 했으면 좋았겠다는 아쉬움이 있다. 향후에 조금 더 깔끔한 이름을 고안해야 할 것으로 보인다.

<br>



