# Introduction

**onsengo**♨(音泉go) is a CLI program that allows listing and browsing radio shows on https://onsen.ag/. (音泉)

It's written in Go, including a data wrapper library with command frontends that use it.

Onsen.ag website, since last refactoring, uses a server-side rendering which is very friendly to developers,
we could get all the data with only one fetch. So the concept is easy:

1. Every time when it needs the data, it requests to onsen.ag.
2. Run the obfuscated Js code ([dop251/goja](https://github.com/dop251/goja)) to get the data json.
3. Parses data and creates a decorator to manipulate with.
4. Cmd is implemented with [spf13/cobra](https://github.com/spf13/cobra), and [adios/pprint](https://github.com/adios/pprint) handles boilerplate typesetting.

**onsengo**♨ currently supports the following commands:

* `onsengo ls`

## Global options

`--backend`: set a custom website url. If you have previously index.html archives, you can provide it like this:
```
onsengo ls --backend file:///full/path/to/index.html
```
`--session`: if you are a premium user, you could provide a **session** to the command:
```
onsengo ls --session SESSION_STRING_KEEP_IT_SECURE
```
A session string can be found in the cookie when you get logged-in and starts browsing onsen.ag.
- Open a browser (Firefox as an example) to https://onsen.ag/, make sure you are logged in.
- Open developer tool, reload the page
- In the network tab, copy the first request to onsen.ag with "copy as cURL"
- From the copied string, match the pattern `_session_id=SESSION_STRING_KEEP_IT_SECURE` without the `_session_id=` prefix.

## onsengo ls

`onsengo ls` can list radio shows and episodes on onsen.ag. Execute without arguments gives you all the radio programs the website current has:

```
~/w/onsengo ❯❯❯ ./onsengo ls
...
(omitted)
...
d--*--  4 Apr 13 2021 toshitai        セブン-イレブン presents 佐倉としたい大西
d--*--  4 Apr 13 2021 kakazu          かかずゆみの超輝け！やまと魂！！
d--*--  2 Apr 14 2021 vivy            Vivy -Flourite Eye’s Radio- 
d--*-- 33 Apr 14 2021 yuukiyui        ゆうきとゆいのラジオで２人暮らし♡
d--*--  8 Apr 14 2021 matsui          松井恵理子のにじらじっ！
d--*--  8 Apr 14 2021 frasta          笠間淳・梶原岳人のふらっと紀行！ ～…え？スタジオからは出られないんですか？～
d--*--  8 Apr 14 2021 mushinobu_radio 東海オンエア虫眼鏡・島﨑信長　声YouラジオZ
d--*--  8 Apr 14 2021 jks             会沢紗弥と花井美春の「まったく、女子高生は最高だぜ！！」
d--*--  4 Apr 14 2021 nonpetit        MoeMiののんびりプティフール
d--*--  1 Apr 14 2021 llss            ラブライブ！サンシャイン!! Aqours浦の星女学院RADIO!!!
d--*--  4 Apr 14 2021 saibou          一緒に「はたらく細胞」らじお
```

`onsengo ls -r` gives you all radio shows and their episodes.
 
Provide names to list only selected radio shows and episodes:

```
~/w/onsengo ❯❯❯ ./onsengo ls fujita toshitai gurepap
d----- 8 Apr  8 2021 gurepap       鷲崎健・藤田茜のグレパラジオP
-r-*-- 1 Apr  8 2021 gurepap/3897  第40回 予告 # 日高里菜
---*+$ 1 Apr  8 2021 gurepap/3898  第40回 本編 # 日高里菜
-----$ 1 Mar 25 2021 gurepap/3736  第39回 予告 # 高森奈津美
----+$ 1 Mar 25 2021 gurepap/3737  第39回 本編 # 高森奈津美
-----$ 1 Mar 11 2021 gurepap/3569  第38回 予告 # あじ秋刀魚
----+$ 1 Mar 11 2021 gurepap/3570  第38回 本編 # あじ秋刀魚
-----$ 1 Feb 25 2021 gurepap/3353  第37回 予告 # 山下七海
----+$ 1 Feb 25 2021 gurepap/3354  第37回 本編 # 山下七海
d----- 8 Apr  9 2021 fujita        藤田茜シーズン1
-rv*-- 1 Apr  9 2021 fujita/3919   第83回 予告
--v*+$ 1 Apr  9 2021 fujita/3920   第83回 本編
--v--$ 1 Mar 26 2021 fujita/3765   第82回 予告
--v-+$ 1 Mar 26 2021 fujita/3766   第82回 本編
--v--$ 1 Mar 12 2021 fujita/3598   第81回 予告
--v-+$ 1 Mar 12 2021 fujita/3599   第81回 本編
--v--$ 1 Feb 26 2021 fujita/3383   第80回 予告
--v-+$ 1 Feb 26 2021 fujita/3384   第80回 本編
d--*-- 4 Apr 13 2021 toshitai      セブン-イレブン presents 佐倉としたい大西
-r-*-- 1 Apr 13 2021 toshitai/3946 第263回
-----$ 1 Apr  6 2021 toshitai/3873 第262回
-----$ 1 Mar 30 2021 toshitai/3796 第261回
-----$ 1 Mar 23 2021 toshitai/3708 第260回
```

- `drv*+$`:
  - `d`: indicates the entry is a radio or episode
  - `r`: whether or not the current ***session*** can play the radio episode
  - `v`: includes video stream
  - `*`: just updated
  - `+`: extra content (sometimes extra is main content)
  - `$`: paid content
- For radios, output is sort by upload date. (no perform sorting on episodes)
