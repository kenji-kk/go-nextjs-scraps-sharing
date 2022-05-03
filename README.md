## 概要
Next.jsとGoAPIでアプリを作るためのテンプレート。認証機能まで実装している。

## 構築コマンド
### ・バックエンドAPI (このリポジトリ)
1. `git clone git@github.com:kenji-kk/go-nextjs-scraps-sharing`
2. `cd go-nextjs-scraps-sharing`
3. `docker-compose up`
### ・フロントエンド ( https://github.com/kenji-kk/nextjs-go-scraps-sharing )
1. `git clone git@github.com:kenji-kk/nextjs-go-scraps-sharing`
2. `cd nextjs-go-scraps-sharing`
3. `touch .env`
4. `echo NEXT_PUBLIC_HOST="http://localhost" > .env` *注意：GCEなどクラウドコンピューティングを使用の場合はlocalhostではなくそちらのipアドレスを指定
5. `docker-compose up`
＊一回で立ち上がらなかったらもう一度３を繰り返す


## 確認URL
- http://localhost (Nginx経由)
